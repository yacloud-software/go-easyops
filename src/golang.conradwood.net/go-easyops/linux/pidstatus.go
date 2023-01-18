package linux

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	INITPID        = 1
	STATUS_RUNNING = 1
	STATUS_STOPPED = 2
)

type STATUS int

type ProcessState struct {
	pid             int
	binary          string
	err             error
	parentpid       int
	direct_children []*ProcessState
}

func AllPids() ([]*ProcessState, error) {
	var res []*ProcessState
	root := PidStatus(INITPID)
	res = append(res, root)
	children, err := root.recursivelyGetChildrenOf()
	if err != nil {
		return nil, err
	}
	res = append(res, children...)
	return res, nil
}
func (ps *ProcessState) getChildrenOf() ([]*ProcessState, error) {
	pid := ps.Pid()
	var res []*ProcessState
	uts, err := ioutil.ReadDir(fmt.Sprintf("/proc/%d/task", pid))
	if err != nil {
		return res, nil
	}
	var tids []int
	for _, dir := range uts {
		xpid, err := strconv.Atoi(dir.Name())
		if err != nil {
			continue
		}
		cname := fmt.Sprintf("/proc/%d/task/%d/children", pid, xpid)
		if _, err := os.Stat(cname); errors.Is(err, os.ErrNotExist) {
			continue
		}
		tids = append(tids, xpid)

	}
	for _, tid := range tids {
		chs, err := readProc(fmt.Sprintf("%d/task/%d/children", pid, tid))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, err
		}
		for _, ns := range strings.Split(string(chs), " ") {
			if ns == "" {
				continue
			}
			cp, err := strconv.Atoi(ns)
			if err != nil {
				return nil, err
			}
			if cp == pid {
				continue
			}
			res = append(res, PidStatus(cp))
			//		fmt.Printf("ChildPid: %d\n", cp)
		}
	}
	ps.direct_children = res
	return res, nil
}
func (ps *ProcessState) recursivelyGetChildrenOf() ([]*ProcessState, error) {
	res, err := ps.getChildrenOf()
	if err != nil {
		return nil, err
	}
	var childchilds []*ProcessState
	for _, cps := range res {
		chp, err := cps.recursivelyGetChildrenOf()
		if err != nil {
			return nil, err
		}
		childchilds = append(childchilds, chp...)
	}
	res = append(res, childchilds...)
	return res, nil
}

func PidStatus(pid int) *ProcessState {
	ps := &ProcessState{pid: pid}
	b, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	if err != nil {
		if ps.Status() != STATUS_RUNNING {
			return ps
		}
		ps.fail(err)
		return ps
	}
	ps.binary = b
	return ps
}
func (ps *ProcessState) fail(err error) {
	fmt.Printf("[go-easyops] linux error: %s\n", err)
	ps.err = err
}
func (ps *ProcessState) Pid() int {
	return ps.pid
}
func (ps *ProcessState) Binary() string {
	return ps.binary
}
func (ps *ProcessState) Status() STATUS {
	_, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", ps.Pid()))
	if err == nil {
		return STATUS_RUNNING
	}
	return STATUS_STOPPED
}

func (ps *ProcessState) String() string {
	return fmt.Sprintf("#%d (%s)", ps.Pid(), ps.Binary())
}
func (ps *ProcessState) Children() ([]*ProcessState, error) {
	if ps.direct_children == nil {
		_, err := ps.getChildrenOf()
		if err != nil {
			return nil, err
		}
	}
	return ps.direct_children, nil
}

func (s STATUS) String() string {
	if s == 1 {
		return "RUNNING"
	}
	if s == 2 {
		return "STOPPED"
	}
	return fmt.Sprintf("STATUS=%d", s)
}
