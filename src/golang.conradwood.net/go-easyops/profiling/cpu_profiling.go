package profiling

import (
	"bytes"
	"flag"
	"fmt"
	_ "net/http/pprof"
	_ "os"
	"runtime/pprof"
	"time"
)

var (
	cpuprofile     = flag.Bool("ge_cpuprofile", false, "enable profiling")
	profiling      = false
	buf            bytes.Buffer
	profileChannel = make(chan *profileInfo, 10)
	chan_started   = false
)

type profileInfo struct {
	mtype   int // 0->start/stop
	started bool
}

func GetBuf() *bytes.Buffer {
	return &buf
}

func Toggle() {
	*cpuprofile = !*cpuprofile
	ProfilingCheckStart()
}

func IsActive() bool {
	return profiling
}

// called on startup and if cpuprofile flag changes (see server.go)
func ProfilingCheckStart() {
	if !chan_started {
		go profiler_watcher()
		chan_started = true
	}
	if profiling == *cpuprofile {
		return
	}
	if profiling {
		ProfilingStop()
		return
	}
	if *cpuprofile {
		/*
			f, err := os.Create("cpuprofile")
			if err != nil {
				fmt.Printf("[go-easyops] cpuprofile: %s", err)
				return
			}
		*/
		fmt.Printf("Starting CPU Profiling...\n")
		buf.Reset()
		pprof.StartCPUProfile(&buf)
		profiling = true
		profileChannel <- &profileInfo{mtype: 0, started: true}
	}
}

// called when this application shuts down. at most once.
func ProfilingStop() {
	if profiling {
		fmt.Printf("Stopping CPU Profiling...\n")
		profileChannel <- &profileInfo{mtype: 0, started: false}
		pprof.StopCPUProfile()
		profiling = false
	}
}

func profiler_watcher() {
	for {
		c := <-profileChannel
		fmt.Printf("Notification: %v\n", c)
		for profiling {
			if buf.Len() > (1024 * 1024 * 10) {
				fmt.Printf("%v Bufsize: %d\n", profiling, buf.Len())
				ProfilingStop()
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}
