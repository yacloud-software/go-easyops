package mysql

// this package opens and maintains database connections
// to postgres and provide some metrics for us

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	pp "golang.conradwood.net/go-easyops/profiling"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/utils"
	"golang.org/x/net/context"
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
	dbhost          = flag.String("mysql_host", "localhost", "hostname of the postgres database rdbms")
	dbdb            = flag.String("mysql_db", "", "database to use")
	dbuser          = flag.String("mysql_user", "root", "username for the database to use")
	dbpw            = flag.String("mysql_pw", "pw", "password for the database to use")
	sqldebug        = flag.Bool("debug_mysql", false, "debug mysql stuff")
	sqlTotalQueries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mysql_queries_executed",
			Help: "V=1 UNIT=ops DESC=total number of sql queries started",
		},
		[]string{"database", "queryname"},
	)
	sqlPerformance = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.015, 0.99: 0.001},
			Name:       "mysql_query_performance",
			Help:       "V=1 UNIT=durations DESC=timing information for sql performance in seconds",
		},
		[]string{"database", "queryname"},
	)
	sqlFailedQueries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mysql_queries_failed",
			Help: "V=1 UNIT=ops DESC=total number of sql queries failed",
		},
		[]string{"database", "queryname"},
	)
	/*
		poolSize = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "sql_pool_size",
				Help: "how many connections are open",
			},
			[]string{"database"},
		)
	*/
	metricsRegistered   = false
	metricsRegisterLock sync.Mutex
	databases           []*DB
	opendblock          sync.Mutex
)

type DB struct {
	dbcon           *sql.DB
	dbname          string
	dbinfo          string
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

	var err error
	var now string
	if *dbdb == "" {
		return nil, fmt.Errorf("Please specify -mysql_db flag")
	}
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", *dbuser, *dbpw, *dbhost, *dbdb)

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
			prometheus.MustRegister(sqlPerformance, sqlTotalQueries, sqlFailedQueries)
			metricsRegistered = true
		}
		metricsRegisterLock.Unlock()
	}

	dbcon, err := sql.Open("mysql", dbinfo)
	if err != nil {
		fmt.Printf("Failed to connect to %s on host \"%s\" as \"%s\"\n", *dbdb, *dbhost, *dbuser)
		return nil, err
	}
	dbcon.SetMaxIdleConns(maxIdle())
	dbcon.SetMaxOpenConns(maxConnections()) // max connections per instance by default
	dbcon.SetConnMaxLifetime(time.Second * time.Duration(90))

	// force at least one connection to initialize
	err = dbcon.QueryRow("SELECT NOW() as now").Scan(&now)
	if err != nil {
		fmt.Printf("Failed to query db %s: %s\n", *dbdb, err)
		return nil, err
	}
	c := &DB{dbcon: dbcon, dbname: *dbdb, dbinfo: dbinfo, MaxQueryTimeout: DEFAULT_MAX_QUERY_MILLIS}
	databases = append(databases, c)
	if len(databases) > 2 {
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
	l := prometheus.Labels{"database": d.dbname, "queryname": name}
	sqlTotalQueries.With(l).Inc()
	started := time.Now()
	r, err := d.dbcon.QueryContext(ctx, query, args...)
	sqlPerformance.With(l).Observe(time.Since(started).Seconds())
	// return err if occured, or context-error if such occured
	if err == nil && ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		if *sqldebug {
			fmt.Printf("[sql] Query %s failed (%s)\n", query, err)
		}
		sqlFailedQueries.With(l).Inc()
	}
	return r, err
}

// "name" will be used to provide timing information as prometheus metric.
func (d *DB) ExecContext(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error) {
	pp.SqlEntered()
	defer pp.SqlDone()
	l := prometheus.Labels{"database": d.dbname, "queryname": name}
	if *sqldebug {
		fmt.Printf("[sql] Exec %s (%v)\n", query, args)
	}
	sqlTotalQueries.With(l).Inc()
	started := time.Now()
	r, err := d.dbcon.ExecContext(ctx, query, args...)
	sqlPerformance.With(l).Observe(time.Since(started).Seconds())
	// return err if occured, or context-error if such occured
	if err == nil && ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		if *sqldebug {
			fmt.Printf("[sql] Query %s failed (%s)\n", query, err)
		}
		sqlFailedQueries.With(l).Inc()
	}
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
	l := prometheus.Labels{"database": d.dbname, "queryname": name}
	sqlTotalQueries.With(l).Inc()
	return d.dbcon.QueryRowContext(ctx, query, args...)
}
