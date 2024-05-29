package database

import (
	"context"
	"fmt"

	migration "userservice/migration"

	"github.com/uptrace/bun/migrate"
	_ "github.com/uptrace/bun/migrate"
)

var migrations = migrate.NewMigrations()

func init() {
	// Register SQL migrations.
	if err := migrations.Discover(migration.SqlMigrations); err != nil {
		panic(err)
	}
}

func RunMigrations(ctx context.Context) error {
	migrator := migrate.NewMigrator(db, migrations)

	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if group.IsZero() {
		fmt.Println("there are no new migrations to run")
		return nil
	}

	fmt.Printf("migrated to %s\n", group)
	return nil
}

func RollbackMigration(ctx context.Context) error {
	migrator := migrate.NewMigrator(db, migrations)

	group, err := migrator.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	if group.IsZero() {
		fmt.Println("there are no groups to roll back")
		return nil
	}

	fmt.Printf("rolled back %s\n", group)
	return nil
}
