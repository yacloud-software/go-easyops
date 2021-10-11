// client create: PrometheusAPIServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/prometheusapi/prometheusapi.proto
   gopackage : golang.conradwood.net/apis/prometheusapi
   importname: ai_0
   varname   : client_PrometheusAPIServiceClient_0
   clientname: PrometheusAPIServiceClient
   servername: PrometheusAPIServiceServer
   gscvname  : prometheusapi.PrometheusAPIService
   lockname  : lock_PrometheusAPIServiceClient_0
   activename: active_PrometheusAPIServiceClient_0
*/

package prometheusapi

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_PrometheusAPIServiceClient_0 sync.Mutex
  client_PrometheusAPIServiceClient_0 PrometheusAPIServiceClient
)

func GetPrometheusAPIClient() PrometheusAPIServiceClient { 
    if client_PrometheusAPIServiceClient_0 != nil {
        return client_PrometheusAPIServiceClient_0
    }

    lock_PrometheusAPIServiceClient_0.Lock() 
    if client_PrometheusAPIServiceClient_0 != nil {
       lock_PrometheusAPIServiceClient_0.Unlock()
       return client_PrometheusAPIServiceClient_0
    }

    client_PrometheusAPIServiceClient_0 = NewPrometheusAPIServiceClient(client.Connect("prometheusapi.PrometheusAPIService"))
    lock_PrometheusAPIServiceClient_0.Unlock()
    return client_PrometheusAPIServiceClient_0
}

func GetPrometheusAPIServiceClient() PrometheusAPIServiceClient { 
    if client_PrometheusAPIServiceClient_0 != nil {
        return client_PrometheusAPIServiceClient_0
    }

    lock_PrometheusAPIServiceClient_0.Lock() 
    if client_PrometheusAPIServiceClient_0 != nil {
       lock_PrometheusAPIServiceClient_0.Unlock()
       return client_PrometheusAPIServiceClient_0
    }

    client_PrometheusAPIServiceClient_0 = NewPrometheusAPIServiceClient(client.Connect("prometheusapi.PrometheusAPIService"))
    lock_PrometheusAPIServiceClient_0.Unlock()
    return client_PrometheusAPIServiceClient_0
}

