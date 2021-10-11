// client create: GerritHooksClient
/* geninfo:
   filename  : golang.conradwood.net/apis/gerrithooks/gerrithooks.proto
   gopackage : golang.conradwood.net/apis/gerrithooks
   importname: ai_0
   varname   : client_GerritHooksClient_0
   clientname: GerritHooksClient
   servername: GerritHooksServer
   gscvname  : gerrithooks.GerritHooks
   lockname  : lock_GerritHooksClient_0
   activename: active_GerritHooksClient_0
*/

package gerrithooks

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_GerritHooksClient_0 sync.Mutex
  client_GerritHooksClient_0 GerritHooksClient
)

func GetGerritHooksClient() GerritHooksClient { 
    if client_GerritHooksClient_0 != nil {
        return client_GerritHooksClient_0
    }

    lock_GerritHooksClient_0.Lock() 
    if client_GerritHooksClient_0 != nil {
       lock_GerritHooksClient_0.Unlock()
       return client_GerritHooksClient_0
    }

    client_GerritHooksClient_0 = NewGerritHooksClient(client.Connect("gerrithooks.GerritHooks"))
    lock_GerritHooksClient_0.Unlock()
    return client_GerritHooksClient_0
}

