// client create: TimeseriesBackendClient
/* geninfo:
   filename  : golang.conradwood.net/apis/timeseries/timeseries.proto
   gopackage : golang.conradwood.net/apis/timeseries
   importname: ai_1
   varname   : client_TimeseriesBackendClient_1
   clientname: TimeseriesBackendClient
   servername: TimeseriesBackendServer
   gscvname  : timeseries.TimeseriesBackend
   lockname  : lock_TimeseriesBackendClient_1
   activename: active_TimeseriesBackendClient_1
*/

package timeseries

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_TimeseriesBackendClient_1 sync.Mutex
  client_TimeseriesBackendClient_1 TimeseriesBackendClient
)

func GetTimeseriesBackendClient() TimeseriesBackendClient { 
    if client_TimeseriesBackendClient_1 != nil {
        return client_TimeseriesBackendClient_1
    }

    lock_TimeseriesBackendClient_1.Lock() 
    if client_TimeseriesBackendClient_1 != nil {
       lock_TimeseriesBackendClient_1.Unlock()
       return client_TimeseriesBackendClient_1
    }

    client_TimeseriesBackendClient_1 = NewTimeseriesBackendClient(client.Connect("timeseries.TimeseriesBackend"))
    lock_TimeseriesBackendClient_1.Unlock()
    return client_TimeseriesBackendClient_1
}

