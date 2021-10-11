// client create: SCUtilsServerClient
/* geninfo:
   filename  : golang.singingcat.net/apis/scutils/scutils.proto
   gopackage : golang.singingcat.net/apis/scutils
   importname: ai_0
   varname   : client_SCUtilsServerClient_0
   clientname: SCUtilsServerClient
   servername: SCUtilsServerServer
   gscvname  : scutils.SCUtilsServer
   lockname  : lock_SCUtilsServerClient_0
   activename: active_SCUtilsServerClient_0
*/

package scutils

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SCUtilsServerClient_0 sync.Mutex
  client_SCUtilsServerClient_0 SCUtilsServerClient
)

func GetSCUtilsServerClient() SCUtilsServerClient { 
    if client_SCUtilsServerClient_0 != nil {
        return client_SCUtilsServerClient_0
    }

    lock_SCUtilsServerClient_0.Lock() 
    if client_SCUtilsServerClient_0 != nil {
       lock_SCUtilsServerClient_0.Unlock()
       return client_SCUtilsServerClient_0
    }

    client_SCUtilsServerClient_0 = NewSCUtilsServerClient(client.Connect("scutils.SCUtilsServer"))
    lock_SCUtilsServerClient_0.Unlock()
    return client_SCUtilsServerClient_0
}

