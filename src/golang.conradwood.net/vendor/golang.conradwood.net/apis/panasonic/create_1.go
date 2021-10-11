// client create: PanasonicServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/panasonic/panasonic.proto
   gopackage : golang.conradwood.net/apis/panasonic
   importname: ai_0
   varname   : client_PanasonicServiceClient_0
   clientname: PanasonicServiceClient
   servername: PanasonicServiceServer
   gscvname  : panasonic.PanasonicService
   lockname  : lock_PanasonicServiceClient_0
   activename: active_PanasonicServiceClient_0
*/

package panasonic

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_PanasonicServiceClient_0 sync.Mutex
  client_PanasonicServiceClient_0 PanasonicServiceClient
)

func GetPanasonicClient() PanasonicServiceClient { 
    if client_PanasonicServiceClient_0 != nil {
        return client_PanasonicServiceClient_0
    }

    lock_PanasonicServiceClient_0.Lock() 
    if client_PanasonicServiceClient_0 != nil {
       lock_PanasonicServiceClient_0.Unlock()
       return client_PanasonicServiceClient_0
    }

    client_PanasonicServiceClient_0 = NewPanasonicServiceClient(client.Connect("panasonic.PanasonicService"))
    lock_PanasonicServiceClient_0.Unlock()
    return client_PanasonicServiceClient_0
}

func GetPanasonicServiceClient() PanasonicServiceClient { 
    if client_PanasonicServiceClient_0 != nil {
        return client_PanasonicServiceClient_0
    }

    lock_PanasonicServiceClient_0.Lock() 
    if client_PanasonicServiceClient_0 != nil {
       lock_PanasonicServiceClient_0.Unlock()
       return client_PanasonicServiceClient_0
    }

    client_PanasonicServiceClient_0 = NewPanasonicServiceClient(client.Connect("panasonic.PanasonicService"))
    lock_PanasonicServiceClient_0.Unlock()
    return client_PanasonicServiceClient_0
}

