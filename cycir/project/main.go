package main

import (
	"log"
	"github.com/gobuffalo/pop"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	tx, err := pop.Connect("development")
	if err != nil {
		log.Fatal(err)
	}
	err = runPopMigrations(tx)
		if err != nil {
			log.Fatal(err)
		}
}

func runPopMigrations(tx *pop.Connection) error {
	var migrationPath = "./migrations"

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fm.Up()
	if err != nil {
		return err
	}
	return nil
}