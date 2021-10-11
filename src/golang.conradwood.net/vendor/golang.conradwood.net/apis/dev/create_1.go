// client create: DevServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/dev/dev.proto
   gopackage : golang.conradwood.net/apis/dev
   importname: ai_0
   varname   : client_DevServiceClient_0
   clientname: DevServiceClient
   servername: DevServiceServer
   gscvname  : dev.DevService
   lockname  : lock_DevServiceClient_0
   activename: active_DevServiceClient_0
*/

package dev

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_DevServiceClient_0 sync.Mutex
  client_DevServiceClient_0 DevServiceClient
)

func GetDevClient() DevServiceClient { 
    if client_DevServiceClient_0 != nil {
        return client_DevServiceClient_0
    }

    lock_DevServiceClient_0.Lock() 
    if client_DevServiceClient_0 != nil {
       lock_DevServiceClient_0.Unlock()
       return client_DevServiceClient_0
    }

    client_DevServiceClient_0 = NewDevServiceClient(client.Connect("dev.DevService"))
    lock_DevServiceClient_0.Unlock()
    return client_DevServiceClient_0
}

func GetDevServiceClient() DevServiceClient { 
    if client_DevServiceClient_0 != nil {
        return client_DevServiceClient_0
    }

    lock_DevServiceClient_0.Lock() 
    if client_DevServiceClient_0 != nil {
       lock_DevServiceClient_0.Unlock()
       return client_DevServiceClient_0
    }

    client_DevServiceClient_0 = NewDevServiceClient(client.Connect("dev.DevService"))
    lock_DevServiceClient_0.Unlock()
    return client_DevServiceClient_0
}

