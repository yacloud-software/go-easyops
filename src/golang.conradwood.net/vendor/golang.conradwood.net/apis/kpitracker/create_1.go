// client create: KPITrackerClient
/* geninfo:
   filename  : golang.conradwood.net/apis/kpitracker/kpitracker.proto
   gopackage : golang.conradwood.net/apis/kpitracker
   importname: ai_0
   varname   : client_KPITrackerClient_0
   clientname: KPITrackerClient
   servername: KPITrackerServer
   gscvname  : kpitracker.KPITracker
   lockname  : lock_KPITrackerClient_0
   activename: active_KPITrackerClient_0
*/

package kpitracker

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_KPITrackerClient_0 sync.Mutex
  client_KPITrackerClient_0 KPITrackerClient
)

func GetKPITrackerClient() KPITrackerClient { 
    if client_KPITrackerClient_0 != nil {
        return client_KPITrackerClient_0
    }

    lock_KPITrackerClient_0.Lock() 
    if client_KPITrackerClient_0 != nil {
       lock_KPITrackerClient_0.Unlock()
       return client_KPITrackerClient_0
    }

    client_KPITrackerClient_0 = NewKPITrackerClient(client.Connect("kpitracker.KPITracker"))
    lock_KPITrackerClient_0.Unlock()
    return client_KPITrackerClient_0
}

