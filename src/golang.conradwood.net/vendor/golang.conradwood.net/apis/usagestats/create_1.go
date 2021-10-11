// client create: UsageStatsServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/usagestats/usagestats.proto
   gopackage : golang.conradwood.net/apis/usagestats
   importname: ai_0
   varname   : client_UsageStatsServiceClient_0
   clientname: UsageStatsServiceClient
   servername: UsageStatsServiceServer
   gscvname  : usagestats.UsageStatsService
   lockname  : lock_UsageStatsServiceClient_0
   activename: active_UsageStatsServiceClient_0
*/

package usagestats

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_UsageStatsServiceClient_0 sync.Mutex
  client_UsageStatsServiceClient_0 UsageStatsServiceClient
)

func GetUsageStatsClient() UsageStatsServiceClient { 
    if client_UsageStatsServiceClient_0 != nil {
        return client_UsageStatsServiceClient_0
    }

    lock_UsageStatsServiceClient_0.Lock() 
    if client_UsageStatsServiceClient_0 != nil {
       lock_UsageStatsServiceClient_0.Unlock()
       return client_UsageStatsServiceClient_0
    }

    client_UsageStatsServiceClient_0 = NewUsageStatsServiceClient(client.Connect("usagestats.UsageStatsService"))
    lock_UsageStatsServiceClient_0.Unlock()
    return client_UsageStatsServiceClient_0
}

func GetUsageStatsServiceClient() UsageStatsServiceClient { 
    if client_UsageStatsServiceClient_0 != nil {
        return client_UsageStatsServiceClient_0
    }

    lock_UsageStatsServiceClient_0.Lock() 
    if client_UsageStatsServiceClient_0 != nil {
       lock_UsageStatsServiceClient_0.Unlock()
       return client_UsageStatsServiceClient_0
    }

    client_UsageStatsServiceClient_0 = NewUsageStatsServiceClient(client.Connect("usagestats.UsageStatsService"))
    lock_UsageStatsServiceClient_0.Unlock()
    return client_UsageStatsServiceClient_0
}

