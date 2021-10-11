// client create: SCFunctionsServerClient
/* geninfo:
   filename  : golang.singingcat.net/apis/scfunctions/scfunctions.proto
   gopackage : golang.singingcat.net/apis/scfunctions
   importname: ai_0
   varname   : client_SCFunctionsServerClient_0
   clientname: SCFunctionsServerClient
   servername: SCFunctionsServerServer
   gscvname  : scfunctions.SCFunctionsServer
   lockname  : lock_SCFunctionsServerClient_0
   activename: active_SCFunctionsServerClient_0
*/

package scfunctions

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SCFunctionsServerClient_0 sync.Mutex
  client_SCFunctionsServerClient_0 SCFunctionsServerClient
)

func GetSCFunctionsServerClient() SCFunctionsServerClient { 
    if client_SCFunctionsServerClient_0 != nil {
        return client_SCFunctionsServerClient_0
    }

    lock_SCFunctionsServerClient_0.Lock() 
    if client_SCFunctionsServerClient_0 != nil {
       lock_SCFunctionsServerClient_0.Unlock()
       return client_SCFunctionsServerClient_0
    }

    client_SCFunctionsServerClient_0 = NewSCFunctionsServerClient(client.Connect("scfunctions.SCFunctionsServer"))
    lock_SCFunctionsServerClient_0.Unlock()
    return client_SCFunctionsServerClient_0
}

