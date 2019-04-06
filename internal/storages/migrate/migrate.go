package migrate

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/spf13/pflag"

	"github.com/onetwoclimb/internal/utils"

	"github.com/olekukonko/tablewriter"
	"github.com/onetwoclimb/internal/storages"
	migrate "github.com/rubenv/sql-migrate"
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
	Step  int
}

func (mc *MigrationConfig) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("MigrationCfg", pflag.PanicOnError)

	f.StringVar(&mc.Mode, "mode", "", "db mode release,debug,test")
	f.AddFlagSet(mc.MySQL.Flags("my_sql"))
	f.IntVar(&mc.Step, "migration_step", 0, "migration step")
	return f
}

func migrateCommandHandler(migrationTable string, c *MigrationConfig, f func() []*migrate.Migration) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		utils.BindEnv(cmd)

		mySql, err := storages.New(&c.MySQL)
		if err != nil {
			log.Fatalln("error while creating db connection", err)
		}

		migrate.SetTable(migrationTable)

		if cmd.Name() == "up" {
			doMigrate(mySql.DB(), migrate.Up, c.Step, f)
		} else if cmd.Name() == "down" {
			doMigrate(mySql.DB(), migrate.Down, c.Step, f)
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
	migrateCmd.PersistentFlags().IntVar(&config.Step, "limit", 0, "Maximum migration steps (--limit 0)")

	return migrateCmd
}

