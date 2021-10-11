// client create: SCAclServiceClient
/* geninfo:
   filename  : golang.singingcat.net/apis/scacl/scacl.proto
   gopackage : golang.singingcat.net/apis/scacl
   importname: ai_0
   varname   : client_SCAclServiceClient_0
   clientname: SCAclServiceClient
   servername: SCAclServiceServer
   gscvname  : scacl.SCAclService
   lockname  : lock_SCAclServiceClient_0
   activename: active_SCAclServiceClient_0
*/

package scacl

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SCAclServiceClient_0 sync.Mutex
  client_SCAclServiceClient_0 SCAclServiceClient
)

func GetSCAclClient() SCAclServiceClient { 
    if client_SCAclServiceClient_0 != nil {
        return client_SCAclServiceClient_0
    }

    lock_SCAclServiceClient_0.Lock() 
    if client_SCAclServiceClient_0 != nil {
       lock_SCAclServiceClient_0.Unlock()
       return client_SCAclServiceClient_0
    }

    client_SCAclServiceClient_0 = NewSCAclServiceClient(client.Connect("scacl.SCAclService"))
    lock_SCAclServiceClient_0.Unlock()
    return client_SCAclServiceClient_0
}

func GetSCAclServiceClient() SCAclServiceClient { 
    if client_SCAclServiceClient_0 != nil {
        return client_SCAclServiceClient_0
    }

    lock_SCAclServiceClient_0.Lock() 
    if client_SCAclServiceClient_0 != nil {
       lock_SCAclServiceClient_0.Unlock()
       return client_SCAclServiceClient_0
    }

    client_SCAclServiceClient_0 = NewSCAclServiceClient(client.Connect("scacl.SCAclService"))
    lock_SCAclServiceClient_0.Unlock()
    return client_SCAclServiceClient_0
}

