// client create: StarlingServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/starling/starling.proto
   gopackage : golang.conradwood.net/apis/starling
   importname: ai_0
   varname   : client_StarlingServiceClient_0
   clientname: StarlingServiceClient
   servername: StarlingServiceServer
   gscvname  : starling.StarlingService
   lockname  : lock_StarlingServiceClient_0
   activename: active_StarlingServiceClient_0
*/

package starling

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_StarlingServiceClient_0 sync.Mutex
  client_StarlingServiceClient_0 StarlingServiceClient
)

func GetStarlingClient() StarlingServiceClient { 
    if client_StarlingServiceClient_0 != nil {
        return client_StarlingServiceClient_0
    }

    lock_StarlingServiceClient_0.Lock() 
    if client_StarlingServiceClient_0 != nil {
       lock_StarlingServiceClient_0.Unlock()
       return client_StarlingServiceClient_0
    }

    client_StarlingServiceClient_0 = NewStarlingServiceClient(client.Connect("starling.StarlingService"))
    lock_StarlingServiceClient_0.Unlock()
    return client_StarlingServiceClient_0
}

func GetStarlingServiceClient() StarlingServiceClient { 
    if client_StarlingServiceClient_0 != nil {
        return client_StarlingServiceClient_0
    }

    lock_StarlingServiceClient_0.Lock() 
    if client_StarlingServiceClient_0 != nil {
       lock_StarlingServiceClient_0.Unlock()
       return client_StarlingServiceClient_0
    }

    client_StarlingServiceClient_0 = NewStarlingServiceClient(client.Connect("starling.StarlingService"))
    lock_StarlingServiceClient_0.Unlock()
    return client_StarlingServiceClient_0
}

