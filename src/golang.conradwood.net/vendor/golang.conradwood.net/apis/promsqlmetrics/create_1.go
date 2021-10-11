// client create: PromSQLMetricsClient
/* geninfo:
   filename  : golang.conradwood.net/apis/promsqlmetrics/promsqlmetrics.proto
   gopackage : golang.conradwood.net/apis/promsqlmetrics
   importname: ai_0
   varname   : client_PromSQLMetricsClient_0
   clientname: PromSQLMetricsClient
   servername: PromSQLMetricsServer
   gscvname  : promsqlmetrics.PromSQLMetrics
   lockname  : lock_PromSQLMetricsClient_0
   activename: active_PromSQLMetricsClient_0
*/

package promsqlmetrics

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_PromSQLMetricsClient_0 sync.Mutex
  client_PromSQLMetricsClient_0 PromSQLMetricsClient
)

func GetPromSQLMetricsClient() PromSQLMetricsClient { 
    if client_PromSQLMetricsClient_0 != nil {
        return client_PromSQLMetricsClient_0
    }

    lock_PromSQLMetricsClient_0.Lock() 
    if client_PromSQLMetricsClient_0 != nil {
       lock_PromSQLMetricsClient_0.Unlock()
       return client_PromSQLMetricsClient_0
    }

    client_PromSQLMetricsClient_0 = NewPromSQLMetricsClient(client.Connect("promsqlmetrics.PromSQLMetrics"))
    lock_PromSQLMetricsClient_0.Unlock()
    return client_PromSQLMetricsClient_0
}

