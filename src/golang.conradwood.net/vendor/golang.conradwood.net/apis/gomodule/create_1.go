// client create: GoModuleServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/gomodule/gomodule.proto
   gopackage : golang.conradwood.net/apis/gomodule
   importname: ai_0
   varname   : client_GoModuleServiceClient_0
   clientname: GoModuleServiceClient
   servername: GoModuleServiceServer
   gscvname  : gomodule.GoModuleService
   lockname  : lock_GoModuleServiceClient_0
   activename: active_GoModuleServiceClient_0
*/

package gomodule

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_GoModuleServiceClient_0 sync.Mutex
  client_GoModuleServiceClient_0 GoModuleServiceClient
)

func GetGoModuleClient() GoModuleServiceClient { 
    if client_GoModuleServiceClient_0 != nil {
        return client_GoModuleServiceClient_0
    }

    lock_GoModuleServiceClient_0.Lock() 
    if client_GoModuleServiceClient_0 != nil {
       lock_GoModuleServiceClient_0.Unlock()
       return client_GoModuleServiceClient_0
    }

    client_GoModuleServiceClient_0 = NewGoModuleServiceClient(client.Connect("gomodule.GoModuleService"))
    lock_GoModuleServiceClient_0.Unlock()
    return client_GoModuleServiceClient_0
}

func GetGoModuleServiceClient() GoModuleServiceClient { 
    if client_GoModuleServiceClient_0 != nil {
        return client_GoModuleServiceClient_0
    }

    lock_GoModuleServiceClient_0.Lock() 
    if client_GoModuleServiceClient_0 != nil {
       lock_GoModuleServiceClient_0.Unlock()
       return client_GoModuleServiceClient_0
    }

    client_GoModuleServiceClient_0 = NewGoModuleServiceClient(client.Connect("gomodule.GoModuleService"))
    lock_GoModuleServiceClient_0.Unlock()
    return client_GoModuleServiceClient_0
}

