// client create: Workshop1SQLServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/workshop_1/workshop_1.proto
   gopackage : golang.conradwood.net/apis/workshop_1
   importname: ai_1
   varname   : client_Workshop1SQLServiceClient_1
   clientname: Workshop1SQLServiceClient
   servername: Workshop1SQLServiceServer
   gscvname  : workshop_1.Workshop1SQLService
   lockname  : lock_Workshop1SQLServiceClient_1
   activename: active_Workshop1SQLServiceClient_1
*/

package workshop_1

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_Workshop1SQLServiceClient_1 sync.Mutex
  client_Workshop1SQLServiceClient_1 Workshop1SQLServiceClient
)

func GetWorkshop1SQLClient() Workshop1SQLServiceClient { 
    if client_Workshop1SQLServiceClient_1 != nil {
        return client_Workshop1SQLServiceClient_1
    }

    lock_Workshop1SQLServiceClient_1.Lock() 
    if client_Workshop1SQLServiceClient_1 != nil {
       lock_Workshop1SQLServiceClient_1.Unlock()
       return client_Workshop1SQLServiceClient_1
    }

    client_Workshop1SQLServiceClient_1 = NewWorkshop1SQLServiceClient(client.Connect("workshop_1.Workshop1SQLService"))
    lock_Workshop1SQLServiceClient_1.Unlock()
    return client_Workshop1SQLServiceClient_1
}

func GetWorkshop1SQLServiceClient() Workshop1SQLServiceClient { 
    if client_Workshop1SQLServiceClient_1 != nil {
        return client_Workshop1SQLServiceClient_1
    }

    lock_Workshop1SQLServiceClient_1.Lock() 
    if client_Workshop1SQLServiceClient_1 != nil {
       lock_Workshop1SQLServiceClient_1.Unlock()
       return client_Workshop1SQLServiceClient_1
    }

    client_Workshop1SQLServiceClient_1 = NewWorkshop1SQLServiceClient(client.Connect("workshop_1.Workshop1SQLService"))
    lock_Workshop1SQLServiceClient_1.Unlock()
    return client_Workshop1SQLServiceClient_1
}

