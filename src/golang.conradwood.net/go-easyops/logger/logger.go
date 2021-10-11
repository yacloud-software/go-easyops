package logger

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/logservice"
	"golang.conradwood.net/go-easyops/client"
	"golang.org/x/net/context"
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
	line    string
	status  string
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

func NewAsyncLogQueue(appname string, repoid uint64, group, namespace, deplid string) (*AsyncLogQueue, error) {
	if appname == "" {
		return nil, fmt.Errorf("[go-easyops] Will not instantiate an AsyncLogQueue without appname")
	}
	repo := fmt.Sprintf("FIX_ME_AUTODEPLOYER_LOGGER_%d", repoid)
	if repo == "" {
		return nil, fmt.Errorf("[go-easyops] Will not instantiate an AsyncLogQueue without repo")
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
			Repository:   repo,
			RepoID:       repoid,
			Groupname:    group,
			Namespace:    namespace,
			DeploymentID: deplid,
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

func (alq *AsyncLogQueue) LogCommandStdout(line string, status string) error {
	if *log_debug {
		fmt.Printf("app:\"%s\" LOGGED: %s\n", alq.appDef.Appname, line)
	}
	qe := QueueEntry{
		created: time.Now().Unix(),
		line:    line,
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
func (alq *AsyncLogQueue) Log(status string, format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	alq.LogCommandStdout(status, s)
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
	_, e := grpcClient.CloseLog(context.Background(), &cl)
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

	flushies := alq.entries
	alq.entries = &([]*QueueEntry{})

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
				Time:   qe.created,
				Line:   qe.line,
				Status: qe.status,
			},
		)
	}

	_, err := grpcClient.LogCommandStdout(context.Background(), logRequest)
	if err != nil {
		if time.Since(alq.lastErrPrinted) > (10 * time.Second) {
			fmt.Printf("Failed to send log: %s\n", err)
			alq.lastErrPrinted = time.Now()
		}
	}

	// all done, so clear the array so we free up the memory

	return nil
}
