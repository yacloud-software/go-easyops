package logger

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/logservice"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/ctx"
	"sync"
	"time"
)

var (
	log_debug  = flag.Bool("logger_debug", false, "set to true to debug logging")
	grpcClient logservice.LogServiceClient
	inp        = false
	logLock    sync.Mutex
)

type QueueEntry struct {
	created int64
	status  string
	binline []byte
}

type AsyncLogQueue struct {
	closed         bool
	appDef         *logservice.LogAppDef
	entries        *[]*QueueEntry
	lastErrPrinted time.Time
	MaxSize        int
	sync.Mutex
}

func getClient() error {
	if inp {
		return fmt.Errorf("Logservice initialisation already in progress\n")
	}
	if grpcClient != nil {
		return nil
	}
	logLock.Lock()
	inp = true
	grpcClient = logservice.NewLogServiceClient(client.Connect("logservice.LogService"))
	inp = false
	logLock.Unlock()
	return nil
}

func NewAsyncLogQueue(appname string, buildid, repoid uint64, group, namespace, deplid string) (*AsyncLogQueue, error) {
	if appname == "" {
		return nil, fmt.Errorf("[go-easyops] Will not instantiate an AsyncLogQueue without appname")
	}
	if group == "" {
		return nil, fmt.Errorf("[go-easyops] Will not instantiate an AsyncLogQueue without group ")
	}
	if namespace == "" {
		return nil, fmt.Errorf("[go-easyops] Will not instantiate an AsyncLogQueue without namespace")
	}
	if deplid == "" {
		return nil, fmt.Errorf("[go-easyops] Will not instantiate an AsyncLogQueue without deploymentid")
	}
	alq := &AsyncLogQueue{
		appDef: &logservice.LogAppDef{
			Appname:      appname,
			RepoID:       repoid,
			Groupname:    group,
			Namespace:    namespace,
			DeploymentID: deplid,
			BuildID:      buildid,
		},
		closed:  false,
		MaxSize: 5000,
	}

	alq.entries = &([]*QueueEntry{})
	t := time.NewTicker(1 * time.Second)

	go func(a *AsyncLogQueue) {
		for _ = range t.C {
			err := a.Flush()
			if (*log_debug) && (err != nil) {
				fmt.Printf("Error flushing logqueue:%s\n", err)
			}
		}
	}(alq)

	return alq, nil
}

func (alq *AsyncLogQueue) String() string {
	if alq == nil {
		return "empty_asynclogqueue"
	}
	ad := alq.appDef
	if ad == nil {
		return "new_asynclogqueue"
	}
	return fmt.Sprintf("Log for %s, repoid %d, build %d (deplid: %s)", ad.Appname, ad.RepoID, ad.BuildID, ad.DeploymentID)
}
func (alq *AsyncLogQueue) LogCommandStdout(line string, status string) error {
	if *log_debug {
		fmt.Printf("app:\"%s\" LOGGED: %s\n", alq.appDef.Appname, line)
	}
	qe := QueueEntry{
		created: time.Now().Unix(),
		binline: []byte(line),
		status:  status,
	}
	alq.Lock()
	defer alq.Unlock()
	if len(*alq.entries) > alq.MaxSize {
		if *log_debug {
			fmt.Printf("queue size larger than %d (it is %d) - discarding log entries\n", alq.MaxSize, len(*alq.entries))
		}
		alq.entries = &([]*QueueEntry{})
	}

	*alq.entries = append(*alq.entries, &qe)

	return nil
}
func (alq *AsyncLogQueue) Write(status string, buf []byte) {
	qe := QueueEntry{
		created: time.Now().Unix(),
		binline: buf,
		status:  status,
	}
	alq.Lock()
	defer alq.Unlock()
	if len(*alq.entries) > alq.MaxSize {
		if *log_debug {
			fmt.Printf("queue size larger than %d (it is %d) - discarding log entries\n", alq.MaxSize, len(*alq.entries))
		}
		alq.entries = &([]*QueueEntry{})
	}

	*alq.entries = append(*alq.entries, &qe)
}
func (alq *AsyncLogQueue) Log(status string, format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	alq.LogCommandStdout(s, status)
}

func (alq *AsyncLogQueue) Close(exitcode int) error {
	err := alq.Flush()
	if alq.closed {
		return fmt.Errorf("Closed already!")
	}
	cl := logservice.CloseLogRequest{AppDef: alq.appDef, ExitCode: int32(exitcode)}
	lerr := getClient()
	if lerr != nil {
		return lerr
	}
	_, e := grpcClient.CloseLog(getctx(), &cl)
	if e != nil {
		return e
	}
	alq.closed = true
	return err
}
func (alq *AsyncLogQueue) SetStartupID(s string) {
	if alq.appDef == nil {
		return
	}
	alq.appDef.StartupID = s
}
func (alq *AsyncLogQueue) Flush() error {
	lerr := getClient()
	if lerr != nil {
		return lerr
	}

	// all done, so clear the array so we free up the memory
	alq.Lock()
	flushies := alq.entries
	alq.entries = &([]*QueueEntry{})
	alq.Unlock()

	if len(*flushies) == 0 {
		// save ourselves from dialing and stuff
		return nil
	}

	logRequest := &logservice.LogRequest{
		AppDef: alq.appDef,
	}

	for _, qe := range *flushies {
		logRequest.Lines = append(
			logRequest.Lines,
			&logservice.LogLine{
				Time:    qe.created,
				BinLine: qe.binline,
				Status:  qe.status,
			},
		)
	}

	_, err := grpcClient.LogCommandStdout(getctx(), logRequest)
	if err != nil {
		if time.Since(alq.lastErrPrinted) > (10 * time.Second) {
			fmt.Printf("%s: Failed to send log: %s\n", alq.String(), err)
			alq.lastErrPrinted = time.Now()

			// try and stick something into the logserver (unlikely to work, unless a logline causes trouble)
			olines := logRequest.Lines
			lc := 0
			bc := 0
			for _, ol := range olines {
				lc++
				bc = bc + len(ol.Line) + len(ol.BinLine)
			}
			logRequest.Lines = []*logservice.LogLine{
				&logservice.LogLine{
					Time:    time.Now().Unix(),
					Status:  "LOGFAILURE",
					BinLine: []byte(fmt.Sprintf("failed to log %d lines and %d bytes: %s", lc, bc, err)),
				},
			}
			grpcClient.LogCommandStdout(getctx(), logRequest)

		}
	}

	return nil
}

func getctx() context.Context {
	cb := ctx.NewContextBuilder()
	return cb.ContextWithAutoCancel()
}
