package sql

import "time"

var (
	DriverType = struct {
		MYSQL    string
		POSTGRES string
	}{
		MYSQL:    "mysql",
		POSTGRES: "postgres",
	}

	embedInit = "init.sql"
)

// Config is the database connection configuration.
type Config struct {
	// ConnMaxLifetime indicates the maximum amount of time a connection may be reused. Expired
	// connections may be closed lazily before reuse.
	//
	// If d <= 0, connections are reused forever. By default, it is 5 * time.Minute.
	ConnMaxLifetime time.Duration

	// ConnMaxIdleTime indicates the maximum amount of time a connection may be idle. Expired
	// connections may be closed lazily before reuse.
	//
	// If d <= 0, connections are not closed due to a connection's idle time. By default, it is
	// 5 * time.Minute.
	ConnMaxIdleTime time.Duration

	// DriverName indicates the SQL driver to use, currently only supports:
	// 	- mysql
	// 	- postgres
	DriverName string

	// MaxIdleConns indicates the maximum number of connections in the idle connection pool. By
	// default, it is 16.
	//
	// Note: MaxIdleConns will automatically be updated to use the same value as MaxOpenConns if
	// MaxIdleConns is greater than MaxOpenConns.
	MaxIdleConns int

	// MaxOpenConns indicates the maximum number of open connections to the database. By default,
	// it is 16.
	MaxOpenConns int

	// Name is an alias to uniquely identify the database for commands:
	// 	- db:migrate
	// 	- db:migrate:status
	// 	- db:rollback
	// 	- db:seed
	//	- gen:migration
	//
	// Note: This is also used to decide how the database migrate/seed path should look like. For
	// example, if the value is 'payment', the DB commands will refer to:
	//	- db/migrate/payment/*.sql for DB migrations
	//	- db/seed/payment/*.sql for DB seeds
	Name string

	// SchemaSearchPath indicates the schema search path which is only used with "postgres". By
	// default, it is "public".
	SchemaSearchPath string

	// SchemaMigrationsTable indicates the table name for storing the schema migration versions.
	// By default, it is "schema_migrations".
	SchemaMigrationsTable string

	// URI indicates the master database connection string to connect.
	//
	// URI connection string documentation:
	//   - mysql: https://dev.mysql.com/doc/refman/8.0/en/connecting-using-uri-or-key-value-pairs.html#connecting-using-uri
	//   - postgres: https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
	URI string

	// ReplicaURIs indicates the replica database connections string to connect.
	ReplicaURIs []string
}

func defaultConfig(c *Config) *Config {
	if c.ConnMaxIdleTime == 0 {
		c.ConnMaxIdleTime = 5 * time.Minute
	}

	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = 5 * time.Minute
	}

	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 16
	}

	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 16
	}

	if c.MaxIdleConns > c.MaxOpenConns {
		c.MaxIdleConns = c.MaxOpenConns
	}

	if c.SchemaSearchPath == "" {
		c.SchemaSearchPath = "public"
	}

	if c.SchemaMigrationsTable == "" {
		c.SchemaMigrationsTable = "schema_migrations"
	}

	return c
}
