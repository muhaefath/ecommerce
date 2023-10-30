package sql

import (
	"embed"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	DBSeedPath = "db/seed"
)

// Seed is the DB seeder.
type Seed struct {
	config  *Config
	db      *sqlx.DB
	embedFS embed.FS
	files   []string
}

// NewSeed initialises the DB seeder with open connection.
func NewSeed(config *Config, embedFS embed.FS) (*Seed, error) {
	dir := fmt.Sprintf("%s/%s", DBSeedPath, config.Name)
	entries, err := embedFS.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := []string{}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") && entry.Name() != embedInit {
			files = append(files, fmt.Sprintf("%s/%s", dir, entry.Name()))
		}
	}

	if len(files) < 1 {
		return nil, errors.New("no seed files found")
	}

	sort.Strings(files)

	db, err := sqlx.Open(config.DriverName, config.URI)
	if err != nil {
		return nil, err
	}

	return &Seed{config, db, embedFS, files}, nil
}

// Close closes the database and prevents new queries from starting.
func (s *Seed) Close() error {
	return s.db.Close()
}

// Run executes the seeding.
func (s *Seed) Run() error {
	for _, file := range s.files {
		b, err := s.embedFS.ReadFile(file)
		if err != nil {
			return err
		}

		if len(b) == 0 {
			continue
		}

		if _, err := s.db.Exec(string(b)); err != nil {
			return fmt.Errorf("unable to seed: '%s'. error: %w", file, err)
		}
	}

	return nil
}
