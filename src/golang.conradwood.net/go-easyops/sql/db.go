/*
Package sql provides safe and managed access to (postgres) databases
*/
package sql

// this package opens and maintains database connections
// to postgres and provide some metrics for us

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	pq "github.com/lib/pq"
	pb "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/errors"
	pp "golang.conradwood.net/go-easyops/profiling"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/go-easyops/utils"
)

const (
	DEFAULT_MAX_QUERY_MILLIS = 3000
)

var (
	e_dbhost  = cmdline.ENV("GE_POSTGRES_HOST", "the postgresql hostname to connect to")
	e_dbdb    = cmdline.ENV("GE_POSTGRES_DB", "the postgresql database to connect to")
	e_dbuser  = cmdline.ENV("GE_POSTGRES_USER", "the postgresql user to connect with")
	e_dbpw    = cmdline.ENV("GE_POSTGRES_PW", "the postgresql password to connect with")
	e_dbproto = cmdline.ENV("GE_POSTGRES_PROTO", "the postgresql connection details as base64 encoded proto goeasyops.PostgresConfig")

	failure_action = flag.String("ge_sql_failure_action", "report", "one of [report|quit|retry], report means to report it to the application (this is the default), quit means to quit the process, retry means to block until the connection is open")
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
	print_errors = flag.Bool("ge_print_sql_errors", true, "print all sql errors (all failed queries)")
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
	failurectr      *utils.SlidingAverage
	lastReconnect   time.Time
	reconnectLock   sync.RWMutex
}

func maxConnections() int {
	return 5
}
func maxIdle() int {
	return 4
}

// make sure we catch missing database configuration issues (fail-fast)
// that is stupid, because it makes it fail even if we include it but not execute it (import with _)
func init() {
	/*
		go func() {
			time.Sleep(time.Duration(3) * time.Second)
			_, err := Open()
			if err != nil {
				fmt.Printf("Application %s error\n", cmdline.SourceCodePath())
			}
			utils.Bail("database error", err)
		}()
	*/
}

