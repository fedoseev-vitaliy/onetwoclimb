package migrate

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/onetwoclimb/internal/storages"
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

type statusRow struct {
	ID        string
	Migrated  bool
	AppliedAt time.Time
}

type MigrationConfig struct {
	Mode  string
	MySQL storages.Config
	Limit int
}

func migrateCommandHandler(migrationTable string, c *MigrationConfig, f func() []*migrate.Migration) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		flag.Parse()

		mySql, err := storages.New(&c.MySQL)
		if err != nil {
			log.Fatalln("error while creating db connection", err)
		}

		migrate.SetTable(migrationTable)

		if cmd.Name() == "up" {
			doMigrate(mySql.DB(), migrate.Up, c.Limit, f)
		} else if cmd.Name() == "down" {
			doMigrate(mySql.DB(), migrate.Down, c.Limit, f)
		} else {
			getMigrateStatus(mySql.DB(), f)
		}
	}
}

func doMigrate(pg *sql.DB, direction migrate.MigrationDirection, max int, f func() []*migrate.Migration) {
	migrations := f()

	source := migrate.MemoryMigrationSource{
		Migrations: migrations,
	}

	num, err := migrate.ExecMax(pg, "mysql", source, direction, max)
	if err != nil {
		log.Fatalln("migration error", err)
	}

	log.Println("Applied", num)
}

func getMigrateStatus(mySql *sql.DB, f func() []*migrate.Migration) {
	migrations := f()

	records, err := migrate.GetMigrationRecords(mySql, "mysql")
	if err != nil {
		log.Fatalln("error while getting migration", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Migration", "Applied"})
	table.SetColWidth(60)

	rows := make(map[string]*statusRow)
	for _, m := range migrations {
		rows[m.Id] = &statusRow{
			ID:       m.Id,
			Migrated: false,
		}
	}

	for _, r := range records {
		//When migration exist in DB but not in project
		if _, exist := rows[r.Id]; !exist {
			rows[r.Id] = &statusRow{
				ID:       r.Id,
				Migrated: false,
			}
		}

		rows[r.Id].Migrated = true
		rows[r.Id].AppliedAt = r.AppliedAt
	}

	for _, m := range rows {
		if m.Migrated {
			table.Append([]string{
				m.ID,
				m.AppliedAt.String(),
			})
		} else {
			table.Append([]string{
				m.ID,
				"no",
			})
		}
	}

	table.Render()
}

func GetMigrationCommand(migrationTable string, config *MigrationConfig, f func() []*migrate.Migration) *cobra.Command {
	handler := migrateCommandHandler(migrationTable, config, f)

	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Database migrations command center",
	}

	var migrateUpCmd = &cobra.Command{
		Use:   "up",
		Short: "Apply migration",
		Run:   handler,
	}

	var migrateDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Downgrade migration",
		Run:   handler,
	}

	var migrateStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Get database state",
		Run:   handler,
	}

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)

	migrateCmd.PersistentFlags().StringVar(&config.Mode, "mode", "debug", "release,debug,test")
	migrateCmd.PersistentFlags().AddFlagSet(config.MySQL.Flags("mysql"))
	migrateCmd.PersistentFlags().IntVar(&config.Limit, "limit", 0, "Maximum migration steps (--limit 0)")

	return migrateCmd
}

func GetMigrations() []*migrate.Migration {
	return []*migrate.Migration{
		{
			Id: "1",
			Up: []string{`
				CREATE TABLE IF NOT EXISTS attempts (
					id int(11) NOT NULL AUTO_INCREMENT,
					user_id int(11) NOT NULL,
					timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
					gym_id int(11) DEFAULT NULL,
					event_id int(11) DEFAULT NULL,
					route_id int(11) NOT NULL,
					isFlash int(1) NOT NULL,
					attemptCount int(11) NOT NULL,
					redpointCount int(11) NOT NULL,
					isBonus tinyint(1) NOT NULL,
					isTop tinyint(1) NOT NULL,
					PRIMARY KEY id,
					UNIQUE KEY id (id),
					UNIQUE KEY route_user_event_unique_index (user_id, route_id, event_id),
					KEY user_id (user_id),
					KEY gym_id (gym_id),
					KEY event_id (event_id),
					KEY route_id (route_id)
				)
			`},
			Down: []string{`
				DROP TABLE IF EXISTS colors CASCADE;
			`},
		},
	}

}
