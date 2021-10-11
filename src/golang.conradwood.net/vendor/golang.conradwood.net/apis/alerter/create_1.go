// client create: AlerterServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/alerter/alerter.proto
   gopackage : golang.conradwood.net/apis/alerter
   importname: ai_0
   varname   : client_AlerterServiceClient_0
   clientname: AlerterServiceClient
   servername: AlerterServiceServer
   gscvname  : alerter.AlerterService
   lockname  : lock_AlerterServiceClient_0
   activename: active_AlerterServiceClient_0
*/

package alerter

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_AlerterServiceClient_0 sync.Mutex
  client_AlerterServiceClient_0 AlerterServiceClient
)

func GetAlerterClient() AlerterServiceClient { 
    if client_AlerterServiceClient_0 != nil {
        return client_AlerterServiceClient_0
    }

    lock_AlerterServiceClient_0.Lock() 
    if client_AlerterServiceClient_0 != nil {
       lock_AlerterServiceClient_0.Unlock()
       return client_AlerterServiceClient_0
    }

    client_AlerterServiceClient_0 = NewAlerterServiceClient(client.Connect("alerter.AlerterService"))
    lock_AlerterServiceClient_0.Unlock()
    return client_AlerterServiceClient_0
}

func GetAlerterServiceClient() AlerterServiceClient { 
    if client_AlerterServiceClient_0 != nil {
        return client_AlerterServiceClient_0
    }

    lock_AlerterServiceClient_0.Lock() 
    if client_AlerterServiceClient_0 != nil {
       lock_AlerterServiceClient_0.Unlock()
       return client_AlerterServiceClient_0
    }

    client_AlerterServiceClient_0 = NewAlerterServiceClient(client.Connect("alerter.AlerterService"))
    lock_AlerterServiceClient_0.Unlock()
    return client_AlerterServiceClient_0
}