// call this once when you startup and cache the result
// only if there is an error you'll need to retry
func Open() (*DB, error) {
	host := *f_dbhost
	db := *f_dbdb
	user := *f_dbuser
	pw := *f_dbpw
	if host == "" {
		host = e_dbhost.Value()
	}
	if db == "" {
		db = e_dbdb.Value()
	}
	if user == "" {
		user = e_dbuser.Value()
	}
	if pw == "" {
		pw = e_dbpw.Value()
	}
	if e_dbproto.Value() != "" {
		pp := &pb.PostgresConfig{}
		err := utils.Unmarshal(e_dbproto.Value(), pp)
		utils.Bail("invalid configuration in "+e_dbproto.Name(), err)
		if host == "" {
			host = pp.Host
		}
		if db == "" {
			db = pp.DB
		}
		if user == "" {
			user = pp.User
		}
		if pw == "" {
			pw = pp.PW
		}
	}
	return OpenWithInfo(host, db, user, pw)
}
func OpenWithInfo(dbhost, dbdb, dbuser, dbpw string) (*DB, error) {
	var err error
	var now string
	if dbdb == "" {
		return nil, errors.Errorf("Please specify -dbdb flag")
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

	var dbcon *sql.DB
	for {
		dbcon, err = sql.Open("postgres", dbinfo)
		if err == nil {
			break
		}
		fmt.Printf("[go-easyops] Failed to connect to %s on host \"%s\" as \"%s\"\n", dbdb, dbhost, dbuser)
		if *failure_action == "quit" {
			os.Exit(10)
		} else if *failure_action == "report" {
			return nil, errors.Errorf("failed to open database \"%s\" on host \"%s\" as \"%s\": %w", dbdb, dbhost, dbuser, err)
		} else if *failure_action == "retry" {
			time.Sleep(time.Duration(2) * time.Second)
		} else {
			fmt.Printf("[go-easyops] ge_sql_failure_action must be one of [report|retry|quit], not \"%s\"\n", *failure_action)
			os.Exit(10)
		}
	}

	dbcon.SetMaxIdleConns(maxIdle())
	dbcon.SetMaxOpenConns(maxConnections()) // max connections per instance by default
	// force at least one connection to initialize
	err = dbcon.QueryRow("SELECT NOW() as now").Scan(&now)
	if err != nil {
		fmt.Printf("[go-easyops] Failed to query db %s: %s\n", dbdb, err)
		return nil, errors.Errorf("failed to open database \"%s\" on host \"%s\" as \"%s\": %w", dbdb, dbhost, dbuser, err)
	}
	if *sqldebug {
		fmt.Printf("[go-easyops] sql now query returned: \"%s\"\n", now)
	}
	names := strings.Split(dbhost, ".")
	dbshort := dbhost
	if len(names) > 0 {
		dbshort = names[0]
	}
	c := &DB{dbcon: dbcon, dbname: dbdb, dbinfo: dbinfo, MaxQueryTimeout: DEFAULT_MAX_QUERY_MILLIS, dbhost: dbhost, dbshorthost: dbshort}
	c.failurectr = utils.NewSlidingAverage()
	c.failurectr.MinSamples = 10
	c.failurectr.MinAge = time.Duration(60) * time.Second
	databases = append(databases, c)
	if len(databases) > 5 {
		fmt.Printf("[go-easyops] WARNING OPENED %d databases\n", len(databases))
		for i, d := range databases {
			fmt.Printf("[go-easyops] Opened database #%d: %s\n", i, d.dbinfo)
		}
		panic("too many databases")
	}
	return c, nil
}
func reopen() {
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
func (d *DB) GetFailureCounter() *utils.SlidingAverage {
	return d.failurectr
}

func (d *DB) reconnect_if_required() {
	if !must_reconnect(d.GetFailureCounter()) {
		return
	}
	if time.Since(d.lastReconnect) < time.Duration(60)*time.Second {
		// recently reconnected, ignore..
		return
	}
	if *sqldebug {
		sa := d.GetFailureCounter()
		fmt.Printf("[go-easyops] sql counters: 0=%d, 1=%d\n", sa.GetCounter(0), sa.GetCounter(1))
	}
	fmt.Printf("[go-easyops] sql reconnect required.\n")
	d.reconnectLock.Lock()
	defer d.reconnectLock.Unlock()
	d.lastReconnect = time.Now()
	fmt.Printf("[go-easyops] sql reconnecting...\n")
	d.dbcon.Close()
	dc, err := sql.Open("postgres", d.dbinfo)
	if err != nil {
		fmt.Printf("[go-easyops] sql failed to reconnect: %s\n", err)
		return
	}
	d.dbcon = dc

}

func must_reconnect(sa *utils.SlidingAverage) bool {
	if sa.GetCounter(0) != 0 {
		// at least one succeeded
		return false
	}
	if sa.GetCounter(1) == 0 {
		// no failures
		return false
	}
	if sa.GetCounter(1) != sa.GetCounts(1) {
		// some failure counts where not "1"?? so we counted a failure as 0 or 2 failures?
		return false
	}
	return true
}

func query_error(ctx context.Context, typ string, name string, query string, err error) {
	if *sqldebug || *print_errors {
		fmt.Printf("[go-easyops] [sql] %s %s (%s) failed (%s)\n", typ, name, query, err)
	}
	if *sqldebug {
		utils.PrintStack("query failed")
	}
}

/*****
// wrapping the calls
/**********/

// "name" will be used to provide timing information as prometheus metric.
func (d *DB) QueryContext(ctx context.Context, name string, query string, args ...interface{}) (*sql.Rows, error) {
	d.reconnect_if_required()
	d.reconnectLock.RLock()
	defer d.reconnectLock.RUnlock()
	pp.SqlEntered()
	defer pp.SqlDone()
	if *sqldebug {
		fmt.Printf("[go-easyops] [sql] Query %s (%v)\n", query, args)
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
		d.failurectr.Add(1, 1)
		query_error(ctx, "select", name, query, err)
		sqlFailedQueries.With(l).Inc()
		err = errors.Wrap(err)
	} else {
		d.failurectr.Add(0, 1)
	}
	sqlTotalTime.With(l).Add(duration)
	pqtime(name, duration)
	return r, err
}

// "name" will be used to provide timing information as prometheus metric.
func (d *DB) ExecContext(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error) {
	rep := *sqldebug || *print_errors
	return d.execContext(ctx, rep, name, query, args...)
}
func (d *DB) ExecContextQuiet(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error) {
	rep := false
	return d.execContext(ctx, rep, name, query, args...)
}
func (d *DB) execContext(ctx context.Context, report_failure bool, name string, query string, args ...interface{}) (sql.Result, error) {
	d.reconnect_if_required()
	d.reconnectLock.RLock()
	defer d.reconnectLock.RUnlock()
	pp.SqlEntered()
	defer pp.SqlDone()
	l := prometheus.Labels{"dbhost": d.dbshorthost, "database": d.dbname, "queryname": name}
	if *sqldebug {
		fmt.Printf("[go-easyops] [sql] Exec %s (%v)\n", query, args)
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
		d.failurectr.Add(1, 1)
		query_error(ctx, "exec", name, query, err)
		sqlFailedQueries.With(l).Inc()
		err = errors.Wrap(err)
	} else {
		d.failurectr.Add(0, 1)
	}
	sqlTotalTime.With(l).Add(duration)
	pqtime(name, duration)
	return r, err
}

// discouraged use. QueryRow() does not provide an error on the query, nor do we get a good timing
// value. Use QueryContext() instead.
func (d *DB) QueryRowContext(ctx context.Context, name string, query string, args ...interface{}) *sql.Row {
	d.reconnect_if_required()
	d.reconnectLock.RLock()
	defer d.reconnectLock.RUnlock()
	pp.SqlEntered()
	defer pp.SqlDone()
	if *sqldebug {
		fmt.Printf("[go-easyops] [sql] QueryRow %s\n", query)
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
	fmt.Printf("[go-easyops] Query \"%s\" completed in %0.2f seconds\n", name, dur)
}
