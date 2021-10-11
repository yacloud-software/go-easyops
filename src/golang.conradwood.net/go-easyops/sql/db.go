package sql

// this package opens and maintains database connections
// to postgres and provide some metrics for us

import (
	"database/sql"
	"flag"
	"fmt"
	pq "github.com/lib/pq"
	"golang.conradwood.net/go-easyops/cmdline"
	pp "golang.conradwood.net/go-easyops/profiling"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/utils"
	"golang.org/x/net/context"
	"strings"
	"sync"
	"time"
)

const (
	DEFAULT_MAX_QUERY_MILLIS = 3000
)

var (
	/* eventually we'll look these up in the datacenter rather than passing
	these as command line parameters.
	this will increase security a little bit (at least obscure it a bit)
	-- Database URL vs Command line parameters: --
	we are not using a DB Url here because the syntax of the url is driver/vendor specific.
	The abstraction into these variables puts the burden of generating a valid url into the code
	rather than requiring the user to know the syntax of the url of the specific driver/version/vendor
	the binary was compiled with.
	*/
	f_dbhost     = flag.String("dbhost", "localhost", "hostname of the postgres database rdbms")
	f_dbdb       = flag.String("dbdb", "", "database to use")
	f_dbuser     = flag.String("dbuser", "root", "username for the database to use")
	f_dbpw       = flag.String("dbpw", "pw", "password for the database to use")
	sqldebug     = flag.Bool("ge_debug_sql", false, "debug sql stuff")
	sqlTotalTime = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sql_total_time",
			Help: "V=1 UNIT=durationms DESC=total time spent in a query",
		},
		[]string{"dbhost", "database", "queryname"},
	)
	sqlTotalQueries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sql_queries_executed",
			Help: "V=1 UNIT=ops DESC=total number of sql queries started",
		},
		[]string{"dbhost", "database", "queryname"},
	)
	sqlPerformance = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.015, 0.99: 0.001},
			Name:       "sql_query_performance",
			Help:       "V=1 UNIT=durations DESC=timing information for sql performance in seconds",
		},
		[]string{"dbhost", "database", "queryname"},
	)
	sqlFailedQueries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sql_queries_failed",
			Help: "V=1 UNIT=ops DESC=total number of sql queries failed",
		},
		[]string{"dbhost", "database", "queryname"},
	)
	metricsRegistered   = false
	metricsRegisterLock sync.Mutex
	databases           []*DB
	opendblock          sync.Mutex
)

type DB struct {
	dbcon           *sql.DB
	dbname          string
	dbinfo          string
	dbhost          string // hostname as specified on commandline
	dbshorthost     string // hostname only (no domain)
	MaxQueryTimeout int
}

func maxConnections() int {
	return 5
}
func maxIdle() int {
	return 4
}

// call this once when you startup and cache the result
// only if there is an error you'll need to retry
func Open() (*DB, error) {
	return OpenWithInfo(cmdline.OptEnvString(*f_dbhost, "GE_DBHOST"),
		cmdline.OptEnvString(*f_dbdb, "GE_DBDB"),
		cmdline.OptEnvString(*f_dbuser, "GE_DBUSER"),
		cmdline.OptEnvString(*f_dbpw, "GE_DBPW"),
	)
}
func OpenWithInfo(dbhost, dbdb, dbuser, dbpw string) (*DB, error) {
	var err error
	var now string
	if dbdb == "" {
		return nil, fmt.Errorf("Please specify -dbdb flag")
	}
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", dbhost, dbuser, dbpw, dbdb)

	// check if we already have an sql object that matches, if so return it
	for _, db := range databases {
		if db.dbinfo == dbinfo {
			return db, nil
		}
	}
	opendblock.Lock()
	defer opendblock.Unlock()
	// check again, with lock
	for _, db := range databases {
		if db.dbinfo == dbinfo {
			return db, nil
		}
	}

	if !metricsRegistered {
		metricsRegisterLock.Lock()
		if !metricsRegistered {
			prometheus.MustRegister(sqlTotalTime, sqlPerformance, sqlTotalQueries, sqlFailedQueries, NewPoolSizeCollector())
			metricsRegistered = true
		}
		metricsRegisterLock.Unlock()
	}

	dbcon, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Printf("Failed to connect to %s on host \"%s\" as \"%s\"\n", dbdb, dbhost, dbuser)
		return nil, err
	}
	dbcon.SetMaxIdleConns(maxIdle())
	dbcon.SetMaxOpenConns(maxConnections()) // max connections per instance by default
	// force at least one connection to initialize
	err = dbcon.QueryRow("SELECT NOW() as now").Scan(&now)
	if err != nil {
		fmt.Printf("Failed to query db %s: %s\n", dbdb, err)
		return nil, err
	}
	names := strings.Split(dbhost, ".")
	dbshort := dbhost
	if len(names) > 0 {
		dbshort = names[0]
	}
	c := &DB{dbcon: dbcon, dbname: dbdb, dbinfo: dbinfo, MaxQueryTimeout: DEFAULT_MAX_QUERY_MILLIS, dbhost: dbhost, dbshorthost: dbshort}
	databases = append(databases, c)
	if len(databases) > 5 {
		fmt.Printf("[go-easyops] WARNING OPENED %d databases\n", len(databases))
		for i, d := range databases {
			fmt.Printf("Opened database #%d: %s\n", i, d.dbinfo)
		}
		panic("too many databases")
	}
	return c, nil
}

