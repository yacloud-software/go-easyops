// client create: SCPrometheusServiceClient
/* geninfo:
   filename  : golang.singingcat.net/apis/scprometheus/scprometheus.proto
   gopackage : golang.singingcat.net/apis/scprometheus
   importname: ai_0
   varname   : client_SCPrometheusServiceClient_0
   clientname: SCPrometheusServiceClient
   servername: SCPrometheusServiceServer
   gscvname  : scprometheus.SCPrometheusService
   lockname  : lock_SCPrometheusServiceClient_0
   activename: active_SCPrometheusServiceClient_0
*/

package scprometheus

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SCPrometheusServiceClient_0 sync.Mutex
  client_SCPrometheusServiceClient_0 SCPrometheusServiceClient
)

func GetSCPrometheusClient() SCPrometheusServiceClient { 
    if client_SCPrometheusServiceClient_0 != nil {
        return client_SCPrometheusServiceClient_0
    }

    lock_SCPrometheusServiceClient_0.Lock() 
    if client_SCPrometheusServiceClient_0 != nil {
       lock_SCPrometheusServiceClient_0.Unlock()
       return client_SCPrometheusServiceClient_0
    }

    client_SCPrometheusServiceClient_0 = NewSCPrometheusServiceClient(client.Connect("scprometheus.SCPrometheusService"))
    lock_SCPrometheusServiceClient_0.Unlock()
    return client_SCPrometheusServiceClient_0
}

func GetSCPrometheusServiceClient() SCPrometheusServiceClient { 
    if client_SCPrometheusServiceClient_0 != nil {
        return client_SCPrometheusServiceClient_0
    }

    lock_SCPrometheusServiceClient_0.Lock() 
    if client_SCPrometheusServiceClient_0 != nil {
       lock_SCPrometheusServiceClient_0.Unlock()
       return client_SCPrometheusServiceClient_0
    }

    client_SCPrometheusServiceClient_0 = NewSCPrometheusServiceClient(client.Connect("scprometheus.SCPrometheusService"))
    lock_SCPrometheusServiceClient_0.Unlock()
    return client_SCPrometheusServiceClient_0
}

