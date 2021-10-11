// client create: GrafanaDataClient
/* geninfo:
   filename  : conradwood.net/apis/grafanadata/grafanadata.proto
   gopackage : conradwood.net/apis/grafanadata
   importname: ai_0
   varname   : client_GrafanaDataClient_0
   clientname: GrafanaDataClient
   servername: GrafanaDataServer
   gscvname  : grafanadata.GrafanaData
   lockname  : lock_GrafanaDataClient_0
   activename: active_GrafanaDataClient_0
*/

package grafanadata

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_GrafanaDataClient_0 sync.Mutex
  client_GrafanaDataClient_0 GrafanaDataClient
)

func GetGrafanaDataClient() GrafanaDataClient { 
    if client_GrafanaDataClient_0 != nil {
        return client_GrafanaDataClient_0
    }

    lock_GrafanaDataClient_0.Lock() 
    if client_GrafanaDataClient_0 != nil {
       lock_GrafanaDataClient_0.Unlock()
       return client_GrafanaDataClient_0
    }

    client_GrafanaDataClient_0 = NewGrafanaDataClient(client.Connect("grafanadata.GrafanaData"))
    lock_GrafanaDataClient_0.Unlock()
    return client_GrafanaDataClient_0
}

