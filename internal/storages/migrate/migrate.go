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
        	--
        	-- Database: '12climb'
        	-- Table structure for table 'attempts'
        	--

        	CREATE TABLE attempts (
				id int(11) NOT NULL,
				user_id int(11) NOT NULL,
				timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				gym_id int(11) DEFAULT NULL,
				event_id int(11) DEFAULT NULL,
				route_id int(11) NOT NULL,
				isFlash int(1) NOT NULL,
				attemptCount int(11) NOT NULL,
				redpointCount int(11) NOT NULL,
				isBonus tinyint(1) NOT NULL,
				isTop tinyint(1) NOT NULL
        	) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        	--
        	-- Triggers 'attempts'
        	--
        	DELIMITER $$
        	CREATE TRIGGER attempt_log_delete BEFORE DELETE ON attempts FOR EACH ROW BEGIN

          		DELETE FROM attempt_log WHERE attemptId=OLD.id;

          	UPDATE 'users' SET 'rating' = 'rating' - ( (OLD.redpointCount + OLD.isFlash) * ( SELECT cat.points
                                                                                           FROM attempts AS att
                                                                                                  INNER JOIN routes AS ro ON att.route_id = ro.id
                                                                                                  INNER JOIN complexity AS cat ON ro.complexity_id = cat.id
                                                                                           WHERE att.id =OLD.id ) )
          	WHERE  'users'.'id' =OLD.user_id;

          	SET @v1 := ( SELECT FIND_IN_SET( 'rating', (
            	SELECT GROUP_CONCAT( 'rating'
                	                 ORDER BY 'rating' DESC )
            FROM 'users' )
                                ) AS rank
                       FROM 'users'
                       WHERE id = OLD.user_id );

          UPDATE 'users'
          SET 'rank_last'=IF('rank' <> @v1, 'rank', 'rank_last'), 'rank'= @v1
          WHERE  'users'.'id' =OLD.user_id;

        END
        $$
        DELIMITER ;
        DELIMITER $$
        CREATE TRIGGER 'attempt_log_insert' AFTER INSERT ON 'attempts' FOR EACH ROW BEGIN
          INSERT INTO 'attempt_log'('isFlash', 'isRedpoint', 'isAttempt', 'attemptId', 'date', 'userId')
          VALUES (IF(NEW.isFlash > 0, 1, null),IF(NEW.redpointCount > 0, 1, null),IF(NEW.attemptCount > 0, 1, null),NEW.id,DATE(NEW.timestamp),NEW.user_id);

          IF (NEW.isFlash > 0 OR NEW.redpointCount > 0) THEN
            UPDATE 'users'
            SET  'rating' =  'rating' + ( SELECT cat.points
                                          FROM attempts AS att
                                                 INNER JOIN 'routes' AS ro ON att.route_id = ro.id
                                                 INNER JOIN 'complexity' AS cat ON ro.complexity_id = cat.id
                                          WHERE att.id =NEW.id )
            WHERE  'users'.'id'=NEW.user_id;

            SET @v1 := ( SELECT FIND_IN_SET( 'rating', (
              SELECT GROUP_CONCAT( 'rating'
                                   ORDER BY 'rating' DESC )
              FROM 'users' )
                                  ) AS rank
                         FROM 'users'
                         WHERE id = NEW.user_id );

            UPDATE 'users'
            SET 'rank_last'=IF('rank' <> @v1, 'rank', 'rank_last'), 'rank'= @v1
            WHERE  'users'.'id' =NEW.user_id;
          END IF;
        END
        $$
        DELIMITER ;
        DELIMITER $$
        CREATE TRIGGER 'attempt_log_update' AFTER UPDATE ON 'attempts' FOR EACH ROW BEGIN
          IF (NEW.redpointCount > OLD.redpointCount OR NEW.attemptCount > OLD.attemptCount OR NEW.isFlash > OLD.isFlash) THEN

            INSERT INTO 'attempt_log'('isFlash', 'isRedpoint', 'isAttempt', 'attemptId', 'date', 'userId')
            VALUES (IF(NEW.isFlash > OLD.isFlash, 1, null),IF(NEW.redpointCount > OLD.redpointCount, 1, null),IF(NEW.attemptCount > OLD.attemptCount, 1, null),NEW.id,DATE(NEW.timestamp),NEW.user_id);
          END IF;

          IF NEW.redpointCount > OLD.redpointCount THEN

            UPDATE 'users'
            SET  'rating' =  'rating' + ( SELECT cat.points
                                          FROM attempts AS att
                                                 INNER JOIN 'routes' AS ro ON att.route_id = ro.id
                                                 INNER JOIN 'complexity' AS cat ON ro.complexity_id = cat.id
                                          WHERE att.id =OLD.id )
            WHERE  'users'.'id' =OLD.user_id;

          END IF;

          IF NEW.redpointCount < OLD.redpointCount THEN

            UPDATE 'users'
            SET  'rating' =  'rating' - ( SELECT cat.points
                                          FROM attempts AS att
                                                 INNER JOIN 'routes' AS ro ON att.route_id = ro.id
                                                 INNER JOIN 'complexity' AS cat ON ro.complexity_id = cat.id
                                          WHERE att.id =OLD.id )
            WHERE  'users'.'id' =OLD.user_id;

          END IF;

          IF NEW.isFlash > OLD.isFlash THEN

            UPDATE 'users'
            SET  'rating' =  'rating' + ( SELECT cat.points
                                          FROM attempts AS att
                                                 INNER JOIN 'routes' AS ro ON att.route_id = ro.id
                                                 INNER JOIN 'complexity' AS cat ON ro.complexity_id = cat.id
                                          WHERE att.id =OLD.id )
            WHERE  'users'.'id' =OLD.user_id;

          END IF;

          IF NEW.isFlash < OLD.isFlash THEN

            UPDATE 'users'
            SET  'rating' =  'rating' - ( SELECT cat.points
                                          FROM attempts AS att
                                                 INNER JOIN 'routes' AS ro ON att.route_id = ro.id
                                                 INNER JOIN 'complexity' AS cat ON ro.complexity_id = cat.id
                                          WHERE att.id =OLD.id )
            WHERE  'users'.'id' =OLD.user_id;

          END IF;

          SET @v1 := ( SELECT FIND_IN_SET( 'rating', (
            SELECT GROUP_CONCAT( 'rating'
                                 ORDER BY 'rating' DESC )
            FROM 'users' )
                                ) AS rank
                       FROM 'users'
                       WHERE id = NEW.user_id );

          UPDATE 'users'
          SET 'rank_last'=IF('rank' <> @v1, 'rank', 'rank_last'), 'rank'= @v1
          WHERE  'users'.'id' =OLD.user_id;

        END
        $$
        DELIMITER ;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'attempt_log'
        --

        CREATE TABLE 'attempt_log' (
                                     'id' int(11) NOT NULL,
                                     'isFlash' tinyint(1) DEFAULT NULL,
                                     'isRedpoint' tinyint(1) DEFAULT NULL,
                                     'isAttempt' tinyint(1) DEFAULT NULL,
                                     'attemptId' int(11) NOT NULL,
                                     'userId' int(11) NOT NULL,
                                     'points' int(11) NOT NULL,
                                     'datetime' timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                     'date' date NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'attempt_type'
        --

        CREATE TABLE 'attempt_type' (
                                      'id' int(11) NOT NULL,
                                      'name' varchar(11) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'board_size'
        --

        CREATE TABLE 'board_size' (
                                    'id' int(11) NOT NULL,
                                    'parameter' varchar(255) NOT NULL,
                                    'def' varchar(255) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'climb_type'
        --

        CREATE TABLE 'climb_type' (
                                    'id' int(11) NOT NULL,
                                    'name' varchar(255) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'colors'
        --

        CREATE TABLE 'colors' (
                                'id' int(11) NOT NULL,
                                'name' varchar(20) CHARACTER SET latin1 NOT NULL,
                                'pin_code' varchar(3) CHARACTER SET latin1 NOT NULL,
                                'hex' varchar(255) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='цвета стенда';

        -- --------------------------------------------------------

        --
        -- Table structure for table 'complexity'
        --

        CREATE TABLE 'complexity' (
                                    'id' int(11) NOT NULL,
                                    'name' varchar(20) NOT NULL,
                                    'color' varchar(255) NOT NULL,
                                    'points' int(11) NOT NULL,
                                    'calories' int(11) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'countries'
        --

        CREATE TABLE 'countries' (
                                   'id' int(11) NOT NULL,
                                   'name' varchar(44) DEFAULT NULL,
                                   'code' varchar(4) DEFAULT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'events'
        --

        CREATE TABLE 'events' (
                                'id' int(11) NOT NULL,
                                'name' varchar(255) NOT NULL,
                                'address' varchar(255) NOT NULL,
                                'shortDescriprion' varchar(255) NOT NULL,
                                'descriprion' text NOT NULL,
                                'dateFrom' timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                'dateTill' timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
                                'image' text NOT NULL,
                                'price' varchar(255) NOT NULL,
                                'fbLink' text NOT NULL,
                                'isFlashEnable' tinyint(1) NOT NULL,
                                'gym_id' int(11) DEFAULT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'event_attempts'
        --

        CREATE TABLE 'event_attempts' (
                                        'id' int(11) NOT NULL,
                                        'event_id' int(11) NOT NULL,
                                        'user_id' int(11) NOT NULL,
                                        'eventRoute_id' int(11) NOT NULL,
                                        'isFlash' tinyint(1) NOT NULL,
                                        'eventAttemptType_id' int(11) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'event_attempt_type'
        --

        CREATE TABLE 'event_attempt_type' (
                                            'id' int(11) NOT NULL,
                                            'name' varchar(255) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'event_levels'
        --

        CREATE TABLE 'event_levels' (
                                      'id' int(11) NOT NULL,
                                      'event_id' int(11) NOT NULL,
                                      'level_id' int(11) NOT NULL,
                                      'description' text NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'event_routes'
        --

        CREATE TABLE 'event_routes' (
                                      'id' int(11) NOT NULL,
                                      'name' varchar(255) NOT NULL,
                                      'event_id' int(11) NOT NULL,
                                      'color_id' int(11) DEFAULT NULL,
                                      'complexity_id' int(11) NOT NULL,
                                      'topRate' int(11) NOT NULL,
                                      'bonusRate' int(11) DEFAULT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'event_route_colors'
        --

        CREATE TABLE 'event_route_colors' (
                                            'id' int(11) NOT NULL,
                                            'name' varchar(255) NOT NULL,
                                            'hex' varchar(255) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'event_users'
        --

        CREATE TABLE 'event_users' (
                                     'id' int(11) NOT NULL,
                                     'user_id' int(11) NOT NULL,
                                     'event_id' int(11) NOT NULL,
                                     'user_level_id' int(11) DEFAULT NULL,
                                     'isConfirmed' tinyint(1) NOT NULL,
                                     'isEventAdmin' tinyint(1) DEFAULT NULL,
                                     'isRouteSetter' tinyint(1) NOT NULL,
                                     'isJudge' tinyint(1) NOT NULL,
                                     'judgeCode' text NOT NULL,
                                     'isOnline' tinyint(1) NOT NULL DEFAULT '0'
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'gyms'
        --

        CREATE TABLE 'gyms' (
                              'id' int(11) NOT NULL,
                              'name' varchar(255) NOT NULL,
                              'country_id' int(11) NOT NULL,
                              'city' varchar(255) NOT NULL,
                              'address' text NOT NULL,
                              'imageURL' text NOT NULL,
                              'lat' varchar(255) NOT NULL,
                              'long' varchar(255) NOT NULL,
                              'description' text NOT NULL,
                              'fb_link' text NOT NULL,
                              'web_link' text NOT NULL,
                              'phone' varchar(255) NOT NULL,
                              'photo1' text NOT NULL,
                              'photo2' text NOT NULL,
                              'footer' text NOT NULL,
                              'hasBoulder' tinyint(1) DEFAULT NULL,
                              'hasClimbingWall' tinyint(1) DEFAULT NULL,
                              'login' varchar(255) NOT NULL,
                              'pass' varchar(255) NOT NULL,
                              'isVerified' tinyint(1) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'gym_hours'
        --

        CREATE TABLE 'gym_hours' (
                                   'id' int(11) NOT NULL,
                                   'gym_id' int(11) NOT NULL,
                                   'monday' varchar(255) NOT NULL,
                                   'tuesday' varchar(255) NOT NULL,
                                   'wednesday' varchar(255) NOT NULL,
                                   'thursday' varchar(255) NOT NULL,
                                   'friday' varchar(255) NOT NULL,
                                   'saturday' varchar(255) NOT NULL,
                                   'sunday' varchar(255) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'gym_options'
        --

        CREATE TABLE 'gym_options' (
                                     'id' int(11) NOT NULL,
                                     'gym_id' int(11) NOT NULL,
                                     'option_id' int(11) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'holds'
        --

        CREATE TABLE 'holds' (
                               'id' int(11) NOT NULL,
                               'route_id' int(11) NOT NULL,
                               'x' int(11) NOT NULL,
                               'y' int(11) NOT NULL,
                               'radius' int(11) NOT NULL,
                               'color_id' int(11) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'levels'
        --

        CREATE TABLE 'levels' (
                                'id' int(11) NOT NULL,
                                'name' varchar(255) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'options'
        --

        CREATE TABLE 'options' (
                                 'id' int(11) NOT NULL,
                                 'name' varchar(255) NOT NULL,
                                 'image_url' varchar(255) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Options in gym';

        -- --------------------------------------------------------

        --
        -- Table structure for table 'outdoors'
        --

        CREATE TABLE 'outdoors' (
                                  'id' int(11) NOT NULL,
                                  'name' varchar(255) NOT NULL,
                                  'lat' varchar(255) NOT NULL,
                                  'long' varchar(255) NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'routes'
        --

        CREATE TABLE 'routes' (
                                'id' int(11) NOT NULL,
                                'name' varchar(70) NOT NULL,
                                'comments' text NOT NULL,
                                'complexity_id' int(11) NOT NULL,
                                'author_id' int(11) NOT NULL,
                                'event_id' int(11) DEFAULT NULL,
                                'gym_id' int(11) DEFAULT NULL,
                                'is12C' tinyint(1) NOT NULL,
                                'imageURL' text NOT NULL,
                                'topRate' int(11) NOT NULL,
                                'bonusRate' int(11) NOT NULL,
                                'flashRate' int(11) NOT NULL,
                                'isArchive' tinyint(1) NOT NULL,
                                'isDefault' tinyint(1) NOT NULL,
                                'isJudged' tinyint(1) NOT NULL,
                                'isConfirmed' tinyint(4) NOT NULL DEFAULT '0',
                                'dt' timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                'hex_v1' text
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        --
        -- Triggers 'routes'
        --
        DELIMITER $$
        CREATE TRIGGER 'route_archived' AFTER UPDATE ON 'routes' FOR EACH ROW BEGIN
          IF (NEW.isArchive > OLD.isArchive) THEN

            IF ((SELECT count(id) FROM attempts WHERE route_id = OLD.id) = 0) THEN
              DELETE FROM holds WHERE route_id = OLD.id;
            END IF;
          END IF;
        END
        $$
        DELIMITER ;
        DELIMITER $$
        CREATE TRIGGER 'routes_delete' BEFORE DELETE ON 'routes' FOR EACH ROW BEGIN
          DELETE FROM holds WHERE route_id = OLD.id;
        END
        $$
        DELIMITER ;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'shop'
        --

        CREATE TABLE 'shop' (
                              'id' int(11) NOT NULL,
                              'header' text NOT NULL,
                              'price' text NOT NULL,
                              'currency' varchar(5) NOT NULL,
                              'link' text NOT NULL,
                              'image' text NOT NULL
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        -- --------------------------------------------------------

        --
        -- Table structure for table 'users'
        --

        CREATE TABLE 'users' (
                               'id' int(11) NOT NULL,
                               'first_name' varchar(255) NOT NULL,
                               'last_name' varchar(255) NOT NULL,
                               'gender' varchar(255) DEFAULT NULL,
                               'fb_id' varchar(255) DEFAULT NULL,
                               'fb_token' varchar(255) DEFAULT NULL,
                               'email' varchar(255) DEFAULT NULL,
                               'picURL' text NOT NULL,
                               'google_id' varchar(255) DEFAULT NULL,
                               'password' varchar(255) DEFAULT NULL,
                               'rating' int(11) NOT NULL DEFAULT '0',
                               'rank' int(11) NOT NULL DEFAULT '0',
                               'rank_last' int(11) NOT NULL DEFAULT '0',
                               'level_id' int(11) NOT NULL DEFAULT '3',
                               'isAdmin' tinyint(1) NOT NULL DEFAULT '0'
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

        --
        -- Indexes for dumped tables
        --

        --
        -- Indexes for table 'attempts'
        --
        ALTER TABLE 'attempts'
          ADD PRIMARY KEY ('id'),
          ADD UNIQUE KEY 'id' ('id'),
          ADD UNIQUE KEY 'route_user_event_unique_index' ('user_id','route_id','event_id'),
          ADD KEY 'user_id' ('user_id'),
          ADD KEY 'gym_id' ('gym_id'),
          ADD KEY 'event_id' ('event_id'),
          ADD KEY 'route_id' ('route_id');

        --
        -- Indexes for table 'attempt_log'
        --
        ALTER TABLE 'attempt_log'
          ADD PRIMARY KEY ('id'),
          ADD KEY 'attemptId' ('attemptId'),
          ADD KEY 'userId' ('userId');

        --
        -- Indexes for table 'attempt_type'
        --
        ALTER TABLE 'attempt_type'
          ADD PRIMARY KEY ('id');

        --
        -- Indexes for table 'board_size'
        --
        ALTER TABLE 'board_size'
          ADD PRIMARY KEY ('id');

        --
        -- Indexes for table 'climb_type'
        --
        ALTER TABLE 'climb_type'
          ADD PRIMARY KEY ('id');

        --
        -- Indexes for table 'colors'
        --
        ALTER TABLE 'colors'
          ADD PRIMARY KEY ('id');

        --
        -- Indexes for table 'complexity'
        --
        ALTER TABLE 'complexity'
          ADD PRIMARY KEY ('id');

        --
        -- Indexes for table 'countries'
        --
        ALTER TABLE 'countries'
          ADD PRIMARY KEY ('id'),
          ADD KEY 'id' ('id');

        --
        -- Indexes for table 'events'
        --
        ALTER TABLE 'events'
          ADD PRIMARY KEY ('id'),
          ADD KEY 'gym_id' ('gym_id'),
          ADD KEY 'gym_id_2' ('gym_id');

        --
        -- Indexes for table 'event_attempts'
        --
        ALTER TABLE 'event_attempts'
          ADD PRIMARY KEY ('id'),
          ADD UNIQUE KEY 'unique_index_2_filds' ('user_id','eventRoute_id'),
          ADD KEY 'user_id' ('user_id'),
          ADD KEY 'eventRoute_id' ('eventRoute_id'),
          ADD KEY 'eventAttemptType_id' ('eventAttemptType_id'),
          ADD KEY 'event_id_3' ('event_id');

        --
        -- Indexes for table 'event_attempt_type'
        --
        ALTER TABLE 'event_attempt_type'
          ADD PRIMARY KEY ('id');

        --
        -- Indexes for table 'event_levels'
        --
        ALTER TABLE 'event_levels'
          ADD PRIMARY KEY ('id'),
          ADD KEY 'event_id' ('event_id'),
          ADD KEY 'level_id' ('level_id');

        --
        -- Indexes for table 'event_routes'
        --
        ALTER TABLE 'event_routes'
          ADD PRIMARY KEY ('id'),
          ADD KEY 'complexity_id' ('complexity_id'),
          ADD KEY 'event_id' ('event_id'),
          ADD KEY 'color_id' ('color_id');

        --
        -- Indexes for table 'event_route_colors'
        --
        ALTER TABLE 'event_route_colors'
          ADD PRIMARY KEY ('id');

        --
        -- Indexes for table 'event_users'
        --
        ALTER TABLE 'event_users'
          ADD PRIMARY KEY ('id'),
          ADD KEY 'user_id' ('user_id'),
          ADD KEY 'user_level' ('user_level_id'),
          ADD KEY 'event_id' ('event_id');

        --
        -- Indexes for table 'gyms'
        --
        ALTER TABLE 'gyms'
          ADD PRIMARY KEY ('id'),
          ADD KEY 'country_id' ('country_id');

        --
        -- Indexes for table 'gym_hours'
        --
        ALTER TABLE 'gym_hours'
          ADD PRIMARY KEY ('id'),
          ADD KEY 'gym_id' ('gym_id');

        --
        -- Indexes for table 'gym_options'
        --
        ALTER TABLE 'gym_options'
          ADD PRIMARY KEY ('id'),
          ADD KEY 'gym_id' ('gym_id','option_id'),
          ADD KEY 'option_id' ('option_id');

        --
        -- Indexes for table 'holds'
        --
        ALTER TABLE 'holds'
          ADD PRIMARY KEY ('id'),
          ADD KEY 'route_id' ('route_id','color_id'),
          ADD KEY 'color_id' ('color_id'),
          ADD KEY 'route_id_2' ('route_id'),
          ADD KEY 'color_id_2' ('color_id');

        --
        -- Indexes for table 'levels'
        --
        ALTER TABLE 'levels'
          ADD PRIMARY KEY ('id');

        --
        -- Indexes for table 'options'
        --
        ALTER TABLE 'options'
          ADD PRIMARY KEY ('id');

        --
        -- Indexes for table 'outdoors'
        --
        ALTER TABLE 'outdoors'
          ADD PRIMARY KEY ('id');

        --
        -- Indexes for table 'routes'
        --
        ALTER TABLE 'routes'
          ADD PRIMARY KEY ('id'),
          ADD KEY 'complexity_id' ('complexity_id'),
          ADD KEY 'author_id' ('author_id'),
          ADD KEY 'event_id' ('event_id'),
          ADD KEY 'gym_id' ('gym_id');

        --
        -- Indexes for table 'shop'
        --
        ALTER TABLE 'shop'
          ADD PRIMARY KEY ('id');

        --
        -- Indexes for table 'users'
        --
        ALTER TABLE 'users'
          ADD PRIMARY KEY ('id'),
          ADD UNIQUE KEY 'fb_id' ('fb_id'),
          ADD KEY 'level_id' ('level_id'),
          ADD KEY 'level_id_2' ('level_id');

        --
        -- AUTO_INCREMENT for dumped tables
        --

        --
        -- AUTO_INCREMENT for table 'attempts'
        --
        ALTER TABLE 'attempts'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10261;
        --
        -- AUTO_INCREMENT for table 'attempt_log'
        --
        ALTER TABLE 'attempt_log'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13704;
        --
        -- AUTO_INCREMENT for table 'attempt_type'
        --
        ALTER TABLE 'attempt_type'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT;
        --
        -- AUTO_INCREMENT for table 'board_size'
        --
        ALTER TABLE 'board_size'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;
        --
        -- AUTO_INCREMENT for table 'climb_type'
        --
        ALTER TABLE 'climb_type'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT;
        --
        -- AUTO_INCREMENT for table 'colors'
        --
        ALTER TABLE 'colors'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;
        --
        -- AUTO_INCREMENT for table 'complexity'
        --
        ALTER TABLE 'complexity'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=39;
        --
        -- AUTO_INCREMENT for table 'countries'
        --
        ALTER TABLE 'countries'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=251;
        --
        -- AUTO_INCREMENT for table 'events'
        --
        ALTER TABLE 'events'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=18;
        --
        -- AUTO_INCREMENT for table 'event_attempts'
        --
        ALTER TABLE 'event_attempts'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=90;
        --
        -- AUTO_INCREMENT for table 'event_attempt_type'
        --
        ALTER TABLE 'event_attempt_type'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;
        --
        -- AUTO_INCREMENT for table 'event_levels'
        --
        ALTER TABLE 'event_levels'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=18;
        --
        -- AUTO_INCREMENT for table 'event_routes'
        --
        ALTER TABLE 'event_routes'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;
        --
        -- AUTO_INCREMENT for table 'event_route_colors'
        --
        ALTER TABLE 'event_route_colors'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT;
        --
        -- AUTO_INCREMENT for table 'event_users'
        --
        ALTER TABLE 'event_users'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=481;
        --
        -- AUTO_INCREMENT for table 'gyms'
        --
        ALTER TABLE 'gyms'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=113;
        --
        -- AUTO_INCREMENT for table 'gym_hours'
        --
        ALTER TABLE 'gym_hours'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;
        --
        -- AUTO_INCREMENT for table 'gym_options'
        --
        ALTER TABLE 'gym_options'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=63;
        --
        -- AUTO_INCREMENT for table 'holds'
        --
        ALTER TABLE 'holds'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=29152;
        --
        -- AUTO_INCREMENT for table 'levels'
        --
        ALTER TABLE 'levels'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;
        --
        -- AUTO_INCREMENT for table 'options'
        --
        ALTER TABLE 'options'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;
        --
        -- AUTO_INCREMENT for table 'outdoors'
        --
        ALTER TABLE 'outdoors'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT;
        --
        -- AUTO_INCREMENT for table 'routes'
        --
        ALTER TABLE 'routes'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1968;
        --
        -- AUTO_INCREMENT for table 'shop'
        --
        ALTER TABLE 'shop'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;
        --
        -- AUTO_INCREMENT for table 'users'
        --
        ALTER TABLE 'users'
          MODIFY 'id' int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=673;
        --
        -- Constraints for dumped tables
        --

        --
        -- Constraints for table 'attempts'
        --
        ALTER TABLE 'attempts'
          ADD CONSTRAINT 'attempts_ibfk_1' FOREIGN KEY ('user_id') REFERENCES 'users' ('id'),
          ADD CONSTRAINT 'attempts_ibfk_3' FOREIGN KEY ('gym_id') REFERENCES 'gyms' ('id'),
          ADD CONSTRAINT 'attempts_ibfk_4' FOREIGN KEY ('event_id') REFERENCES 'events' ('id'),
          ADD CONSTRAINT 'attempts_ibfk_5' FOREIGN KEY ('route_id') REFERENCES 'routes' ('id');

        --
        -- Constraints for table 'attempt_log'
        --
        ALTER TABLE 'attempt_log'
          ADD CONSTRAINT 'attempt_log_ibfk_1' FOREIGN KEY ('attemptId') REFERENCES 'attempts' ('id'),
          ADD CONSTRAINT 'attempt_log_ibfk_2' FOREIGN KEY ('userId') REFERENCES 'users' ('id');

        --
        -- Constraints for table 'events'
        --
        ALTER TABLE 'events'
          ADD CONSTRAINT 'events_ibfk_1' FOREIGN KEY ('gym_id') REFERENCES 'gyms' ('id');

        --
        -- Constraints for table 'event_attempts'
        --
        ALTER TABLE 'event_attempts'
          ADD CONSTRAINT 'event_attempts_ibfk_1' FOREIGN KEY ('event_id') REFERENCES 'events' ('id'),
          ADD CONSTRAINT 'event_attempts_ibfk_2' FOREIGN KEY ('user_id') REFERENCES 'users' ('id'),
          ADD CONSTRAINT 'event_attempts_ibfk_3' FOREIGN KEY ('eventRoute_id') REFERENCES 'event_routes' ('id'),
          ADD CONSTRAINT 'event_attempts_ibfk_4' FOREIGN KEY ('eventAttemptType_id') REFERENCES 'event_attempt_type' ('id');

        --
        -- Constraints for table 'event_levels'
        --
        ALTER TABLE 'event_levels'
          ADD CONSTRAINT 'event_levels_ibfk_1' FOREIGN KEY ('event_id') REFERENCES 'events' ('id'),
          ADD CONSTRAINT 'event_levels_ibfk_2' FOREIGN KEY ('level_id') REFERENCES 'levels' ('id');

        --
        -- Constraints for table 'event_routes'
        --
        ALTER TABLE 'event_routes'
          ADD CONSTRAINT 'event_routes_ibfk_1' FOREIGN KEY ('event_id') REFERENCES 'events' ('id'),
          ADD CONSTRAINT 'event_routes_ibfk_2' FOREIGN KEY ('complexity_id') REFERENCES 'complexity' ('id'),
          ADD CONSTRAINT 'event_routes_ibfk_3' FOREIGN KEY ('color_id') REFERENCES 'event_route_colors' ('id');

        --
        -- Constraints for table 'event_users'
        --
        ALTER TABLE 'event_users'
          ADD CONSTRAINT 'event_users_ibfk_1' FOREIGN KEY ('user_id') REFERENCES 'users' ('id'),
          ADD CONSTRAINT 'event_users_ibfk_3' FOREIGN KEY ('event_id') REFERENCES 'events' ('id'),
          ADD CONSTRAINT 'event_users_ibfk_4' FOREIGN KEY ('user_level_id') REFERENCES 'levels' ('id');

        --
        -- Constraints for table 'gym_hours'
        --
        ALTER TABLE 'gym_hours'
          ADD CONSTRAINT 'ck_hours_gym' FOREIGN KEY ('gym_id') REFERENCES 'gyms' ('id');

        --
        -- Constraints for table 'gym_options'
        --
        ALTER TABLE 'gym_options'
          ADD CONSTRAINT 'ck_gym_options_gym' FOREIGN KEY ('gym_id') REFERENCES 'gyms' ('id'),
          ADD CONSTRAINT 'ck_gym_options_option' FOREIGN KEY ('option_id') REFERENCES 'options' ('id');

        --
        -- Constraints for table 'holds'
        --
        ALTER TABLE 'holds'
          ADD CONSTRAINT 'FK4bfbw02mxlisfura0p50s32gk' FOREIGN KEY ('route_id') REFERENCES 'routes' ('id'),
          ADD CONSTRAINT 'FKhy7536eph98d3bbwo1bqt6e1f' FOREIGN KEY ('color_id') REFERENCES 'colors' ('id'),
          ADD CONSTRAINT 'color' FOREIGN KEY ('color_id') REFERENCES 'colors' ('id'),
          ADD CONSTRAINT 'rout' FOREIGN KEY ('route_id') REFERENCES 'routes' ('id');

        --
        -- Constraints for table 'routes'
        --
        ALTER TABLE 'routes'
          ADD CONSTRAINT 'FK718folur1fr8yq7n1rg56us16' FOREIGN KEY ('author_id') REFERENCES 'users' ('id'),
          ADD CONSTRAINT 'FKhjm6jo3j0rlmkrc1wyfinwi42' FOREIGN KEY ('complexity_id') REFERENCES 'complexity' ('id'),
          ADD CONSTRAINT 'author' FOREIGN KEY ('author_id') REFERENCES 'users' ('id'),
          ADD CONSTRAINT 'complexity' FOREIGN KEY ('complexity_id') REFERENCES 'complexity' ('id'),
          ADD CONSTRAINT 'routes_ibfk_1' FOREIGN KEY ('event_id') REFERENCES 'events' ('id'),
          ADD CONSTRAINT 'routes_ibfk_2' FOREIGN KEY ('gym_id') REFERENCES 'gyms' ('id');

        --
        -- Constraints for table 'users'
        --
        ALTER TABLE 'users'
          ADD CONSTRAINT 'user_levels_ck' FOREIGN KEY ('level_id') REFERENCES 'levels' ('id');

        /*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
        /*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
        /*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
			`},
			Down: []string{`
				DROP TABLE IF EXISTS colors CASCADE;
			`},
		},
	}

}
