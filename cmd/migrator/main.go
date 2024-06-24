package main

import (
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	var migrationsPath, migrationTable string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationTable, "migration-table", "migrations", "name of migration table")
	flag.Parse()

	if migrationsPath == "" {
		log.Fatal("migrations path is required")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&x-migrate-table=%s",
		"postgres", "123", "localhost", 5432, "todo_list", migrationTable)

	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrate: %v", err)
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to apply")
		} else {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
	}
	fmt.Println("Migrations applied successfully")
}