/*****
// Helpers
/**********/
// returns true if this string is sql safe (no special characters
func IsSQLSafe(txt string) bool {
	return utils.IsOnlyChars(txt, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}
func (d *DB) GetDatabaseName() string {
	return d.dbname
}

/*****
// wrapping the calls
/**********/

// "name" will be used to provide timing information as prometheus metric.
func (d *DB) QueryContext(ctx context.Context, name string, query string, args ...interface{}) (*sql.Rows, error) {
	pp.SqlEntered()
	defer pp.SqlDone()
	if *sqldebug {
		fmt.Printf("[sql] Query %s (%v)\n", query, args)
	}
	l := prometheus.Labels{"dbhost": d.dbshorthost, "database": d.dbname, "queryname": name}
	sqlTotalQueries.With(l).Inc()
	started := time.Now()
	r, err := d.dbcon.QueryContext(ctx, query, args...)
	duration := time.Since(started).Seconds()
	sqlPerformance.With(l).Observe(duration)
	// return err if occured, or context-error if such occured
	if err == nil && ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		if *sqldebug {
			fmt.Printf("[sql] Query %s (%s) failed (%s)\n", name, query, err)
		}
		sqlFailedQueries.With(l).Inc()
	}
	sqlTotalTime.With(l).Add(duration)
	pqtime(name, duration)
	return r, err
}

// "name" will be used to provide timing information as prometheus metric.
func (d *DB) ExecContext(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error) {
	pp.SqlEntered()
	defer pp.SqlDone()
	l := prometheus.Labels{"dbhost": d.dbshorthost, "database": d.dbname, "queryname": name}
	if *sqldebug {
		fmt.Printf("[sql] Exec %s (%v)\n", query, args)
	}
	sqlTotalQueries.With(l).Inc()
	started := time.Now()
	r, err := d.dbcon.ExecContext(ctx, query, args...)
	duration := time.Since(started).Seconds()
	sqlPerformance.With(l).Observe(duration)
	// return err if occured, or context-error if such occured
	if err == nil && ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		if *sqldebug {
			fmt.Printf("[sql] Query %s (%s) failed (%s)\n", name, query, err)
		}
		sqlFailedQueries.With(l).Inc()
	}
	sqlTotalTime.With(l).Add(duration)
	pqtime(name, duration)
	return r, err
}

// discouraged use. QueryRow() does not provide an error on the query, nor do we get a good timing
// value. Use QueryContext() instead.
func (d *DB) QueryRowContext(ctx context.Context, name string, query string, args ...interface{}) *sql.Row {
	pp.SqlEntered()
	defer pp.SqlDone()
	if *sqldebug {
		fmt.Printf("[sql] QueryRow %s\n", query)
	}
	l := prometheus.Labels{"dbhost": d.dbshorthost, "database": d.dbname, "queryname": name}
	sqlTotalQueries.With(l).Inc()
	r := d.dbcon.QueryRowContext(ctx, query, args...)
	return r
}

func (d *DB) Conn(ctx context.Context) (*sql.Conn, error) {
	return d.dbcon.Conn(ctx)
}

func (d *DB) CheckDuplicateRowError(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" {
			return true
		}
	}

	return false
}

func pqtime(name string, dur float64) {
	if !*sqldebug {
		return
	}
	fmt.Printf("Query \"%s\" completed in %0.2f seconds\n", name, dur)
}
