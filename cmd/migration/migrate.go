package migration

import (
	"github.com/onetwoclimb/internal/storages/migrate"
)

var migrationConfig migrate.MigrationConfig

var Migration = migrate.GetMigrationCommand("migration-onetwoclimb", &migrationConfig, migrate.GetMigrations)
