//go:generate mockery --all --outpkg=mock --output=./mock
package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	sdksqlx "ecommerce/utils/sqlx"

	logger "ecommerce/utils/logger"

	"github.com/XSAM/otelsql"
	"github.com/jmoiron/sqlx"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zapadapter"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
)

// DB is a wrapper around sdksqlx.DBer which keeps track of the driverName upon Open, used mostly to
// automatically bind named queries using the right bindvars.
type DB struct {
	master   sdksqlx.DBer
	replicas []sdksqlx.DBer
	config   *Config
	logger   *logger.Logger

	// currReplicaIdx is helper variable round robin slave.
	currReplicaIdx uint64
}

// Conn is a *sql.Conn wrapper for all the master/replicas connections.
type Conn struct {
	Master  *sql.Conn
	Replica []*sql.Conn
}

// Driver is a wrapper for driver.Driver which wraps all the master/replicas drivers.
type Driver struct {
	Master  driver.Driver
	Replica []driver.Driver
}

// DBStats is a sql.DBStats wrapper for all the master/replicas connections.
type DBStats struct {
	Master  sql.DBStats
	Replica []sql.DBStats
}

// DBer is the interface that implements all the DB functions.
type DBer interface {
	// Begin starts a transaction. The default isolation level is dependent on the driver.
	Begin() (*sqlx.Tx, error)

	// BeginContext starts a transaction.
	//
	// The provided context is used until the transaction is committed or rolled back. If the context
	// is canceled, the sql package will roll back the transaction. Tx.Commit will return an error if
	// the context provided to BeginContext is canceled.
	//
	// The provided TxOptions is optional and may be nil if defaults should be used. If a non-default
	// isolation level is used that the driver doesn't support, an error will be returned.
	BeginContext(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)

	// Close returns the connection to the connection pool.
	// All operations after a Close will return with ErrConnDone.
	// Close is safe to call concurrently with other operations and will
	// block until all other operations finish. It may be useful to first
	// cancel any used context and then call close directly after.
	Close() error

	// Config returns the database handle's config.
	Config() *Config

	// Conn returns the wrapper for all the master/replicas connections.
	Conn(ctx context.Context) (Conn, error)

	// Connect establishes a connection to the database specified in URI and assign
	// the database handle which is safe for concurrent use by multiple goroutines
	// and maintains its own connection pool.
	Connect() error

	// Driver returns the database's underlying driver.
	Driver() Driver

	// DriverName returns the database handle's SQL driver name.
	DriverName() string

	// Exec executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	//
	// Exec uses context.Background internally; to specify the context, use
	// ExecContext.
	Exec(query string, args ...interface{}) (sql.Result, error)

	// ExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Get does a QueryRow using the provided Queryer, and scans the resulting row
	// to dest.  If dest is scannable, the result must only have one column.  Otherwise,
	// StructScan is used.  Get will return sql.ErrNoRows like row.Scan would.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	Get(dest interface{}, query string, args ...interface{}) error

	// GetContext using this DB.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// NamedExec using this DB.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedExec(query string, arg interface{}) (sql.Result, error)

	// NamedExecContext using this DB.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)

	// NamedQuery using this DB. Any named placeholder parameters are replaced with fields from arg.
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)

	// NamedQueryContext using this DB. Any named placeholder parameters are replaced with fields from arg.
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)

	// Open opens a database specified by its database driver name and a
	// driver-specific data source name, usually consisting of at least a
	// database name and connection information.
	//
	// Most users will open a database via a driver-specific connection
	// helper function that returns db *DB. No database drivers are included
	// in the Go standard library. See https://golang.org/s/sqldrivers for
	// a list of third-party drivers.
	//
	// Open may just validate its arguments without creating a connection
	// to the database. To verify that the data source name is valid, call
	// Ping.
	//
	// The returned DB is safe for concurrent use by multiple goroutines
	// and maintains its own pool of idle connections. Thus, the Open
	// function should be called just once. It is rarely necessary to
	// close a DB.
	Open() error

	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	//
	// Ping uses context.Background internally; to specify the context, use
	// PingContext.
	Ping() error

	// Prepare creates a prepared statement for later queries or executions. Multiple queries or
	// executions may be run concurrently from the returned statement. The caller must call the
	// statement's Close method when the statement is no longer needed.
	Prepare(query string) (*sqlx.Stmt, error)

	// PrepareContext creates a prepared statement for later queries or executions. Multiple queries
	// or executions may be run concurrently from the returned statement. The caller must call the
	// statement's Close method when the statement is no longer needed.
	//
	// The provided context is used for the preparation of the statement, not for the execution of
	// the statement.
	PrepareContext(ctx context.Context, query string) (*sqlx.Stmt, error)

	// PrepareNamed returns an *sqlx.NamedStmt.
	PrepareNamed(query string) (*sqlx.NamedStmt, error)

	// PrepareNamedContext returns an *sqlx.NamedStmt.
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)

	// Query executes a query that returns rows, typically a SELECT. The args are for any placeholder
	// parameters in the query.
	Query(query string, args ...interface{}) (*sqlx.Rows, error)

	// QueryContext executes a query that returns rows, typically a SELECT. The args are for any
	// placeholder parameters in the query.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)

	// QueryRow executes a query that is expected to return at most one row. QueryRow always returns a
	// non-nil value. Errors are deferred until Row's Scan method is called.
	//
	// If the query selects no rows, the Row's Scan will return ErrNoRows. Otherwise, the Row's Scan
	// scans the first selected row and discards the rest.
	QueryRow(query string, args ...interface{}) *sqlx.Row

	// QueryRowContext executes a query that is expected to return at most one row. QueryRowContext
	// always returns a non-nil value. Errors are deferred until Row's Scan method is called.
	//
	// If the query selects no rows, the Row's Scan will return ErrNoRows. Otherwise, the Row's Scan
	// scans the first selected row and discards the rest.
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row

	// Rebind transforms a query from QUESTION to the DB driver's bindvar type.
	Rebind(query string) string

	// Select using this DB. Any placeholder parameters are replaced with supplied args.
	Select(dest interface{}, query string, args ...interface{}) error

	// SelectContext using this DB. Any placeholder parameters are replaced with supplied args.
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	//
	// Expired connections may be closed lazily before reuse.
	//
	// If d <= 0, connections are not closed due to a connection's age.
	SetConnMaxLifetime(d time.Duration)

	// SetConnMaxIdleTime sets the maximum amount of time a connection may be idle.
	//
	// Expired connections may be closed lazily before reuse.
	//
	// If d <= 0, connections are not closed due to a connection's idle time.
	SetConnMaxIdleTime(d time.Duration)

	// SetMaxIdleConns sets the maximum number of connections in the idle
	// connection pool.
	//
	// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
	// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
	//
	// If n <= 0, no idle connections are retained.
	SetMaxIdleConns(n int)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	//
	// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
	// MaxIdleConns, then MaxIdleConns will be reduced to match the new
	// MaxOpenConns limit.
	//
	// If n <= 0, then there is no limit on the number of open connections.
	SetMaxOpenConns(n int)

	// Master returns the master database.
	Master() sdksqlx.DBer

	// Replica returns 1 of the replica databases. It returns the master database instead if the
	// replica list is empty.
	Replica() sdksqlx.DBer

	// Stats returns database statistics for all the master/replicas connections.
	Stats() DBStats
}

