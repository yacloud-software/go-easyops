/*
	The FanoutRouter distributes work (evenly) between available instances, dynamically adjusting to instances coming and going.

fanoutrouter maintains a go-routine per instance. each go-routine listens on a channel for work, if received it calls a function with a grpcConnection as parameter. The result is passed to another channel the number of go-routines changes dynamically as instances come and go each go-routine, if it has work to do, will call a "processor" (a user defined function) with a ProcessRequest. Once the processor completed its work, the result will be send to a function (perhaps even multi-threaded!)
*/
package router

import (
	"fmt"
	"sync"
	"time"

	"golang.conradwood.net/go-easyops/authremote"
	"google.golang.org/grpc"
)

const (
	state_Starting = 1
	state_Started  = 2
	state_Stopped  = 3
)

type FanoutRouter struct {
	cm             *ConnectionManager
	requests       chan *fanout_router_process_request
	cn             func(*CompletionNotification)
	proc           func(*ProcessRequest) error
	processor_wg   *sync.WaitGroup
	stopping       bool
	cur_processors []*fanout_router_processor
	stop_lock      sync.Mutex
}
type CompletionNotification struct {
	pr  *ProcessRequest
	err error
}
type ProcessRequest struct {
	req  *fanout_router_process_request
	proc *fanout_router_processor
}

// one processor per target
type fanout_router_processor struct {
	state           int
	fr              *FanoutRouter
	target          *ConnectionTarget
	control_channel chan *fanout_router_control_request
	processed       int
}

type fanout_router_process_request struct {
	o    interface{} // whatever the user wants to process
	quit bool        // special flag to stop a go-routine from processing more
}
type fanout_router_control_request struct {
	quit bool
}

func NewFanoutRouter(cm *ConnectionManager, processor func(*ProcessRequest) error, consumer func(*CompletionNotification)) *FanoutRouter {
	res := &FanoutRouter{
		cm:           cm,
		requests:     make(chan *fanout_router_process_request, 1),
		proc:         processor,
		cn:           consumer,
		processor_wg: &sync.WaitGroup{},
	}
	go res.poll_target_list()
	return res
}
func (fr *FanoutRouter) SubmitWork(object interface{}) {
	if fr.stopping {
		fr.debugf("WARNING - submitted work to fanoutrouter, which is in the process of stopping\n")
	}
	if len(fr.cur_processors) == 0 {
		fr.debugf("WARNING - submitted work to fanoutrouter, which has no backends atm\n")
	}
	pr := &fanout_router_process_request{o: object}
	fr.requests <- pr
}

// this can take a long time, because we wait for all pending requests to finish before returning
func (fr *FanoutRouter) Stop() {
	fr.debugf("Stopping...\n")
	fr.stopping = true
	fr.stop_lock.Lock()
	defer fr.stop_lock.Unlock()
	for i := 0; i < len(fr.cur_processors); i++ {
		fr.requests <- &fanout_router_process_request{quit: true}
	}
	fr.processor_wg.Wait()
	fr.debugf("Stopped\n")
}

func (fr *FanoutRouter) poll_target_list() {
	fr.stop_lock.Lock()
	fr.debugf("starting polling...\n")
	ctx := authremote.Context()
	ct := fr.cm.GetCurrentTargets(ctx)
	fr.debugf("first polling got %d targets\n", len(ct))
	fr.compare_current_targets(ct)
	fr.stop_lock.Unlock()

	for {
		if fr.stopping {
			break
		}
		time.Sleep(time.Duration(15) * time.Second)
		fr.stop_lock.Lock()
		if fr.stopping {
			fr.stop_lock.Unlock()
			break
		}
		ctx := authremote.Context()
		fr.debugf("polling...\n")
		ct := fr.cm.GetCurrentTargets(ctx)
		fr.debugf("got %d targets\n", len(ct))
		fr.compare_current_targets(ct)
		fr.stop_lock.Unlock()

	}
}
func (fr *FanoutRouter) compare_current_targets(ct []*ConnectionTarget) {
	targets := make(map[string]*ConnectionTarget)
	for _, c := range ct {
		targets[c.Address()] = c
	}
	// find new ones to start
	for _, proc := range fr.cur_processors {
		proc_adr := proc.address()
		delete(targets, proc_adr)
	}
	//now start those in targets
	for _, v := range targets {
		fp := &fanout_router_processor{fr: fr, target: v, control_channel: make(chan *fanout_router_control_request, 10)}
		fr.start_processor(fp)
	}

	// find ones to stop
	targets = make(map[string]*ConnectionTarget)
	for _, c := range ct {
		targets[c.Address()] = c
	}
	for _, proc := range fr.cur_processors {
		_, valid := targets[proc.address()]
		if !valid {
			proc.control_channel <- &fanout_router_control_request{quit: true}
		}
	}
}
func (fr *FanoutRouter) start_processor(pr *fanout_router_processor) {
	pr.state = state_Starting
	fr.cur_processors = append(fr.cur_processors, pr)
	go pr.process_requests()
	fr.processor_wg.Add(1)
}
func (fp *fanout_router_processor) process_requests() {
	prefix := fmt.Sprintf("[%s] ", fp.address())
	fmt.Printf("%sstarted\n", prefix)
	fp.state = state_Started
	for {
		select {
		case ctrl := <-fp.control_channel:
			if ctrl.quit {
				goto out
			}
		case req := <-fp.fr.requests:
			if req.quit {
				goto out
			}
			pr := &ProcessRequest{proc: fp, req: req}
			fmt.Printf("%sprocessing...\n", prefix)
			err := fp.fr.proc(pr)
			fmt.Printf("%scomplete...\n", prefix)
			cn := &CompletionNotification{pr: pr, err: err}
			fp.processed++
			fp.fr.cn(cn)
			//
		}

	}
out:
	fmt.Printf("%sFinished (after %d requests)\n", prefix, fp.processed)
	fp.fr.processor_wg.Done()
	fp.state = state_Stopped
}
func (fp *fanout_router_processor) address() string {
	return fp.target.Address()
}
func (p *ProcessRequest) Object() interface{} {
	return p.req.o
}
func (p *ProcessRequest) GRPCConnection() *grpc.ClientConn {
	rcon, err := p.proc.target.Connection()
	if err != nil {
		fmt.Printf("Failed to get Connection: %s\n", err)
		return nil
	}
	gcon, err := rcon.GRPCConnection()
	if err != nil {
		fmt.Printf("Failed to get GRPCConnection: %s\n", err)
		return nil
	}
	return gcon

}
func (p *CompletionNotification) Error() error {
	return p.err
}
func (p *CompletionNotification) Object() interface{} {
	return p.pr.req.o
}

/**************** debugf *********************/
func (fr *FanoutRouter) debugf(format string, args ...interface{}) {
	s := fmt.Sprintf("[fanoutrouter for %s] ", fr.cm.ServiceName())
	s2 := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s", s, s2)
}
