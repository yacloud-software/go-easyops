// client create: AlertStoreServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/alertstore/alertstore.proto
   gopackage : golang.conradwood.net/apis/alertstore
   importname: ai_0
   varname   : client_AlertStoreServiceClient_0
   clientname: AlertStoreServiceClient
   servername: AlertStoreServiceServer
   gscvname  : alertstore.AlertStoreService
   lockname  : lock_AlertStoreServiceClient_0
   activename: active_AlertStoreServiceClient_0
*/

package alertstore

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_AlertStoreServiceClient_0 sync.Mutex
  client_AlertStoreServiceClient_0 AlertStoreServiceClient
)

func GetAlertStoreClient() AlertStoreServiceClient { 
    if client_AlertStoreServiceClient_0 != nil {
        return client_AlertStoreServiceClient_0
    }

    lock_AlertStoreServiceClient_0.Lock() 
    if client_AlertStoreServiceClient_0 != nil {
       lock_AlertStoreServiceClient_0.Unlock()
       return client_AlertStoreServiceClient_0
    }

    client_AlertStoreServiceClient_0 = NewAlertStoreServiceClient(client.Connect("alertstore.AlertStoreService"))
    lock_AlertStoreServiceClient_0.Unlock()
    return client_AlertStoreServiceClient_0
}

func GetAlertStoreServiceClient() AlertStoreServiceClient { 
    if client_AlertStoreServiceClient_0 != nil {
        return client_AlertStoreServiceClient_0
    }

    lock_AlertStoreServiceClient_0.Lock() 
    if client_AlertStoreServiceClient_0 != nil {
       lock_AlertStoreServiceClient_0.Unlock()
       return client_AlertStoreServiceClient_0
    }

    client_AlertStoreServiceClient_0 = NewAlertStoreServiceClient(client.Connect("alertstore.AlertStoreService"))
    lock_AlertStoreServiceClient_0.Unlock()
    return client_AlertStoreServiceClient_0
}