var (
	logLevel = sqldblogger.WithMinimumLevel(sqldblogger.LevelDebug)
)

func init() {
	if os.Getenv("APP_ENV") != "" && os.Getenv("APP_ENV") != "development" {
		logLevel = sqldblogger.WithMinimumLevel(sqldblogger.LevelInfo)
	}
}

// NewDB initialises the DB handle that is used to connect to a database.
func NewDB(c *Config, logger *logger.Logger) DBer {
	return &DB{
		config:         defaultConfig(c),
		currReplicaIdx: 0,
		logger:         logger,
	}
}

// Begin starts a transaction on master DB. The default isolation level is dependent on the driver.
func (db *DB) Begin() (*sqlx.Tx, error) {
	return db.master.Beginx()
}

// BeginContext starts a transaction.
//
// The provided context is used until the transaction is committed or rolled back. If the context
// is canceled, the sql package will roll back the transaction. Tx.Commit will return an error if
// the context provided to BeginContext is canceled.
//
// The provided TxOptions is optional and may be nil if defaults should be used. If a non-default
// isolation level is used that the driver doesn't support, an error will be returned.
func (db *DB) BeginContext(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return db.master.BeginTxx(ctx, opts)
}

// Close returns the connection to the connection pool.
func (db *DB) Close() error {
	err := db.master.Close()
	if err != nil {
		return err
	}

	for i := range db.replicas {
		err := db.replicas[i].Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// Config returns the database config.
func (db *DB) Config() *Config {
	return db.config
}

// Conn return the database connection
func (db *DB) Conn(ctx context.Context) (Conn, error) {
	var (
		result Conn
		err    error
	)

	// get master connections
	result.Master, err = db.master.Conn(ctx)
	if err != nil {
		return result, err
	}

	// get replica connections
	for i := range db.replicas {
		replicaConn, err := db.replicas[i].Conn(ctx)
		if err != nil {
			return result, err
		}

		result.Replica = append(result.Replica, replicaConn)
	}

	return result, nil
}

// Connect establishes a connection to the database and verify with a ping.
func (db *DB) Connect() error {
	dbSystem := ""
	switch db.Config().DriverName {
	case "mysql":
		dbSystem = semconv.DBSystemMySQL.Value.AsString()
	case "postgres":
		dbSystem = semconv.DBSystemPostgreSQL.Value.AsString()
	}

	driverName, err := otelsql.Register(db.Config().DriverName, dbSystem)
	if err != nil {
		return err
	}

	masterDB, err := sqlx.Connect(driverName, db.Config().URI)
	if err != nil {
		return fmt.Errorf("unable to connect to '%s', error: %s", db.Config().URI, err.Error())
	}

	wrappedMasterDB, err := newDBWithLogger(masterDB.DriverName(), db.Config().URI, db.Config(), masterDB.Driver(), db.logger.Desugar())
	if err != nil {
		return err
	}
	db.master = wrappedMasterDB

	for _, replicaURI := range db.Config().ReplicaURIs {
		replicaDB, err := sqlx.Connect(driverName, replicaURI)
		if err != nil {
			return fmt.Errorf("unable to connect to '%s', error: %s", replicaURI, err.Error())
		}

		wrappedReplicaDB, err := newDBWithLogger(replicaDB.DriverName(), db.Config().URI, db.Config(), replicaDB.Driver(), db.logger.Desugar())
		if err != nil {
			return err
		}
		db.replicas = append(db.replicas, wrappedReplicaDB)
	}

	return nil
}

// Driver returns the database's underlying driver.
func (db *DB) Driver() Driver {
	var result Driver

	// get master connections
	result.Master = db.master.Driver()

	// get replica connections
	for i := range db.replicas {
		replicaConn := db.replicas[i].Driver()
		result.Replica = append(result.Replica, replicaConn)
	}

	return result
}

// DriverName returns the database handle's SQL driver name.
func (db *DB) DriverName() string {
	// get master connections
	result := db.master.DriverName()

	return result
}

// Exec executes a query using the master database without returning any rows. The args are for
// any placeholder parameters in the query.
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.master.Exec(query, args...)
}

// ExecContext executes a query using the master database without returning any rows. The args are
// for any placeholder parameters in the query.
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.master.ExecContext(ctx, query, args...)
}

