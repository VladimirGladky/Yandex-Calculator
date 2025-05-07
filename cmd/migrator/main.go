package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var StoragePath, MigrationsPath, MigrationsTable string
	flag.StringVar(&StoragePath, "storage-path", "", "storage path")
	flag.StringVar(&MigrationsPath, "migrations-path", "", "migrations path")
	flag.StringVar(&MigrationsTable, "migrations-table", "", "migrations table")
	flag.Parse()

	if StoragePath == "" {
		panic("storage path is empty")
	}
	if MigrationsPath == "" {
		panic("migrations path is empty")
	}
	m, err := migrate.New("file://"+MigrationsPath, fmt.Sprintf("sqlite3://%s", StoragePath))
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return
		}
		panic(err)
	}
	fmt.Println("migrations applied")
}
