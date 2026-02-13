package postgres

import (
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

var (
	//go:embed migrations/shortener/*.sql
	migrations embed.FS
)

func MigrateDB(db *Postgres) error {
	if err := migrate(db, "migrations/shortener"); err != nil {
		return fmt.Errorf("err on migrate: %v", err)
	}

	return nil
}

func migrate(db *Postgres, dir string) error {
	goose.SetBaseFS(migrations)

	if err := goose.Up(db.SqlDB(), dir); err != nil {
		return fmt.Errorf("goose up: %v", err)
	}

	return nil
}