// Get using 1 of the replica databases. Any placeholder parameters are replaced with supplied
// args. An error is returned if the result set is empty.
func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	return db.Replica().Get(dest, query, args...)
}

// GetContext using 1 of the replica databases. Any placeholder parameters are replaced with
// supplied args. An error is returned if the result set is empty.
func (db *DB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.Replica().GetContext(ctx, dest, query, args...)
}

// NamedExec using the master database. Any named placeholder parameters are replaced with fields
// from arg.
func (db *DB) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return db.master.NamedExec(query, arg)
}

// NamedExecContext using the master database. Any named placeholder parameters are replaced with
// fields from arg.
func (db *DB) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return db.master.NamedExecContext(ctx, query, arg)
}

// NamedQuery using 1 of the replica databases. Any named placeholder parameters are replaced with
// fields from arg.
func (db *DB) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return db.Replica().NamedQuery(query, arg)
}

// NamedQueryContext using 1 of the replica databases. Any named placeholder parameters are
// replaced with fields from arg.
func (db *DB) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	return db.Replica().NamedQueryContext(ctx, query, arg)
}

// Open opens a database specified by its database driver name and a
// driver-specific data source name, usually consisting of at least a
// database name and connection information.
//
// Most users will open a database via a driver-specific connection
// helper function that returns db *DB. No database drivers are included
// in the Go standard library. See https://golang.org/s/sqldrivers for
// a list of third-party drivers.
//
// Open may just validate its arguments without creating a connection
// to the database. To verify that the data source name is valid, call
// Ping.
//
// The returned DB is safe for concurrent use by multiple goroutines
// and maintains its own pool of idle connections. Thus, the Open
// function should be called just once. It is rarely necessary to
// close a DB.
func (db *DB) Open() error {
	dbSystem := ""
	switch db.Config().DriverName {
	case "mysql":
		dbSystem = semconv.DBSystemMySQL.Value.AsString()
	case "postgres":
		dbSystem = semconv.DBSystemPostgreSQL.Value.AsString()
	}

	driverName, err := otelsql.Register(db.Config().DriverName, dbSystem)
	if err != nil {
		return err
	}

	masterDB, err := sqlx.Open(driverName, db.Config().URI)
	if err != nil {
		return fmt.Errorf("unable to connect to '%s', error: %s", db.Config().URI, err.Error())
	}

	wrappedMasterDB, err := newDBWithLogger(masterDB.DriverName(), db.Config().URI, db.Config(), masterDB.Driver(), db.logger.Desugar())
	if err != nil {
		return err
	}
	db.master = wrappedMasterDB

	for _, replicaURI := range db.Config().ReplicaURIs {
		replicaDB, err := sqlx.Open(driverName, replicaURI)
		if err != nil {
			return fmt.Errorf("unable to connect to '%s', error: %s", replicaURI, err.Error())
		}

		wrappedReplicaDB, err := newDBWithLogger(replicaDB.DriverName(), db.Config().URI, db.Config(), replicaDB.Driver(), db.logger.Desugar())
		if err != nil {
			return err
		}
		db.replicas = append(db.replicas, wrappedReplicaDB)
	}

	return nil
}

