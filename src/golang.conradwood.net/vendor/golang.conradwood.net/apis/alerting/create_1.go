// client create: AlertingServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/alerting/alerting.proto
   gopackage : golang.conradwood.net/apis/alerting
   importname: ai_0
   varname   : client_AlertingServiceClient_0
   clientname: AlertingServiceClient
   servername: AlertingServiceServer
   gscvname  : alerting.AlertingService
   lockname  : lock_AlertingServiceClient_0
   activename: active_AlertingServiceClient_0
*/

package alerting

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_AlertingServiceClient_0 sync.Mutex
  client_AlertingServiceClient_0 AlertingServiceClient
)

func GetAlertingClient() AlertingServiceClient { 
    if client_AlertingServiceClient_0 != nil {
        return client_AlertingServiceClient_0
    }

    lock_AlertingServiceClient_0.Lock() 
    if client_AlertingServiceClient_0 != nil {
       lock_AlertingServiceClient_0.Unlock()
       return client_AlertingServiceClient_0
    }

    client_AlertingServiceClient_0 = NewAlertingServiceClient(client.Connect("alerting.AlertingService"))
    lock_AlertingServiceClient_0.Unlock()
    return client_AlertingServiceClient_0
}

func GetAlertingServiceClient() AlertingServiceClient { 
    if client_AlertingServiceClient_0 != nil {
        return client_AlertingServiceClient_0
    }

    lock_AlertingServiceClient_0.Lock() 
    if client_AlertingServiceClient_0 != nil {
       lock_AlertingServiceClient_0.Unlock()
       return client_AlertingServiceClient_0
    }

    client_AlertingServiceClient_0 = NewAlertingServiceClient(client.Connect("alerting.AlertingService"))
    lock_AlertingServiceClient_0.Unlock()
    return client_AlertingServiceClient_0
}

