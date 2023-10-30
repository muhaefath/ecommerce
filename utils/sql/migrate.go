package sql

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type migrateDriver struct {
	migratePath string
	seedPath    string
	db          DBer
}

// Migrate is the DB migrator.
type Migrate struct {
	*migrate.Migrate
	migrateDriver
}

// NewMigrate initialises the DB migrator with open connection.
func NewMigrate(db DBer, migrateDBPath, seedDBPath string) (*Migrate, error) {
	md := migrateDriver{}
	md.db = db
	md.migratePath = migrateDBPath
	md.seedPath = seedDBPath

	m, err := migrate.New(fmt.Sprintf("file://%s", migrateDBPath), db.Config().URI)
	if err != nil {
		return nil, err
	}

	return &Migrate{m, md}, nil
}

// Up is method to start migrate database
func (m *Migrate) Up() error {
	return m.Migrate.Up()
}

// SeedDB is method to seed the database
func (m *Migrate) SeedDB() error {
	files, err := ioutil.ReadDir(m.seedPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := ioutil.ReadFile(path.Join(m.seedPath, file.Name()))
		if err != nil {
			return err
		}
		_, _ = m.db.Exec(string(content))
	}

	return nil
}

// Down is method to purge migrated database
func (m *Migrate) Down() error {
	return m.Migrate.Down()
}

// Rollback is method to rolling back migration to 1 step behind
func (m *Migrate) Rollback() error {
	return m.Migrate.Steps(-1)
}

// MigrateStatus returns the migration status.
func (m *Migrate) MigrateStatus() ([][]string, error) {
	var status [][]string

	lastMigratedVersion, _, _ := m.Migrate.Version()
	files, err := ioutil.ReadDir(m.seedPath)
	if err != nil {
		return status, err
	}

	for _, entry := range files {
		if !entry.IsDir() {
			fn := entry.Name()
			splits := strings.Split(fn, "_")
			if strings.HasSuffix(fn, ".up.sql") && len(splits) > 1 {
				version, err := strconv.Atoi(splits[0])
				if err != nil {
					return nil, err
				}

				state := "no"
				if version <= int(lastMigratedVersion) {
					state = "yes"
				}

				status = append(status, []string{fmt.Sprintf("%s/%s", m.migratePath, fn), state})
			}
		}
	}

	return status, nil
}
