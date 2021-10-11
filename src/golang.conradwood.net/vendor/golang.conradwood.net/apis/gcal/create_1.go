// client create: GCALServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/gcal/gcal.proto
   gopackage : golang.conradwood.net/apis/gcal
   importname: ai_0
   varname   : client_GCALServiceClient_0
   clientname: GCALServiceClient
   servername: GCALServiceServer
   gscvname  : gcal.GCALService
   lockname  : lock_GCALServiceClient_0
   activename: active_GCALServiceClient_0
*/

package gcal

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_GCALServiceClient_0 sync.Mutex
  client_GCALServiceClient_0 GCALServiceClient
)

func GetGCALClient() GCALServiceClient { 
    if client_GCALServiceClient_0 != nil {
        return client_GCALServiceClient_0
    }

    lock_GCALServiceClient_0.Lock() 
    if client_GCALServiceClient_0 != nil {
       lock_GCALServiceClient_0.Unlock()
       return client_GCALServiceClient_0
    }

    client_GCALServiceClient_0 = NewGCALServiceClient(client.Connect("gcal.GCALService"))
    lock_GCALServiceClient_0.Unlock()
    return client_GCALServiceClient_0
}

func GetGCALServiceClient() GCALServiceClient { 
    if client_GCALServiceClient_0 != nil {
        return client_GCALServiceClient_0
    }

    lock_GCALServiceClient_0.Lock() 
    if client_GCALServiceClient_0 != nil {
       lock_GCALServiceClient_0.Unlock()
       return client_GCALServiceClient_0
    }

    client_GCALServiceClient_0 = NewGCALServiceClient(client.Connect("gcal.GCALService"))
    lock_GCALServiceClient_0.Unlock()
    return client_GCALServiceClient_0
}

