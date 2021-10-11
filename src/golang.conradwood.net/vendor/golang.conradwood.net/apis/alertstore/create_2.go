// client create: AlertStoreMgrServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/alertstore/alertstore.proto
   gopackage : golang.conradwood.net/apis/alertstore
   importname: ai_1
   varname   : client_AlertStoreMgrServiceClient_1
   clientname: AlertStoreMgrServiceClient
   servername: AlertStoreMgrServiceServer
   gscvname  : alertstore.AlertStoreMgrService
   lockname  : lock_AlertStoreMgrServiceClient_1
   activename: active_AlertStoreMgrServiceClient_1
*/

package alertstore

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_AlertStoreMgrServiceClient_1 sync.Mutex
  client_AlertStoreMgrServiceClient_1 AlertStoreMgrServiceClient
)

func GetAlertStoreMgrClient() AlertStoreMgrServiceClient { 
    if client_AlertStoreMgrServiceClient_1 != nil {
        return client_AlertStoreMgrServiceClient_1
    }

    lock_AlertStoreMgrServiceClient_1.Lock() 
    if client_AlertStoreMgrServiceClient_1 != nil {
       lock_AlertStoreMgrServiceClient_1.Unlock()
       return client_AlertStoreMgrServiceClient_1
    }

    client_AlertStoreMgrServiceClient_1 = NewAlertStoreMgrServiceClient(client.Connect("alertstore.AlertStoreMgrService"))
    lock_AlertStoreMgrServiceClient_1.Unlock()
    return client_AlertStoreMgrServiceClient_1
}

func GetAlertStoreMgrServiceClient() AlertStoreMgrServiceClient { 
    if client_AlertStoreMgrServiceClient_1 != nil {
        return client_AlertStoreMgrServiceClient_1
    }

    lock_AlertStoreMgrServiceClient_1.Lock() 
    if client_AlertStoreMgrServiceClient_1 != nil {
       lock_AlertStoreMgrServiceClient_1.Unlock()
       return client_AlertStoreMgrServiceClient_1
    }

    client_AlertStoreMgrServiceClient_1 = NewAlertStoreMgrServiceClient(client.Connect("alertstore.AlertStoreMgrService"))
    lock_AlertStoreMgrServiceClient_1.Unlock()
    return client_AlertStoreMgrServiceClient_1
}

