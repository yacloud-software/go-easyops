// client create: Workshop1ServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/workshop_1/workshop_1.proto
   gopackage : golang.conradwood.net/apis/workshop_1
   importname: ai_0
   varname   : client_Workshop1ServiceClient_0
   clientname: Workshop1ServiceClient
   servername: Workshop1ServiceServer
   gscvname  : workshop_1.Workshop1Service
   lockname  : lock_Workshop1ServiceClient_0
   activename: active_Workshop1ServiceClient_0
*/

package workshop_1

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_Workshop1ServiceClient_0 sync.Mutex
  client_Workshop1ServiceClient_0 Workshop1ServiceClient
)

func GetWorkshop1Client() Workshop1ServiceClient { 
    if client_Workshop1ServiceClient_0 != nil {
        return client_Workshop1ServiceClient_0
    }

    lock_Workshop1ServiceClient_0.Lock() 
    if client_Workshop1ServiceClient_0 != nil {
       lock_Workshop1ServiceClient_0.Unlock()
       return client_Workshop1ServiceClient_0
    }

    client_Workshop1ServiceClient_0 = NewWorkshop1ServiceClient(client.Connect("workshop_1.Workshop1Service"))
    lock_Workshop1ServiceClient_0.Unlock()
    return client_Workshop1ServiceClient_0
}

func GetWorkshop1ServiceClient() Workshop1ServiceClient { 
    if client_Workshop1ServiceClient_0 != nil {
        return client_Workshop1ServiceClient_0
    }

    lock_Workshop1ServiceClient_0.Lock() 
    if client_Workshop1ServiceClient_0 != nil {
       lock_Workshop1ServiceClient_0.Unlock()
       return client_Workshop1ServiceClient_0
    }

    client_Workshop1ServiceClient_0 = NewWorkshop1ServiceClient(client.Connect("workshop_1.Workshop1Service"))
    lock_Workshop1ServiceClient_0.Unlock()
    return client_Workshop1ServiceClient_0
}

