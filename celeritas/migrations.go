package celeritas

import (
	"log"

	"github.com/gobuffalo/pop"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (c *Celeritas) PopConnect() (*pop.Connection, error) {
	tx, err := pop.Connect("development")
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// the buffalo binaries allow empty migration
// but here we want to add few arguments
func (c *Celeritas) CreatePopMigration(up, down []byte, migrationName, migrationType string) error {
	var migrationPath = c.RootPath + "/migrations"
	err := pop.MigrationCreate(migrationPath, migrationName, migrationType, up, down)
	if err != nil {
		return err
	}
	return nil
}

func (c *Celeritas) RunPopMigrations(tx *pop.Connection) error {

	var migrationPath = c.RootPath + "/migrations"

	// create a file migrator
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

// variatic parameters can receive 0 or more variables
// so it's a good work around for potential inexistant parameters
func (c *Celeritas) PopMigrateDown(tx *pop.Connection, steps ...int) error {

	var migrationPath = c.RootPath + "/migrations"
	step := 1
	if len(steps) > 0 {
		//  if steps == 2 as parameter
		// then step will be == 2
		step = steps[0]
	}

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fm.Down(step)
	if err != nil {
		return err
	}
	return nil
}

func (c *Celeritas) PopMigrateReset(tx *pop.Connection) error {

	var migrationPath = c.RootPath + "/migrations"

	fm, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fm.Reset()
	if err != nil {
		return err
	}
	return nil
}

// migration working with the package "github.com/golang-migrate/migrate/v4"
func (c *Celeritas) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		log.Println("Error running migration:", err)
		return err
	}
	return nil
}

func (c *Celeritas) MigrateDownAll(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) Steps(n int, dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(n); err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Force(-1); err != nil {
		return err
	}

	return nil
}
