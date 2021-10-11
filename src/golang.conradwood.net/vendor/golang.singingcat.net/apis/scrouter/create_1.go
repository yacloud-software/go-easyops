// client create: SCRouterClient
/* geninfo:
   filename  : golang.singingcat.net/apis/scrouter/scrouter.proto
   gopackage : golang.singingcat.net/apis/scrouter
   importname: ai_0
   varname   : client_SCRouterClient_0
   clientname: SCRouterClient
   servername: SCRouterServer
   gscvname  : scrouter.SCRouter
   lockname  : lock_SCRouterClient_0
   activename: active_SCRouterClient_0
*/

package scrouter

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SCRouterClient_0 sync.Mutex
  client_SCRouterClient_0 SCRouterClient
)

func GetSCRouterClient() SCRouterClient { 
    if client_SCRouterClient_0 != nil {
        return client_SCRouterClient_0
    }

    lock_SCRouterClient_0.Lock() 
    if client_SCRouterClient_0 != nil {
       lock_SCRouterClient_0.Unlock()
       return client_SCRouterClient_0
    }

    client_SCRouterClient_0 = NewSCRouterClient(client.Connect("scrouter.SCRouter"))
    lock_SCRouterClient_0.Unlock()
    return client_SCRouterClient_0
}

