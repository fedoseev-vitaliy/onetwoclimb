package migration

import (
	"github.com/onetwoclimb/internal/storages/migrate"
)

var migrationConfig migrate.MigrationConfig

func init() {
	Migration.Flags().AddFlagSet(migrationConfig.Flags())
}

var Migration = migrate.GetMigrationCommand("migration-onetwoclimb", &migrationConfig, migrate.GetMigrations)
