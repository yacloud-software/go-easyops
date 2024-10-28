package common

import (
	"fmt"
	"sync"
)

var (
	infolock       sync.Mutex
	info_providers []*infoProviderEntry
)

type infoProviderEntry struct {
	name string
	f    func() []*InfoValue
}

type InfoValue struct {
	Name  string
	Value float64
}

func RegisterInfoProvider(name string, ip func() []*InfoValue) {
	ipe := &infoProviderEntry{name: name, f: ip}
	infolock.Lock()
	info_providers = append(info_providers, ipe)
	infolock.Unlock()
	return
}

func GetText() map[string]string {
	res := make(map[string]string)
	infolock.Lock()
	for _, ip := range info_providers {
		s := ""
		ivs := ip.f()
		for _, iv := range ivs {
			x := fmt.Sprintf("%s:%0.2f", iv.Name, iv.Value)
			s = s + x
		}
	}
	infolock.Unlock()
	return res
}
