// client create: ExeWithMetricsClient
/* geninfo:
   filename  : conradwood.net/apis/exewithmetrics/exewithmetrics.proto
   gopackage : conradwood.net/apis/exewithmetrics
   importname: ai_0
   varname   : client_ExeWithMetricsClient_0
   clientname: ExeWithMetricsClient
   servername: ExeWithMetricsServer
   gscvname  : exewithmetrics.ExeWithMetrics
   lockname  : lock_ExeWithMetricsClient_0
   activename: active_ExeWithMetricsClient_0
*/

package exewithmetrics

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ExeWithMetricsClient_0 sync.Mutex
  client_ExeWithMetricsClient_0 ExeWithMetricsClient
)

func GetExeWithMetricsClient() ExeWithMetricsClient { 
    if client_ExeWithMetricsClient_0 != nil {
        return client_ExeWithMetricsClient_0
    }

    lock_ExeWithMetricsClient_0.Lock() 
    if client_ExeWithMetricsClient_0 != nil {
       lock_ExeWithMetricsClient_0.Unlock()
       return client_ExeWithMetricsClient_0
    }

    client_ExeWithMetricsClient_0 = NewExeWithMetricsClient(client.Connect("exewithmetrics.ExeWithMetrics"))
    lock_ExeWithMetricsClient_0.Unlock()
    return client_ExeWithMetricsClient_0
}

