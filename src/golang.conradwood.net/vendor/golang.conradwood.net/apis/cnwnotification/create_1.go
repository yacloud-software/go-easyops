// client create: CNWNotificationServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/cnwnotification/cnwnotification.proto
   gopackage : golang.conradwood.net/apis/cnwnotification
   importname: ai_0
   varname   : client_CNWNotificationServiceClient_0
   clientname: CNWNotificationServiceClient
   servername: CNWNotificationServiceServer
   gscvname  : cnwnotification.CNWNotificationService
   lockname  : lock_CNWNotificationServiceClient_0
   activename: active_CNWNotificationServiceClient_0
*/

package cnwnotification

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_CNWNotificationServiceClient_0 sync.Mutex
  client_CNWNotificationServiceClient_0 CNWNotificationServiceClient
)

func GetCNWNotificationClient() CNWNotificationServiceClient { 
    if client_CNWNotificationServiceClient_0 != nil {
        return client_CNWNotificationServiceClient_0
    }

    lock_CNWNotificationServiceClient_0.Lock() 
    if client_CNWNotificationServiceClient_0 != nil {
       lock_CNWNotificationServiceClient_0.Unlock()
       return client_CNWNotificationServiceClient_0
    }

    client_CNWNotificationServiceClient_0 = NewCNWNotificationServiceClient(client.Connect("cnwnotification.CNWNotificationService"))
    lock_CNWNotificationServiceClient_0.Unlock()
    return client_CNWNotificationServiceClient_0
}

func GetCNWNotificationServiceClient() CNWNotificationServiceClient { 
    if client_CNWNotificationServiceClient_0 != nil {
        return client_CNWNotificationServiceClient_0
    }

    lock_CNWNotificationServiceClient_0.Lock() 
    if client_CNWNotificationServiceClient_0 != nil {
       lock_CNWNotificationServiceClient_0.Unlock()
       return client_CNWNotificationServiceClient_0
    }

    client_CNWNotificationServiceClient_0 = NewCNWNotificationServiceClient(client.Connect("cnwnotification.CNWNotificationService"))
    lock_CNWNotificationServiceClient_0.Unlock()
    return client_CNWNotificationServiceClient_0
}

