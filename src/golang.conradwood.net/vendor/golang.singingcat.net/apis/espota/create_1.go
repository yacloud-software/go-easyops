// client create: ESPOtaServiceClient
/* geninfo:
   filename  : golang.singingcat.net/apis/espota/espota.proto
   gopackage : golang.singingcat.net/apis/espota
   importname: ai_0
   varname   : client_ESPOtaServiceClient_0
   clientname: ESPOtaServiceClient
   servername: ESPOtaServiceServer
   gscvname  : espota.ESPOtaService
   lockname  : lock_ESPOtaServiceClient_0
   activename: active_ESPOtaServiceClient_0
*/

package espota

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ESPOtaServiceClient_0 sync.Mutex
  client_ESPOtaServiceClient_0 ESPOtaServiceClient
)

func GetESPOtaClient() ESPOtaServiceClient { 
    if client_ESPOtaServiceClient_0 != nil {
        return client_ESPOtaServiceClient_0
    }

    lock_ESPOtaServiceClient_0.Lock() 
    if client_ESPOtaServiceClient_0 != nil {
       lock_ESPOtaServiceClient_0.Unlock()
       return client_ESPOtaServiceClient_0
    }

    client_ESPOtaServiceClient_0 = NewESPOtaServiceClient(client.Connect("espota.ESPOtaService"))
    lock_ESPOtaServiceClient_0.Unlock()
    return client_ESPOtaServiceClient_0
}

func GetESPOtaServiceClient() ESPOtaServiceClient { 
    if client_ESPOtaServiceClient_0 != nil {
        return client_ESPOtaServiceClient_0
    }

    lock_ESPOtaServiceClient_0.Lock() 
    if client_ESPOtaServiceClient_0 != nil {
       lock_ESPOtaServiceClient_0.Unlock()
       return client_ESPOtaServiceClient_0
    }

    client_ESPOtaServiceClient_0 = NewESPOtaServiceClient(client.Connect("espota.ESPOtaService"))
    lock_ESPOtaServiceClient_0.Unlock()
    return client_ESPOtaServiceClient_0
}

