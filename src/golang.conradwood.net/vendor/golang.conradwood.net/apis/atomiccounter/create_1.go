// client create: AtomicCounterServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/atomiccounter/atomiccounter.proto
   gopackage : golang.conradwood.net/apis/atomiccounter
   importname: ai_0
   varname   : client_AtomicCounterServiceClient_0
   clientname: AtomicCounterServiceClient
   servername: AtomicCounterServiceServer
   gscvname  : atomiccounter.AtomicCounterService
   lockname  : lock_AtomicCounterServiceClient_0
   activename: active_AtomicCounterServiceClient_0
*/

package atomiccounter

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_AtomicCounterServiceClient_0 sync.Mutex
  client_AtomicCounterServiceClient_0 AtomicCounterServiceClient
)

func GetAtomicCounterClient() AtomicCounterServiceClient { 
    if client_AtomicCounterServiceClient_0 != nil {
        return client_AtomicCounterServiceClient_0
    }

    lock_AtomicCounterServiceClient_0.Lock() 
    if client_AtomicCounterServiceClient_0 != nil {
       lock_AtomicCounterServiceClient_0.Unlock()
       return client_AtomicCounterServiceClient_0
    }

    client_AtomicCounterServiceClient_0 = NewAtomicCounterServiceClient(client.Connect("atomiccounter.AtomicCounterService"))
    lock_AtomicCounterServiceClient_0.Unlock()
    return client_AtomicCounterServiceClient_0
}

func GetAtomicCounterServiceClient() AtomicCounterServiceClient { 
    if client_AtomicCounterServiceClient_0 != nil {
        return client_AtomicCounterServiceClient_0
    }

    lock_AtomicCounterServiceClient_0.Lock() 
    if client_AtomicCounterServiceClient_0 != nil {
       lock_AtomicCounterServiceClient_0.Unlock()
       return client_AtomicCounterServiceClient_0
    }

    client_AtomicCounterServiceClient_0 = NewAtomicCounterServiceClient(client.Connect("atomiccounter.AtomicCounterService"))
    lock_AtomicCounterServiceClient_0.Unlock()
    return client_AtomicCounterServiceClient_0
}

