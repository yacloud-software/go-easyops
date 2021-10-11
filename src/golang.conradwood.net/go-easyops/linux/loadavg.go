package linux

import (
	"bufio"
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/prometheus"
	"golang.conradwood.net/apis/common"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	READ_MILLIS = 5000
)

var (
	loadGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "goeasyops_cpu_percent",
			Help: "ticks spent per type in a cpu",
		},
		[]string{"type"},
	)
	ctr           = 0
	proclock      sync.Mutex
	procstat      = make([]uint64, 10)
	procstat_diff = make([]uint64, 10)
	total         uint64
)

func init() {
	prometheus.MustRegister(loadGauge)
	go load_calc_loop()
}

type Loadavg struct {
}

func (g *Loadavg) GetCPULoad(ctx context.Context, req *common.Void) (*common.CPULoad, error) {
	if ctr < 3 {
		return nil, fmt.Errorf("not ready yet - initialising. try later")
	}
	res := &common.CPULoad{}
	file, err := os.Open("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var line string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		break
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	avgs := strings.Split(line, " ")
	if len(avgs) < 3 {
		return nil, fmt.Errorf("Invalid line read from loadavg: \"%s\" (splits into %d parts)", line, len(avgs))
	}
	res.Avg1, err = strconv.ParseFloat(avgs[0], 32)
	if err != nil {
		return nil, err
	}

	res.Avg5, err = strconv.ParseFloat(avgs[1], 32)
	if err != nil {
		return nil, err
	}

	res.Avg15, err = strconv.ParseFloat(avgs[2], 32)
	if err != nil {
		return nil, err
	}

	// get number of cpus
	res.CPUCount = uint32(runtime.NumCPU())
	res.PerCPU = res.Avg1 / float64(res.CPUCount)
	proclock.Lock()
	res.Sum = 0
	for _, r := range procstat_diff {
		res.Sum = res.Sum + r
	}
	res.RawSum = total
	res.User = procstat_diff[0]
	res.Nice = procstat_diff[1]
	res.System = procstat_diff[2]
	res.Idle = procstat_diff[3]
	res.IOWait = procstat_diff[4]
	res.IRQ = procstat_diff[5]
	res.SoftIRQ = procstat_diff[6]
	proclock.Unlock()
	res.IdleTime = float64(res.Idle) / float64(res.Sum) * 100
	return res, nil
}

func load_calc_loop() {
	def := time.Duration(READ_MILLIS) * time.Millisecond
	for {
		if ctr < 3 { // fast startup
			time.Sleep(time.Duration(100) * time.Millisecond)
			ctr++
		} else {
			time.Sleep(def)
		}
		err := load_calc()
		if err != nil {
			fmt.Printf("Failed to process /proc/stat: %s\n", err)
			continue
		}
	}
}

// we read the usage regularly and compare with previous reading so to get result
func load_calc() error {
	// scan short-term utilisation

	file, err := os.Open("/proc/stat")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		break
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	ns := strings.Fields(line)
	if len(ns) < 8 {
		return fmt.Errorf("invalid /proc/stat: %s (%d nums)\n", line, len(ns))
	}
	var nums []uint64
	for i, n := range ns {
		if i == 0 {
			continue
		}
		num, err := strconv.ParseUint(n, 10, 64)
		if err != nil {
			return err
		}
		nums = append(nums, num)
	}
	proclock.Lock()
	total = 0
	for i, n := range nums {
		total = total + n
		procstat_diff[i] = n - procstat[i]
		procstat[i] = n
	}
	proclock.Unlock()
	dt := uint64(0)
	for _, r := range procstat_diff {
		dt = dt + r
	}
	gauge("idle", float64(procstat_diff[3])/float64(dt)*100)
	gauge("busy", 100.0-(float64(procstat_diff[3])/float64(dt)*100))
	return nil
}

func gauge(name string, perc float64) {
	loadGauge.With(prometheus.Labels{"type": name}).Set(perc)
}