// Prepare creates a prepared statement for later queries or executions. Multiple queries or
// executions may be run concurrently from the returned statement. The caller must call the
// statement's Close method when the statement is no longer needed.
func (db *DB) Prepare(query string) (*sqlx.Stmt, error) {
	return db.master.Preparex(query)
}

// PrepareContext creates a prepared statement for later queries or executions. Multiple queries
// or executions may be run concurrently from the returned statement. The caller must call the
// statement's Close method when the statement is no longer needed.
//
// The provided context is used for the preparation of the statement, not for the execution of
// the statement.
func (db *DB) PrepareContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return db.master.PreparexContext(ctx, query)
}

// PrepareNamed returns an *sqlx.NamedStmt.
func (db *DB) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	return db.master.PrepareNamed(query)
}

// PrepareNamedContext returns an *sqlx.NamedStmt.
func (db *DB) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	return db.master.PrepareNamedContext(ctx, query)
}

// Query executes a query using 1 of the replica databases that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (db *DB) Query(query string, args ...interface{}) (*sqlx.Rows, error) {
	return db.Replica().Queryx(query, args...)
}

// QueryContext executes a query using 1 of the replica databases that returns rows, typically a
// SELECT. The args are for any placeholder parameters in the query.
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return db.Replica().QueryxContext(ctx, query, args...)
}

// QueryRow executes a query using 1 of the replica databases that is expected to return at most
// one row. QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method
// is called.
//
// If the query selects no rows, the Row's Scan will return ErrNoRows. Otherwise, the *sqlx.Row's
// Scan scans the first selected row and discards the rest.
func (db *DB) QueryRow(query string, args ...interface{}) *sqlx.Row {
	return db.Replica().QueryRowx(query, args...)
}

// QueryRowContext executes a query using 1 of the replica databases that is expected to return
// at most one row. QueryRowContext always returns a non-nil value. Errors are deferred until
// *sqlx.Row's Scan method is called.
//
// If the query selects no rows, the *sqlx.Row's Scan will return ErrNoRows. Otherwise, the
// *sqlx.Row's Scan scans the first selected row and discards the rest.
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return db.Replica().QueryRowxContext(ctx, query, args...)
}

// Ping verifies the connections to the master/replica databases are still alive, establishing a
// connection if necessary.
//
// Ping uses context.Background internally; to specify the context, use PingContext.
func (db *DB) Ping() error {
	err := db.master.Ping()
	if err != nil {
		return fmt.Errorf("unable to ping '%s', error: %s", db.config.URI, err.Error())
	}

	for i := range db.replicas {
		err := db.replicas[i].Ping()
		if err != nil {
			return fmt.Errorf("unable to ping '%s', error: %s", db.config.ReplicaURIs[i], err.Error())
		}
	}

	return nil
}

