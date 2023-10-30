//go:generate mockery --all --outpkg=mock --output=./mock
package sqlx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/jmoiron/sqlx"
)

// DBer is the interface that implements all the *sqlx.DB functions.
type DBer interface {
	// Begin starts a transaction. The default isolation level is dependent on the driver.
	//
	// Begin uses context.Background internally; to specify the context, use BeginTx.
	Begin() (*sql.Tx, error)

	// BeginTx starts a transaction.
	//
	// The provided context is used until the transaction is committed or rolled back. If the context
	// is canceled, the sql package will roll back the transaction. Tx.Commit will return an error if
	// the context provided to BeginTx is canceled.
	//
	// The provided TxOptions is optional and may be nil if defaults should be used. If a non-default
	// isolation level is used that the driver doesn't support, an error will be returned.
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

	// BeginTxx begins a transaction and returns an *sqlx.Tx instead of an *sql.Tx.
	//
	// The provided context is used until the transaction is committed or rolled back. If the context
	// is canceled, the sql package will roll back the transaction. Tx.Commit will return an error if
	// the context provided to BeginxContext is canceled.
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)

	// Beginx begins a transaction and returns an *sqlx.Tx instead of an *sql.Tx.
	Beginx() (*sqlx.Tx, error)

	// BindNamed binds a query using the DB driver's bindvar type.
	BindNamed(query string, arg interface{}) (string, []interface{}, error)

	// Close closes the database and prevents new queries from starting. Close then waits for all
	// queries that have started processing on the server to finish.
	//
	// It is rare to Close a DB, as the DB handle is meant to be long-lived and shared between many
	// goroutines.
	Close() error

	// Conn returns a single connection by either opening a new connection or returning an existing
	// connection from the connection pool. Conn will block until either a connection is returned or
	// ctx is canceled. Queries run on the same Conn will be run in the same database session.
	//
	// Every Conn must be returned to the database pool after use by calling Conn.Close.
	Conn(ctx context.Context) (*sql.Conn, error)

	// Connx returns an *sqlx.Conn instead of an *sql.Conn.
	Connx(ctx context.Context) (*sqlx.Conn, error)

	// Driver returns the database's underlying driver.
	Driver() driver.Driver

	// DriverName returns the driverName passed to the Open function for this DB.
	DriverName() string

	// Exec executes a query without returning any rows. The args are for any placeholder parameters
	// in the query.
	//
	// Exec uses context.Background internally; to specify the context, use ExecContext.
	Exec(query string, args ...interface{}) (sql.Result, error)

	// ExecContext executes a query without returning any rows. The args are for any placeholder
	// parameters in the query.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Get using this DB. Any placeholder parameters are replaced with supplied args. An error is
	// returned if the result set is empty.
	Get(dest interface{}, query string, args ...interface{}) error

	// GetContext using this DB. Any placeholder parameters are replaced with supplied args. An error
	// is returned if the result set is empty.
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// MapperFunc sets a new mapper for this db using the default sqlx struct tag and the provided
	// mapper function.
	MapperFunc(mf func(string) string)

	// NamedExec using this DB. Any named placeholder parameters are replaced with fields from arg.
	NamedExec(query string, arg interface{}) (sql.Result, error)

	// NamedExecContext using this DB. Any named placeholder parameters are replaced with fields from
	// arg.
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)

	// NamedQuery using this DB. Any named placeholder parameters are replaced with fields from arg.
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)

	// NamedQueryContext using this DB. Any named placeholder parameters are replaced with fields
	// from arg.
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)

	// Ping verifies a connection to the database is still alive, establishing a connection if
	// necessary.
	//
	// Ping uses context.Background internally; to specify the context, use PingContext.
	Ping() error

	// PingContext verifies a connection to the database is still alive, establishing a connection if
	// necessary.
	PingContext(ctx context.Context) error

	// Prepare creates a prepared statement for later queries or executions. Multiple queries or executions may be run concurrently from the returned statement. The caller must call the statement's Close method when the statement is no longer needed.
	//
	// Prepare uses context.Background internally; to specify the context, use PrepareContext.
	Prepare(query string) (*sql.Stmt, error)

	// PrepareContext creates a prepared statement for later queries or executions. Multiple queries or executions may be run concurrently from the returned statement. The caller must call the statement's Close method when the statement is no longer needed.
	//
	// The provided context is used for the preparation of the statement, not for the execution of the statement.
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	// PrepareNamed returns an sqlx.NamedStmt.
	PrepareNamed(query string) (*sqlx.NamedStmt, error)

	// PrepareNamedContext returns an sqlx.NamedStmt.
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)

	// Preparex returns an sqlx.Stmt instead of a sql.Stmt.
	Preparex(query string) (*sqlx.Stmt, error)

	// PreparexContext returns an sqlx.Stmt instead of a sql.Stmt.
	//
	// The provided context is used for the preparation of the statement, not for the execution of
	// the statement.
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)

	// Query executes a query that returns rows, typically a SELECT. The args are for any placeholder
	// parameters in the query.
	//
	// Query uses context.Background internally; to specify the context, use QueryContext.
	Query(query string, args ...interface{}) (*sql.Rows, error)

	// QueryContext executes a query that returns rows, typically a SELECT. The args are for any
	// placeholder parameters in the query.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// QueryRow executes a query that is expected to return at most one row. QueryRow always returns
	// a non-nil value. Errors are deferred until Row's Scan method is called. If the query selects
	// no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's Scan scans the first
	// selected row and discards the rest.
	//
	// QueryRow uses context.Background internally; to specify the context, use QueryRowContext.
	QueryRow(query string, args ...interface{}) *sql.Row

	// QueryRowContext executes a query that is expected to return at most one row. QueryRowContext
	// always returns a non-nil value. Errors are deferred until Row's Scan method is called. If the
	// query selects no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's Scan scans
	// the first selected row and discards the rest.
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	// QueryRowx queries the database and returns an *sqlx.Row. Any placeholder parameters are
	// replaced with supplied args.
	QueryRowx(query string, args ...interface{}) *sqlx.Row

	// QueryRowxContext queries the database and returns an *sqlx.Row. Any placeholder parameters are
	// replaced with supplied args.
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row

	// Queryx queries the database and returns an *sqlx.Rows. Any placeholder parameters are replaced
	// with supplied args.
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)

	// QueryxContext queries the database and returns an *sqlx.Rows. Any placeholder parameters are
	// replaced with supplied args.
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)

	// Rebind transforms a query from QUESTION to the DB driver's bindvar type.
	Rebind(query string) string

	// Select using this DB. Any placeholder parameters are replaced with supplied args.
	Select(dest interface{}, query string, args ...interface{}) error

	// SelectContext using this DB. Any placeholder parameters are replaced with supplied args.
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// SetConnMaxIdleTime sets the maximum amount of time a connection may be idle.
	//
	// Expired connections may be closed lazily before reuse.
	//
	// If d <= 0, connections are not closed due to a connection's idle time.
	SetConnMaxIdleTime(d time.Duration)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	//
	// Expired connections may be closed lazily before reuse.
	//
	// If d <= 0, connections are not closed due to a connection's age.
	SetConnMaxLifetime(d time.Duration)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	//
	// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns, then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
	//
	// If n <= 0, no idle connections are retained.
	//
	// The default max idle connections is currently 2. This may change in a future release.
	SetMaxIdleConns(n int)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	//
	// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than MaxIdleConns, then MaxIdleConns will be reduced to match the new MaxOpenConns limit.
	//
	// If n <= 0, then there is no limit on the number of open connections. The default is 0 (unlimited).
	SetMaxOpenConns(n int)

	// Stats returns database statistics.
	Stats() sql.DBStats

	// Unsafe returns a version of DB which will silently succeed to scan when columns in the SQL
	// result have no fields in the destination struct. sqlx.Stmt and sqlx.Tx which are created
	// from this DB will inherit its safety behavior.
	Unsafe() *sqlx.DB
}
