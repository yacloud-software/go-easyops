// client create: TimeseriesServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/timeseries/timeseries.proto
   gopackage : golang.conradwood.net/apis/timeseries
   importname: ai_0
   varname   : client_TimeseriesServiceClient_0
   clientname: TimeseriesServiceClient
   servername: TimeseriesServiceServer
   gscvname  : timeseries.TimeseriesService
   lockname  : lock_TimeseriesServiceClient_0
   activename: active_TimeseriesServiceClient_0
*/

package timeseries

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_TimeseriesServiceClient_0 sync.Mutex
  client_TimeseriesServiceClient_0 TimeseriesServiceClient
)

func GetTimeseriesClient() TimeseriesServiceClient { 
    if client_TimeseriesServiceClient_0 != nil {
        return client_TimeseriesServiceClient_0
    }

    lock_TimeseriesServiceClient_0.Lock() 
    if client_TimeseriesServiceClient_0 != nil {
       lock_TimeseriesServiceClient_0.Unlock()
       return client_TimeseriesServiceClient_0
    }

    client_TimeseriesServiceClient_0 = NewTimeseriesServiceClient(client.Connect("timeseries.TimeseriesService"))
    lock_TimeseriesServiceClient_0.Unlock()
    return client_TimeseriesServiceClient_0
}

func GetTimeseriesServiceClient() TimeseriesServiceClient { 
    if client_TimeseriesServiceClient_0 != nil {
        return client_TimeseriesServiceClient_0
    }

    lock_TimeseriesServiceClient_0.Lock() 
    if client_TimeseriesServiceClient_0 != nil {
       lock_TimeseriesServiceClient_0.Unlock()
       return client_TimeseriesServiceClient_0
    }

    client_TimeseriesServiceClient_0 = NewTimeseriesServiceClient(client.Connect("timeseries.TimeseriesService"))
    lock_TimeseriesServiceClient_0.Unlock()
    return client_TimeseriesServiceClient_0
}

