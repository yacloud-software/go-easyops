package server

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"golang.conradwood.net/apis/auth"
)

var (
	usages_tracking_enabled = flag.Bool("ge_track_usage_per_calling_service", false, "if true track usage per calling service")
	usages                  = &usage_info_tracker{usages: make(map[string]*usage_service)}
)

type usage_info_tracker struct {
	sync.Mutex
	usages map[string]*usage_service
}

type usage_service struct {
	sync.Mutex
	name             string
	usage_method_map map[string]*usage_method
}
type usage_method struct {
	sync.Mutex
	name    string
	callers []*usage_caller
}
type usage_caller struct {
	sync.Mutex
	user      *auth.User
	calls     int
	last_call time.Time
}

func (uit *usage_info_tracker) GetServiceByName(name string) *usage_service {
	uit.Lock()
	defer uit.Unlock()
	us := uit.usages[name]
	if us == nil {
		us = &usage_service{name: name, usage_method_map: make(map[string]*usage_method)}
		uit.usages[name] = us
	}
	return us
}
func (uit *usage_info_tracker) Services() []*usage_service {
	uit.Lock()
	defer uit.Unlock()
	var res []*usage_service
	for _, v := range uit.usages {
		res = append(res, v)
	}
	return res
}

func (us *usage_service) Methods() []*usage_method {
	us.Lock()
	defer us.Unlock()
	var res []*usage_method
	for _, v := range us.usage_method_map {
		res = append(res, v)
	}
	return res
}
func (us *usage_service) MethodByName(name string) *usage_method {
	us.Lock()
	defer us.Unlock()
	um := us.usage_method_map[name]
	if um == nil {
		um = &usage_method{name: name}
		us.usage_method_map[name] = um
	}
	return um
}
func (us *usage_service) Name() string {
	return us.name
}

func (um *usage_method) Callers() []*usage_caller {
	um.Lock()
	defer um.Unlock()
	var res []*usage_caller
	for _, v := range um.callers {
		res = append(res, v)
	}
	return res
}
func (um *usage_method) Name() string {
	return um.name
}
func (um *usage_method) CallerByUser(caller *auth.User) *usage_caller {
	um.Lock()
	defer um.Unlock()
	var uc *usage_caller
	for _, ucl := range um.callers {
		if ucl.user.ID == caller.ID {
			uc = ucl
			break
		}
	}
	if uc == nil {
		uc = &usage_caller{user: caller}
		um.callers = append(um.callers, uc)
	}
	return uc
}
func (uc *usage_caller) User() *auth.User {
	return uc.user
}

// how often was it called?
func (uc *usage_caller) Usages() int {
	return uc.calls
}

// when was last time it was called?
func (uc *usage_caller) LastCallTime() time.Time {
	return uc.last_call
}
func (uc *usage_caller) String() string {
	return fmt.Sprintf("%s %d %d (%s)", uc.user.ID, uc.calls, uc.last_call.Unix(), uc.user.Email)
}

func track_inbound_call(service, method string, caller *auth.User) {
	if !*usages_tracking_enabled {
		return
	}
	if caller == nil {
		return
	}
	us := usages.GetServiceByName(service)
	um := us.MethodByName(method)
	uc := um.CallerByUser(caller)

	uc.Lock()
	uc.calls++
	uc.last_call = time.Now()
	uc.Unlock()
}

func GetUsageInfo() *usage_info_tracker {
	return usages
}
