package server

import (
	"context"
	"flag"
	"fmt"
	el "golang.conradwood.net/apis/errorlogger"
	fw "golang.conradwood.net/apis/framework"
	ge "golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc/status"
	"time"
)

var (
	logChan       = make(chan *le, 200)
	els           el.ErrorLoggerClient
	error_looping = false
	debug_elog    = flag.Bool("ge_debug_error_log", false, "if true debug what is being sent to the error logger")
)

type le struct {
	ts  time.Time
	sd  *serverDef
	rc  *rpccall
	ctx context.Context
	err error
}

func (sd *serverDef) logError(ctx context.Context, rc *rpccall, err error) {
	if cmdline.IsStandalone() {
		fmt.Printf("[go-easyops] ERROR: %s\n", err)
	}
	if len(logChan) > 100 {
		fmt.Printf("[go-easyops] Dropping errorlog\n")
		return
	}
	l := &le{sd: sd, ctx: ctx, rc: rc, err: err, ts: time.Now()}
	logChan <- l
}
func error_handler_startup() {
	if cmdline.IsStandalone() {
		return
	}
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
	u := auth.GetUser(l.ctx)
	uid := ""
	if u != nil {
		uid = u.ID
	}
	svc := auth.GetService(l.ctx)
	st := status.Convert(l.err)
	e := &el.ErrorLogRequest{
		UserID:         uid,
		ErrorCode:      uint32(st.Code()),
		ErrorMessage:   fmt.Sprintf("%s", l.err),
		LogMessage:     utils.ErrorString(l.err),
		ServiceName:    l.rc.ServiceName,
		MethodName:     l.rc.MethodName,
		Timestamp:      uint32(l.ts.Unix()),
		RequestID:      "norequestidinerrorhandler",
		CallingService: svc,
		Errors:         &ge.GRPCErrorList{},
	}
	/*
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
	*/
	for _, a := range st.Details() {
		if a == nil {
			continue
		}
		fmd, ok := a.(*ge.GRPCError)
		if !ok {
			continue
		}
		e.Errors.Errors = append(e.Errors.Errors, fmd)
	}
	ctx := authremote.Context()
	if *debug_elog {
		fmt.Printf("[go-easyops] errorlog: %v\n", e)
	}
	els.Log(ctx, e)
}
func AddErrorDetail(st *status.Status, ct *ge.GRPCError) *status.Status {
	// add details (and keep previous)
	odet := st.Details()
	if *debug_rpc_serve {
		fancyPrintf("Error %s (%s) (%s)\n", st.Err(), st.Message(), utils.ErrorString(st.Err()))
	}
	add := &ge.GRPCErrorList{}
	for _, d := range odet {
		if *debug_rpc_serve {
			fancyPrintf("keeping error %v\n", d)
		}
		fmd, ok := d.(*ge.GRPCError)
		if ok {
			add.Errors = append(add.Errors, fmd)
		} else {
			add.Errors = append(add.Errors, &ge.GRPCError{LogMessage: fmt.Sprintf("%v", d)})

		}
	}
	add.Errors = append(add.Errors, ct)
	stn, errx := st.WithDetails(add)

	// if adding details failed, just return the undecorated error message
	if errx != nil {
		if *debug_rpc_serve {
			fancyPrintf("failed to get status with detail: %s", errx)
		}
		return st
	}
	return stn
}
func AddStatusDetail(st *status.Status, ct *fw.CallTrace) *status.Status {
	return st
	/*
		// add details (and keep previous)
		add := &fw.FrameworkMessageDetail{Message: ct.Message}
		odet := st.Details()
		if *debug_rpc_serve {
			fancyPrintf("Error %s (%s) (%s)\n", st.Err(), st.Message(), utils.ErrorString(st.Err()))
		}
		for _, d := range odet {
			if *debug_rpc_serve {
				fancyPrintf("keeping error %v\n", d)
			}
			fmd, ok := d.(*fw.FrameworkMessageDetail)
			if ok {
				add.CallTraces = append(add.CallTraces, fmd.CallTraces...)
			} else {
				add.CallTraces = append(add.CallTraces, &fw.CallTrace{Message: fmt.Sprintf("%v", d)})

			}
		}
		add.CallTraces = append(add.CallTraces, ct)
		stn, errx := st.WithDetails(add)

		// if adding details failed, just return the undecorated error message
		if errx != nil {
			if *debug_rpc_serve {
				fancyPrintf("failed to get status with detail: %s", errx)
			}
			return st
		}
		return stn
	*/
}
