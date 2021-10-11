// client create: HalleGateServiceClient
/* geninfo:
   filename  : golang.singingcat.net/apis/hallegate/hallegate.proto
   gopackage : golang.singingcat.net/apis/hallegate
   importname: ai_0
   varname   : client_HalleGateServiceClient_0
   clientname: HalleGateServiceClient
   servername: HalleGateServiceServer
   gscvname  : hallegate.HalleGateService
   lockname  : lock_HalleGateServiceClient_0
   activename: active_HalleGateServiceClient_0
*/

package hallegate

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_HalleGateServiceClient_0 sync.Mutex
  client_HalleGateServiceClient_0 HalleGateServiceClient
)

func GetHalleGateClient() HalleGateServiceClient { 
    if client_HalleGateServiceClient_0 != nil {
        return client_HalleGateServiceClient_0
    }

    lock_HalleGateServiceClient_0.Lock() 
    if client_HalleGateServiceClient_0 != nil {
       lock_HalleGateServiceClient_0.Unlock()
       return client_HalleGateServiceClient_0
    }

    client_HalleGateServiceClient_0 = NewHalleGateServiceClient(client.Connect("hallegate.HalleGateService"))
    lock_HalleGateServiceClient_0.Unlock()
    return client_HalleGateServiceClient_0
}

func GetHalleGateServiceClient() HalleGateServiceClient { 
    if client_HalleGateServiceClient_0 != nil {
        return client_HalleGateServiceClient_0
    }

    lock_HalleGateServiceClient_0.Lock() 
    if client_HalleGateServiceClient_0 != nil {
       lock_HalleGateServiceClient_0.Unlock()
       return client_HalleGateServiceClient_0
    }

    client_HalleGateServiceClient_0 = NewHalleGateServiceClient(client.Connect("hallegate.HalleGateService"))
    lock_HalleGateServiceClient_0.Unlock()
    return client_HalleGateServiceClient_0
}

