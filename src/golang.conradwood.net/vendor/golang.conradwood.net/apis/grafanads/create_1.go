// client create: GrafanaDSServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/grafanads/grafanads.proto
   gopackage : golang.conradwood.net/apis/grafanads
   importname: ai_0
   varname   : client_GrafanaDSServiceClient_0
   clientname: GrafanaDSServiceClient
   servername: GrafanaDSServiceServer
   gscvname  : grafanads.GrafanaDSService
   lockname  : lock_GrafanaDSServiceClient_0
   activename: active_GrafanaDSServiceClient_0
*/

package grafanads

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_GrafanaDSServiceClient_0 sync.Mutex
  client_GrafanaDSServiceClient_0 GrafanaDSServiceClient
)

func GetGrafanaDSClient() GrafanaDSServiceClient { 
    if client_GrafanaDSServiceClient_0 != nil {
        return client_GrafanaDSServiceClient_0
    }

    lock_GrafanaDSServiceClient_0.Lock() 
    if client_GrafanaDSServiceClient_0 != nil {
       lock_GrafanaDSServiceClient_0.Unlock()
       return client_GrafanaDSServiceClient_0
    }

    client_GrafanaDSServiceClient_0 = NewGrafanaDSServiceClient(client.Connect("grafanads.GrafanaDSService"))
    lock_GrafanaDSServiceClient_0.Unlock()
    return client_GrafanaDSServiceClient_0
}

func GetGrafanaDSServiceClient() GrafanaDSServiceClient { 
    if client_GrafanaDSServiceClient_0 != nil {
        return client_GrafanaDSServiceClient_0
    }

    lock_GrafanaDSServiceClient_0.Lock() 
    if client_GrafanaDSServiceClient_0 != nil {
       lock_GrafanaDSServiceClient_0.Unlock()
       return client_GrafanaDSServiceClient_0
    }

    client_GrafanaDSServiceClient_0 = NewGrafanaDSServiceClient(client.Connect("grafanads.GrafanaDSService"))
    lock_GrafanaDSServiceClient_0.Unlock()
    return client_GrafanaDSServiceClient_0
}

