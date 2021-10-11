// client create: SCServerClient
/* geninfo:
   filename  : golang.singingcat.net/apis/scserver/scserver.proto
   gopackage : golang.singingcat.net/apis/scserver
   importname: ai_0
   varname   : client_SCServerClient_0
   clientname: SCServerClient
   servername: SCServerServer
   gscvname  : scserver.SCServer
   lockname  : lock_SCServerClient_0
   activename: active_SCServerClient_0
*/

package scserver

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SCServerClient_0 sync.Mutex
  client_SCServerClient_0 SCServerClient
)

func GetSCServerClient() SCServerClient { 
    if client_SCServerClient_0 != nil {
        return client_SCServerClient_0
    }

    lock_SCServerClient_0.Lock() 
    if client_SCServerClient_0 != nil {
       lock_SCServerClient_0.Unlock()
       return client_SCServerClient_0
    }

    client_SCServerClient_0 = NewSCServerClient(client.Connect("scserver.SCServer"))
    lock_SCServerClient_0.Unlock()
    return client_SCServerClient_0
}

