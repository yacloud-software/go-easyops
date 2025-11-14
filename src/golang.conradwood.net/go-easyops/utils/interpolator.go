package utils

import (
	"fmt"
	"sort"
	"sync"
)

/*
	given some reference points will interpolate a value,

examples:
 1. reference points 5(50) and 10(100), interpolation: 7=70, 8=80,3=30 (0 is assumed), 11=100 (highest=max)
*/
type Interpolator struct {
	sync.Mutex
	referencepoints []*interpolator_referencepoint
}

type interpolator_referencepoint struct {
	number float64
	value  float64
}

func (ip *Interpolator) AddReferencePoints(points map[float64]float64) {
	for k, v := range points {
		ip.AddReferencePoint(k, v)
	}
}
func (ip *Interpolator) AddReferencePoint(number, value float64) {
	ip.Lock()
	defer ip.Unlock()
	if len(ip.referencepoints) == 0 {
		ip.referencepoints = append(ip.referencepoints, &interpolator_referencepoint{number: 0, value: 0})
	}
	ip.referencepoints = append(ip.referencepoints, &interpolator_referencepoint{number: number, value: value})
	sort.Slice(ip.referencepoints, func(i, j int) bool {
		return ip.referencepoints[i].number < ip.referencepoints[j].number
	})
	var last *interpolator_referencepoint
	for _, ir := range ip.referencepoints {
		if last == nil {
			last = ir
			continue
		}
		if ir.value < last.value {
			fmt.Printf("WARNING: interpolator values out of order: %s\n", ip.String())
		}
		last = ir
	}
}

/*
fast and quick interpolation using linear interpolation
*/
func (ip *Interpolator) LinearInterpolate(number float64) float64 {
	ir1, ir2 := ip.findReferences(number)
	if ir2 == nil {
		// max out
		//	fmt.Printf("Value %0.1f is maxed out, because highest is %s\n", number, ir1.String())
		return ir1.value
	}
	xa := ir1.number
	ya := ir1.value
	xb := ir2.number
	yb := ir2.value
	x := number
	diff := (x - xa) / (xb - xa)
	y := ya + (yb-ya)*diff
	//fmt.Printf("Value %0.1f is between %s and %s, diff=%0.1f, res=%0.1f\n", number, ir1.String(), ir2.String(), diff, y)
	return y
}
func (ip *Interpolator) findReferences(number float64) (*interpolator_referencepoint, *interpolator_referencepoint) {
	var ir1 *interpolator_referencepoint
	var ir2 *interpolator_referencepoint
	for i, ir := range ip.referencepoints {
		if number < ir.number {
			continue
		}
		ir1 = ir
		if i < len(ip.referencepoints)-1 {
			ir2 = ip.referencepoints[i+1]
		} else {
			ir2 = nil
		}
	}
	return ir1, ir2
}
func (ip *Interpolator) String() string {
	deli := ""
	s := ""
	for _, ir := range ip.referencepoints {
		s = s + deli + fmt.Sprintf("%f->%f", ir.number, ir.value)
		deli = ", "
	}
	return s
}

func (ir *interpolator_referencepoint) String() string {
	return fmt.Sprintf("%f->%f", ir.number, ir.value)
}
