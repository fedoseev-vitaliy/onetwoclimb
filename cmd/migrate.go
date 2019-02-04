package cmd

import (
	"github.com/onetwoclimb/internal/storages/migrate"
)

var migrationConfig migrate.MigrationConfig

func init() {
	RootCmd.AddCommand(migrate.GetMigrationCommand("migration-onetwoclimb", &migrationConfig, migrate.GetMigrations))
}
