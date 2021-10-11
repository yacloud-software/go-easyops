// client create: BuilderClient
/* geninfo:
   filename  : yacloud.eu/apis/builder/builder.proto
   gopackage : yacloud.eu/apis/builder
   importname: ai_0
   varname   : client_BuilderClient_0
   clientname: BuilderClient
   servername: BuilderServer
   gscvname  : builder.Builder
   lockname  : lock_BuilderClient_0
   activename: active_BuilderClient_0
*/

package builder

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_BuilderClient_0 sync.Mutex
  client_BuilderClient_0 BuilderClient
)

func GetBuilderClient() BuilderClient { 
    if client_BuilderClient_0 != nil {
        return client_BuilderClient_0
    }

    lock_BuilderClient_0.Lock() 
    if client_BuilderClient_0 != nil {
       lock_BuilderClient_0.Unlock()
       return client_BuilderClient_0
    }

    client_BuilderClient_0 = NewBuilderClient(client.Connect("builder.Builder"))
    lock_BuilderClient_0.Unlock()
    return client_BuilderClient_0
}

