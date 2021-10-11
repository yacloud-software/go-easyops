// client create: SCApiServiceClient
/* geninfo:
   filename  : golang.singingcat.net/apis/scapi/scapi.proto
   gopackage : golang.singingcat.net/apis/scapi
   importname: ai_0
   varname   : client_SCApiServiceClient_0
   clientname: SCApiServiceClient
   servername: SCApiServiceServer
   gscvname  : scapi.SCApiService
   lockname  : lock_SCApiServiceClient_0
   activename: active_SCApiServiceClient_0
*/

package scapi

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SCApiServiceClient_0 sync.Mutex
  client_SCApiServiceClient_0 SCApiServiceClient
)

func GetSCApiClient() SCApiServiceClient { 
    if client_SCApiServiceClient_0 != nil {
        return client_SCApiServiceClient_0
    }

    lock_SCApiServiceClient_0.Lock() 
    if client_SCApiServiceClient_0 != nil {
       lock_SCApiServiceClient_0.Unlock()
       return client_SCApiServiceClient_0
    }

    client_SCApiServiceClient_0 = NewSCApiServiceClient(client.Connect("scapi.SCApiService"))
    lock_SCApiServiceClient_0.Unlock()
    return client_SCApiServiceClient_0
}

func GetSCApiServiceClient() SCApiServiceClient { 
    if client_SCApiServiceClient_0 != nil {
        return client_SCApiServiceClient_0
    }

    lock_SCApiServiceClient_0.Lock() 
    if client_SCApiServiceClient_0 != nil {
       lock_SCApiServiceClient_0.Unlock()
       return client_SCApiServiceClient_0
    }

    client_SCApiServiceClient_0 = NewSCApiServiceClient(client.Connect("scapi.SCApiService"))
    lock_SCApiServiceClient_0.Unlock()
    return client_SCApiServiceClient_0
}

