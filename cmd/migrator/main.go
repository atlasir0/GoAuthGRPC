package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3" // Добавьте этот импорт
)

func main() {
	var storagePath, migrationsPath, migrationTable string

	// Устанавливаем значение по умолчанию для storagePath
	storagePath = "./storage/sso.db"

	flag.StringVar(&migrationsPath, "migrations-path", "", "path migrations ")
	flag.StringVar(&migrationTable, "migration-table", "migrations", "name of migration table")
	flag.Parse()

	if migrationsPath == "" {
		panic("migrations path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite://%s/?x-migrate-table=%s", storagePath, migrationTable),
	)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")
		} else {
			panic(err)
		}
	}
	fmt.Println("Migrations applied")
}