// Rebind transforms a query from QUESTION to the DB driver's bindvar type using 1 of the replica
// databases.
func (db *DB) Rebind(query string) string {
	return db.master.Rebind(query)
}

// Select using 1 of the replica databases. Any placeholder parameters are replaced with supplied
// args.
func (db *DB) Select(dest interface{}, query string, args ...interface{}) error {
	return db.Replica().Select(dest, query, args...)
}

// Select using 1 of the replica databases. Any placeholder parameters are replaced with supplied
// args.
func (db *DB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.Replica().SelectContext(ctx, dest, query, args...)
}

// SetMaxIdleConns sets the maximum number of master/replica connections in their idle connection
// pools.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns, then the new MaxIdleConns
// will be reduced to match the MaxOpenConns limit.
//
// If n <= 0, no idle connections are retained.
func (db *DB) SetMaxIdleConns(max int) {
	db.master.SetMaxIdleConns(max)

	for i := range db.replicas {
		db.replicas[i].SetMaxIdleConns(max)
	}
}

// SetMaxOpenConns sets the maximum number of open connections to the master/replica databases.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than MaxIdleConns, then
// MaxIdleConns will be reduced to match the new MaxOpenConns limit.
//
// If n <= 0, then there is no limit on the number of open connections.
func (db *DB) SetMaxOpenConns(n int) {
	// set max conn idle for master
	db.master.SetMaxOpenConns(n)

	// set max conn idle conn for replica
	for i := range db.replicas {
		db.replicas[i].SetMaxOpenConns(n)
	}
}

// SetConnMaxLifetime sets the maximum amount of time a connection of the master/replica databases
// may be reused.
//
// Expired connections may be closed lazily before reuse.
//
// If d <= 0, connections are reused forever.
func (db *DB) SetConnMaxLifetime(d time.Duration) {
	db.master.SetConnMaxLifetime(d)

	for i := range db.replicas {
		db.replicas[i].SetConnMaxLifetime(d)
	}
}

// SetConnMaxIdleTime sets the maximum amount of time a connection of the master/replica databases
// may be idle.
//
// Expired connections may be closed lazily before reuse.
//
// If d <= 0, connections are not closed due to a connection's idle time.
func (db *DB) SetConnMaxIdleTime(d time.Duration) {
	db.master.SetConnMaxIdleTime(d)

	for i := range db.replicas {
		db.replicas[i].SetConnMaxIdleTime(d)
	}
}

// Stats returns the database statistics for all the master/replica databases.
func (db *DB) Stats() DBStats {
	var (
		result DBStats
	)

	result.Master = db.master.Stats()

	for _, r := range db.replicas {
		result.Replica = append(result.Replica, r.Stats())
	}

	return result
}

// Master returns the master database.
func (db *DB) Master() sdksqlx.DBer {
	return db.master
}

// Replica returns 1 of the replica databases. It returns the master database instead if the
// replica list is empty.
func (db *DB) Replica() sdksqlx.DBer {
	if len(db.replicas) < 1 {
		db.logger.Warn("fallback to master database due to the empty replica database list")
		return db.master
	}

	return db.replicas[db.replica()]
}

// Replica returns the current replica index using round-robin algorithm.
func (db *DB) replica() int {
	if len(db.replicas) <= 1 {
		return 0
	}

	return int((1 + atomic.AddUint64(&db.currReplicaIdx, 1)%uint64(len(db.replicas)-1)))
}

func newDBWithLogger(driverName, uri string, config *Config, driver driver.Driver, logger *zap.Logger) (*sqlx.DB, error) {
	dbWithLogger := sqlx.NewDb(
		sqldblogger.OpenDriver(
			uri,
			driver,
			zapadapter.New(logger),
			sqldblogger.WithIncludeStartTime(true),
			sqldblogger.WithTimeFormat(sqldblogger.TimeFormatRFC3339),
			logLevel,
		),
		driverName,
	)
	if dbWithLogger == nil {
		return nil, fmt.Errorf("unable to create sqlx.DB with logger for '%s'", uri)
	}

	dbWithLogger.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	dbWithLogger.SetConnMaxLifetime(config.ConnMaxLifetime)
	dbWithLogger.SetMaxIdleConns(config.MaxIdleConns)
	dbWithLogger.SetMaxOpenConns(config.MaxOpenConns)
	if err := dbWithLogger.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping '%s', error: %s", uri, err.Error())
	}

	return dbWithLogger, nil
}
