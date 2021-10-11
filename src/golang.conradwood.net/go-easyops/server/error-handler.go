package server

import (
	"fmt"
	el "golang.conradwood.net/apis/errorlogger"
	fw "golang.conradwood.net/apis/framework"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/rpc"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc/status"
	"time"
)

var (
	logChan       = make(chan *le, 200)
	els           el.ErrorLoggerClient
	error_looping = false
)

type le struct {
	ts  time.Time
	sd  *serverDef
	cs  *rpc.CallState
	err error
}

func (sd *serverDef) logError(cs *rpc.CallState, err error) {
	if len(logChan) > 100 {
		fmt.Printf("[go-easyops] Dropping errorlog\n")
		return
	}
	l := &le{sd: sd, cs: cs, err: err, ts: time.Now()}
	logChan <- l
}
func error_handler_startup() {
	if error_looping {
		return
	}
	error_looping = true
	if els == nil {
		els = el.NewErrorLoggerClient(client.Connect("errorlogger.ErrorLogger"))
	}
	go logLoop()
}
func logLoop() {
	for {
		l := <-logChan
		log(l)
	}
}
func log(l *le) {
	u := auth.GetUser(l.cs.Context)
	uid := ""
	if u != nil {
		uid = u.ID
	}
	st := status.Convert(l.err)
	e := &el.ErrorLogRequest{
		UserID:       uid,
		ErrorCode:    uint32(st.Code()),
		ErrorMessage: fmt.Sprintf("%s", l.err),
		LogMessage:   utils.ErrorString(l.err),
		ServiceName:  l.cs.ServiceName,
		MethodName:   l.cs.MethodName,
		Timestamp:    uint32(l.ts.Unix()),
		RequestID:    l.cs.RequestID(),
	}
	for _, a := range st.Details() {
		if a == nil {
			continue
		}
		fmd, ok := a.(*fw.FrameworkMessageDetail)
		if !ok {
			continue
		}
		e.Messages = append(e.Messages, fmd)
	}
	ctx := tokens.ContextWithToken()
	els.Log(ctx, e)
}
