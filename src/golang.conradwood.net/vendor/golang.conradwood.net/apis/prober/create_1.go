// client create: ProberServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/prober/prober.proto
   gopackage : golang.conradwood.net/apis/prober
   importname: ai_0
   varname   : client_ProberServiceClient_0
   clientname: ProberServiceClient
   servername: ProberServiceServer
   gscvname  : prober.ProberService
   lockname  : lock_ProberServiceClient_0
   activename: active_ProberServiceClient_0
*/

package prober

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ProberServiceClient_0 sync.Mutex
  client_ProberServiceClient_0 ProberServiceClient
)

func GetProberClient() ProberServiceClient { 
    if client_ProberServiceClient_0 != nil {
        return client_ProberServiceClient_0
    }

    lock_ProberServiceClient_0.Lock() 
    if client_ProberServiceClient_0 != nil {
       lock_ProberServiceClient_0.Unlock()
       return client_ProberServiceClient_0
    }

    client_ProberServiceClient_0 = NewProberServiceClient(client.Connect("prober.ProberService"))
    lock_ProberServiceClient_0.Unlock()
    return client_ProberServiceClient_0
}

func GetProberServiceClient() ProberServiceClient { 
    if client_ProberServiceClient_0 != nil {
        return client_ProberServiceClient_0
    }

    lock_ProberServiceClient_0.Lock() 
    if client_ProberServiceClient_0 != nil {
       lock_ProberServiceClient_0.Unlock()
       return client_ProberServiceClient_0
    }

    client_ProberServiceClient_0 = NewProberServiceClient(client.Connect("prober.ProberService"))
    lock_ProberServiceClient_0.Unlock()
    return client_ProberServiceClient_0
}

