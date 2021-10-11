// client create: IPManagerServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/ipmanager/ipmanager.proto
   gopackage : golang.conradwood.net/apis/ipmanager
   importname: ai_0
   varname   : client_IPManagerServiceClient_0
   clientname: IPManagerServiceClient
   servername: IPManagerServiceServer
   gscvname  : ipmanager.IPManagerService
   lockname  : lock_IPManagerServiceClient_0
   activename: active_IPManagerServiceClient_0
*/

package ipmanager

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_IPManagerServiceClient_0 sync.Mutex
  client_IPManagerServiceClient_0 IPManagerServiceClient
)

func GetIPManagerClient() IPManagerServiceClient { 
    if client_IPManagerServiceClient_0 != nil {
        return client_IPManagerServiceClient_0
    }

    lock_IPManagerServiceClient_0.Lock() 
    if client_IPManagerServiceClient_0 != nil {
       lock_IPManagerServiceClient_0.Unlock()
       return client_IPManagerServiceClient_0
    }

    client_IPManagerServiceClient_0 = NewIPManagerServiceClient(client.Connect("ipmanager.IPManagerService"))
    lock_IPManagerServiceClient_0.Unlock()
    return client_IPManagerServiceClient_0
}

func GetIPManagerServiceClient() IPManagerServiceClient { 
    if client_IPManagerServiceClient_0 != nil {
        return client_IPManagerServiceClient_0
    }

    lock_IPManagerServiceClient_0.Lock() 
    if client_IPManagerServiceClient_0 != nil {
       lock_IPManagerServiceClient_0.Unlock()
       return client_IPManagerServiceClient_0
    }

    client_IPManagerServiceClient_0 = NewIPManagerServiceClient(client.Connect("ipmanager.IPManagerService"))
    lock_IPManagerServiceClient_0.Unlock()
    return client_IPManagerServiceClient_0
}