func GetMigrations() []*migrate.Migration {
	return []*migrate.Migration{
		{
			Id: "1",
			Up: []string{`
				--
				-- Table structure for table gyms
				--
				CREATE TABLE IF NOT EXISTS gyms (
  					id int(11) NOT NULL AUTO_INCREMENT,
  					name varchar(255) NOT NULL,
  					country_id int(11) NOT NULL,
  					city varchar(255) NOT NULL,
  					address text NOT NULL,
  					imageURL text NOT NULL,
  					lat varchar(255) NOT NULL,
					` + "`long`" + ` varchar(255) NOT NULL, -- dirty hack
  					description text NOT NULL,
  					fb_link text NOT NULL,
  					web_link text NOT NULL,
  					phone varchar(255) NOT NULL,
  					photo1 text NOT NULL,
  					photo2 text NOT NULL,
  					footer text NOT NULL,
  					hasBoulder tinyint(1) DEFAULT NULL,
  					hasClimbingWall tinyint(1) DEFAULT NULL,
  					login varchar(255) NOT NULL,
  					pass varchar(255) NOT NULL,
  					isVerified tinyint(1) NOT NULL,
  					PRIMARY KEY (id),
  					KEY country_id (country_id)
				) ENGINE=InnoDB AUTO_INCREMENT=115 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS gyms;`},
		},
		{
			Id: "2",
			Up: []string{`
			--
			-- Table structure for table events
			--
			CREATE TABLE IF NOT EXISTS events (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				name varchar(255) NOT NULL,
  				address varchar(255) NOT NULL,
  				shortDescriprion varchar(255) NOT NULL,
  				descriprion text NOT NULL,
  				dateFrom timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  				dateTill timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- here was zero, guess it's not crucial to have current date as default'
  				image text NOT NULL,
  				price varchar(255) NOT NULL,
  				fbLink text NOT NULL,
  				isFlashEnable tinyint(1) NOT NULL,
  				gym_id int(11) DEFAULT NULL,
  				PRIMARY KEY (id),
  				KEY gym_id (gym_id),
  				KEY gym_id_2 (gym_id),
  				CONSTRAINT events_ibfk_1 FOREIGN KEY (gym_id) REFERENCES gyms (id)
			) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS events;`},
		},
		{
			Id: "3",
			Up: []string{`
			--
			-- Table structure for table attempt_type
			--
			CREATE TABLE IF NOT EXISTS attempt_type (
  				id int(11) NOT NULL AUTO_INCREMENT,
				name varchar(11) NOT NULL,
				PRIMARY KEY (id)
				) ENGINE=InnoDB DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS attempt_type;`},
		},
		{
			Id: "4",
			Up: []string{`
			--
			-- Table structure for table board_size
			--
			CREATE TABLE IF NOT EXISTS board_size (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				parameter varchar(255) NOT NULL,
  				def varchar(255) NOT NULL,
  				PRIMARY KEY (id)
			) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS board_size;`},
		},
		{
			Id: "5",
			Up: []string{`
			--
			-- Table structure for table climb_type
			--
			CREATE TABLE IF NOT EXISTS climb_type (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				name varchar(255) NOT NULL,
  				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS climb_type;`},
		},
		{
			Id: "6",
			Up: []string{`
			--
			-- Table structure for table colors
			--
			CREATE TABLE IF NOT EXISTS colors (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				name varchar(20) CHARACTER SET latin1 NOT NULL,
  				pin_code varchar(3) CHARACTER SET latin1 NOT NULL,
  				hex varchar(255) NOT NULL,
  			PRIMARY KEY (id)
			) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COMMENT='цвета стенда';
			`},
			Down: []string{`DROP TABLE IF EXISTS colors;`},
		},
		{
			Id: "7",
			Up: []string{`
			--
			-- Table structure for table complexity
			--
			CREATE TABLE IF NOT EXISTS complexity (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				name varchar(20) NOT NULL,
  				color varchar(255) NOT NULL,
  				points int(11) NOT NULL,
  				calories int(11) NOT NULL,
  				PRIMARY KEY (id)
			) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS complexity;`},
		},
		{
			Id: "8",
			Up: []string{`
			--
			-- Table structure for table countries
			--
			CREATE TABLE IF NOT EXISTS countries (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				name varchar(44) DEFAULT NULL,
  				code varchar(4) DEFAULT NULL,
  				PRIMARY KEY (id),
  				KEY id (id)
			) ENGINE=InnoDB AUTO_INCREMENT=251 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS countries;`},
		},
		{
			Id: "9",
			Up: []string{`
			--
			-- Table structure for table event_attempt_type
			--
			CREATE TABLE IF NOT EXISTS event_attempt_type (
  				id int(11) NOT NULL AUTO_INCREMENT,
				name varchar(255) NOT NULL,
  				PRIMARY KEY (id)
			) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS event_attempt_type;`},
		},
		{
			Id: "10",
			Up: []string{`
			--
			-- Table structure for table event_route_colors
			--
			CREATE TABLE IF NOT EXISTS event_route_colors (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				name varchar(255) NOT NULL,
  				hex varchar(255) NOT NULL,
  				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS event_route_colors;`},
		},
		{
			Id: "11",
			Up: []string{`
			--
			-- Table structure for table gym_hours
			--
			CREATE TABLE IF NOT EXISTS gym_hours (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				gym_id int(11) NOT NULL,
  				monday varchar(255) NOT NULL,
  				tuesday varchar(255) NOT NULL,
  				wednesday varchar(255) NOT NULL,
  				thursday varchar(255) NOT NULL,
  				friday varchar(255) NOT NULL,
  				saturday varchar(255) NOT NULL,
  				sunday varchar(255) NOT NULL,
  				PRIMARY KEY (id),
  				KEY gym_id (gym_id),
  				CONSTRAINT ck_hours_gym FOREIGN KEY (gym_id) REFERENCES gyms (id)
				) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS gym_hours;`},
		},
		{
			Id: "12",
			Up: []string{`
			--
			-- Table structure for table levels
			--
			CREATE TABLE IF NOT EXISTS levels (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				name varchar(255) NOT NULL,
  				PRIMARY KEY (id)
			) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS levels;`},
		},
		{
			Id: "13",
			Up: []string{`
			--
			-- Table structure for table options
			--
			CREATE TABLE IF NOT EXISTS options (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				name varchar(255) NOT NULL,
  				image_url varchar(255) NOT NULL,
  				PRIMARY KEY (id)
			) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COMMENT='Options in gym';
			`},
			Down: []string{`DROP TABLE IF EXISTS options;`},
		},
		{
			Id: "14",
			Up: []string{`
			--
			-- Table structure for table shop
			--
			CREATE TABLE IF NOT EXISTS shop (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				header text NOT NULL,
  				price text NOT NULL,
  				currency varchar(5) NOT NULL,
  				link text NOT NULL,
  				image text NOT NULL,
  				PRIMARY KEY (id)
			) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS shop;`},
		},
		{
			Id: "15",
			Up: []string{`
			--
			-- Table structure for table users
			--
			CREATE TABLE IF NOT EXISTS users (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				first_name varchar(255) NOT NULL,
  				last_name varchar(255) NOT NULL,
  				gender varchar(255) DEFAULT NULL,
  				fb_id varchar(255) DEFAULT NULL,
  				fb_token varchar(255) DEFAULT NULL,
  				email varchar(255) DEFAULT NULL,
  				picURL text NOT NULL,
  				google_id varchar(255) DEFAULT NULL,
  				password varchar(255) DEFAULT NULL,
  				rating int(11) NOT NULL DEFAULT '0',
  				rank int(11) NOT NULL DEFAULT '0',
  				rank_last int(11) NOT NULL DEFAULT '0',
				level_id int(11) NOT NULL DEFAULT '3',
  				isAdmin tinyint(1) NOT NULL DEFAULT '0',
  				PRIMARY KEY (id),
  				UNIQUE KEY fb_id (fb_id),
  				KEY level_id (level_id),
  				KEY level_id_2 (level_id),
  				CONSTRAINT user_levels_ck FOREIGN KEY (level_id) REFERENCES levels (id)
			) ENGINE=InnoDB AUTO_INCREMENT=787 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS users;`},
		},
		{
			Id: "16",
			Up: []string{`
			--
			-- Table structure for table outdoors
			--
			CREATE TABLE IF NOT EXISTS outdoors (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				name varchar(255) NOT NULL,
  				lat varchar(255) NOT NULL,
  				` + "`long`" + ` varchar(255) NOT NULL,
  				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS outdoors;`},
		},
		{
			Id: "17",
			Up: []string{`
			--
			-- Table structure for table gym_options
			--
			CREATE TABLE IF NOT EXISTS gym_options (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				gym_id int(11) NOT NULL,
  				option_id int(11) NOT NULL,
  				PRIMARY KEY (id),
  				KEY gym_id (gym_id,option_id),
  				KEY option_id (option_id),
  				CONSTRAINT ck_gym_options_option FOREIGN KEY (option_id) REFERENCES options (id),
  				CONSTRAINT ck_gym_options_gym FOREIGN KEY (gym_id) REFERENCES gyms (id)
			) ENGINE=InnoDB AUTO_INCREMENT=66 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS gym_options;`},
		},
		{
			Id: "18",
			Up: []string{`
			--
			-- Table structure for table event_levels
			--
			CREATE TABLE IF NOT EXISTS event_levels (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				event_id int(11) NOT NULL,
  				level_id int(11) NOT NULL,
  				description text NOT NULL,
  				PRIMARY KEY (id),
  				KEY event_id (event_id),
  				KEY level_id (level_id),
  				CONSTRAINT event_levels_ibfk_1 FOREIGN KEY (event_id) REFERENCES events (id),
  				CONSTRAINT event_levels_ibfk_2 FOREIGN KEY (level_id) REFERENCES levels (id)
			) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS event_levels;`},
		},
		{
			Id: "19",
			Up: []string{`
			--
			-- Table structure for table event_routes
			--
			CREATE TABLE IF NOT EXISTS event_routes (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				name varchar(255) NOT NULL,
  				event_id int(11) NOT NULL,
  				color_id int(11) DEFAULT NULL,
  				complexity_id int(11) NOT NULL,
  				topRate int(11) NOT NULL,
  				bonusRate int(11) DEFAULT NULL,
  				PRIMARY KEY (id),
  				KEY complexity_id (complexity_id),
  				KEY event_id (event_id),
  				KEY color_id (color_id),
  				CONSTRAINT event_routes_ibfk_1 FOREIGN KEY (event_id) REFERENCES events (id),
  				CONSTRAINT event_routes_ibfk_2 FOREIGN KEY (complexity_id) REFERENCES complexity (id),
  				CONSTRAINT event_routes_ibfk_3 FOREIGN KEY (color_id) REFERENCES event_route_colors (id)
			) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS event_routes;`},
		},
		{
			Id: "20",
			Up: []string{`
			--
			-- Table structure for table event_attempts
			--
			CREATE TABLE IF NOT EXISTS event_attempts (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				event_id int(11) NOT NULL,
  				user_id int(11) NOT NULL,
  				eventRoute_id int(11) NOT NULL,
  				isFlash tinyint(1) NOT NULL,
  				eventAttemptType_id int(11) NOT NULL,
  				PRIMARY KEY (id),
  				UNIQUE KEY unique_index_2_filds (user_id,eventRoute_id),
  				KEY user_id (user_id),
  				KEY eventRoute_id (eventRoute_id),
  				KEY eventAttemptType_id (eventAttemptType_id),
  				KEY event_id_3 (event_id),
  				CONSTRAINT event_attempts_ibfk_1 FOREIGN KEY (event_id) REFERENCES events (id),
  				CONSTRAINT event_attempts_ibfk_2 FOREIGN KEY (user_id) REFERENCES users (id),
  				CONSTRAINT event_attempts_ibfk_3 FOREIGN KEY (eventRoute_id) REFERENCES event_routes (id),
  				CONSTRAINT event_attempts_ibfk_4 FOREIGN KEY (eventAttemptType_id) REFERENCES event_attempt_type (id)
			) ENGINE=InnoDB AUTO_INCREMENT=90 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS event_attempts;`},
		},
		{
			Id: "21",
			Up: []string{`
			--
			-- Table structure for table event_users
			--
			CREATE TABLE IF NOT EXISTS event_users (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				user_id int(11) NOT NULL,
  				event_id int(11) NOT NULL,
  				user_level_id int(11) DEFAULT NULL,
  				isConfirmed tinyint(1) NOT NULL,
  				isEventAdmin tinyint(1) DEFAULT NULL,
  				isRouteSetter tinyint(1) NOT NULL,
  				isJudge tinyint(1) NOT NULL,
  				judgeCode text NOT NULL,
  				isOnline tinyint(1) NOT NULL DEFAULT '0',
  				PRIMARY KEY (id),
  				KEY user_id (user_id),
  				KEY user_level (user_level_id),
  				KEY event_id (event_id),
  				CONSTRAINT event_users_ibfk_1 FOREIGN KEY (user_id) REFERENCES users (id),
  				CONSTRAINT event_users_ibfk_3 FOREIGN KEY (event_id) REFERENCES events (id),
  				CONSTRAINT event_users_ibfk_4 FOREIGN KEY (user_level_id) REFERENCES levels (id)
			) ENGINE=InnoDB AUTO_INCREMENT=569 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS event_users;`},
		},
		{
			Id: "22",
			Up: []string{`
			--
			-- Table structure for table routes
			--
			CREATE TABLE IF NOT EXISTS routes (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				name varchar(70) NOT NULL,
  				comments text NOT NULL,
  				complexity_id int(11) NOT NULL,
  				author_id int(11) NOT NULL,
  				event_id int(11) DEFAULT NULL,
  				gym_id int(11) DEFAULT NULL,
  				is12C tinyint(1) NOT NULL,
  				imageURL text NOT NULL,
  				topRate int(11) NOT NULL,
  				bonusRate int(11) NOT NULL,
  				flashRate int(11) NOT NULL,
  				isArchive tinyint(1) NOT NULL,
  				isDefault tinyint(1) NOT NULL,
  				isJudged tinyint(1) NOT NULL,
  				isConfirmed tinyint(4) NOT NULL DEFAULT '0',
  				dt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  				hex_v1 text,
  				PRIMARY KEY (id),
  				KEY complexity_id (complexity_id),
  				KEY author_id (author_id),
  				KEY event_id (event_id),
  				KEY gym_id (gym_id),
  				CONSTRAINT author FOREIGN KEY (author_id) REFERENCES users (id),
  				CONSTRAINT complexity FOREIGN KEY (complexity_id) REFERENCES complexity (id),
  				CONSTRAINT FK718folur1fr8yq7n1rg56us16 FOREIGN KEY (author_id) REFERENCES users (id),
  				CONSTRAINT FKhjm6jo3j0rlmkrc1wyfinwi42 FOREIGN KEY (complexity_id) REFERENCES complexity (id),
  				CONSTRAINT routes_ibfk_1 FOREIGN KEY (event_id) REFERENCES events (id),
  				CONSTRAINT routes_ibfk_2 FOREIGN KEY (gym_id) REFERENCES gyms (id)
			) ENGINE=InnoDB AUTO_INCREMENT=2239 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS routes;`},
		},
		{
			Id: "23",
			Up: []string{`
			--
			-- Table structure for table holds
			--
			CREATE TABLE IF NOT EXISTS holds (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				route_id int(11) NOT NULL,
  				x int(11) NOT NULL,
  				y int(11) NOT NULL,
  				radius int(11) NOT NULL,
  				color_id int(11) NOT NULL,
  				PRIMARY KEY (id),
  				KEY route_id (route_id,color_id),
  				KEY color_id (color_id),
  				KEY route_id_2 (route_id),
  				KEY color_id_2 (color_id),
  				CONSTRAINT color FOREIGN KEY (color_id) REFERENCES colors (id),
  				CONSTRAINT FK4bfbw02mxlisfura0p50s32gk FOREIGN KEY (route_id) REFERENCES routes (id),
  				CONSTRAINT FKhy7536eph98d3bbwo1bqt6e1f FOREIGN KEY (color_id) REFERENCES colors (id),
  				CONSTRAINT rout FOREIGN KEY (route_id) REFERENCES routes (id)
			) ENGINE=InnoDB AUTO_INCREMENT=32354 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS holds;`},
		},
		{
			Id: "24",
			Up: []string{`
			--
			-- Table structure for table attempts
			--
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
  				PRIMARY KEY (id),
  				UNIQUE KEY id (id),
  				UNIQUE KEY route_user_event_unique_index (user_id,route_id,event_id),
  				KEY user_id (user_id),
  				KEY gym_id (gym_id),
  				KEY event_id (event_id),
				KEY route_id (route_id),
  				CONSTRAINT attempts_ibfk_1 FOREIGN KEY (user_id) REFERENCES users (id),
  				CONSTRAINT attempts_ibfk_3 FOREIGN KEY (gym_id) REFERENCES gyms (id),
  				CONSTRAINT attempts_ibfk_4 FOREIGN KEY (event_id) REFERENCES events (id),
  				CONSTRAINT attempts_ibfk_5 FOREIGN KEY (route_id) REFERENCES routes (id)
			) ENGINE=InnoDB AUTO_INCREMENT=12263 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS attempts;`},
		},
		{
			Id: "25",
			Up: []string{`
			--
			-- Table structure for table "attempt_log"
			--
			CREATE TABLE IF NOT EXISTS attempt_log (
  				id int(11) NOT NULL AUTO_INCREMENT,
  				isFlash tinyint(1) DEFAULT NULL,
  				isRedpoint tinyint(1) DEFAULT NULL,
  				isAttempt tinyint(1) DEFAULT NULL,
  				attemptId int(11) NOT NULL,
  				userId int(11) NOT NULL,
  				points int(11) NOT NULL,
  				datetime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  				date date NOT NULL,
  				PRIMARY KEY (id),
  				KEY attemptId (attemptId),
  				KEY userId (userId),
  				CONSTRAINT attempt_log_ibfk_1 FOREIGN KEY (attemptId) REFERENCES attempts (id),
  				CONSTRAINT attempt_log_ibfk_2 FOREIGN KEY (userId) REFERENCES users (id)
			) ENGINE=InnoDB AUTO_INCREMENT=16149 DEFAULT CHARSET=utf8;
			`},
			Down: []string{`DROP TABLE IF EXISTS attempt_log;`},
		},
	}
}
