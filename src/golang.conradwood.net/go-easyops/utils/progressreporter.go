package utils

import (
	"fmt"
	"sync"
	"time"
)

/*
* This is a simple command line "progress reporter".
* It's usage is (intentionally) very simple.
* Example:
* pr := utils.ProgressReporter{}
* pr.SetTotal(1000)
* for i:=0;i<1000;i++ {
*    // do something slow
*    DoSomethingSlow()
*    pr.Inc()
*    pr.Print()
* }
* The above sippet will print a rate and an ETA once a second.
 */
type ProgressReporter struct {
	lastPrinted time.Time
	eta         time.Time
	etaCalced   time.Time
	rate1       *RateCalculator
	rate2       *RateCalculator
	cur_rc      int       // 1==rate1, 2==rate2
	since_rc    time.Time // when was cur_rc last changed ?
	start       time.Time
	total       uint64
	done        uint64
	RawPrint    bool
	addlock     sync.Mutex
	Prefix      string
}

func (p *ProgressReporter) SetTotal(total uint64) {
	p.total = total
}
func (p *ProgressReporter) Set(a uint64) {
	p.Add(a - p.done)
}
func (p *ProgressReporter) Inc() {
	p.Add(1)
}
func (p *ProgressReporter) Add(a uint64) {
	p.fixRates()
	p.addlock.Lock()
	defer p.addlock.Unlock()
	p.rate1.Add(a)
	p.rate2.Add(a)
	p.done = p.done + a
}
func (p *ProgressReporter) Eta() time.Time {
	if time.Since(p.etaCalced) < (time.Duration(5) * time.Second) {
		return p.eta
	}
	left := float64(p.total - p.done)
	r := p.Rate()
	secs_to_go := time.Duration(left/r) * time.Second
	p.eta = time.Now().Add(secs_to_go)
	p.etaCalced = time.Now()
	return p.eta
}

// return true if it actually printed stuff
func (p *ProgressReporter) PrintSingleLine() bool {
	s := p.String()
	if s == "" {
		return false
	}
	fmt.Printf("%c%s", byte(13), s)
	return true
}
func (p *ProgressReporter) Print() bool {
	s := p.String()
	if s == "" {
		return false
	}
	fmt.Println(s)
	return true
}

func (p *ProgressReporter) fixRates() {
	if p.rate1 == nil {
		p.rate1 = &RateCalculator{name: "calc1", start: time.Now()}
	}
	if p.rate2 == nil {
		p.rate2 = &RateCalculator{name: "calc2", start: time.Now()}
	}
	if p.cur_rc == 0 {
		p.cur_rc = 1
	}
	if time.Since(p.since_rc) > time.Duration(5)*time.Second {
		if p.cur_rc == 1 {
			p.cur_rc = 2
			p.rate1.Reset()
		} else {
			p.cur_rc = 1
			p.rate2.Reset()
		}
		p.since_rc = time.Now()
	}
}
func (p *ProgressReporter) Rate() float64 {
	p.fixRates()

	var rc *RateCalculator

	rc = p.rate1
	if p.cur_rc == 2 {
		rc = p.rate2
	}
	//	fmt.Printf("Rate from %s\n", rc.String())
	res := rc.Rate()
	return res
}
func (p *ProgressReporter) String() string {
	if (time.Since(p.lastPrinted)) < (time.Duration(1) * time.Second) {
		return ""
	}
	p.lastPrinted = time.Now()
	eta_s := p.Eta().Format("2006-01-02 15:04:05")
	perc := float32(float32(p.done) / float32(p.total) * float32(100))
	sp := ""
	if p.Prefix != "" {
		sp = fmt.Sprintf("[%s]: ", p.Prefix)
	}
	prefix := fmt.Sprintf("%sProcessing %d", sp, p.done)
	if p.total != 0 {
		prefix = fmt.Sprintf("%sProcessing %d of %d (%2.1f%%), ETA: %v", sp, p.done, p.total, perc, eta_s)
	}
	if p.RawPrint {
		return prefix + fmt.Sprintf(", %.1f/sec", p.Rate())
	} else {
		return prefix + fmt.Sprintf(", %s/sec", PrettyNumber(uint64(p.Rate())))
	}

}

type RateCalculator struct {
	start     time.Time
	counter   uint64
	additions int
	resetted  bool
	name      string
}

func (r *RateCalculator) Add(a uint64) {
	r.additions++
	r.counter = r.counter + a

}
func (r *RateCalculator) Rate() float64 {
	elapsed := time.Since(r.start).Seconds()
	z := float64(r.counter)
	f := z / (elapsed)
	return f
}
func (r *RateCalculator) Reset() {
	r.start = time.Now()
	r.additions = 0
	r.counter = 0
	r.resetted = false
	//	fmt.Printf("Reset \"%s\"\n", r.name)
}

func (r *RateCalculator) String() string {
	return fmt.Sprintf("%s (started %0.1f seconds ago, points=%d", r.name, time.Since(r.start).Seconds(), r.additions)
}
