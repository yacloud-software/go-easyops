// client create: ApiTestServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/apitest/apitest.proto
   gopackage : golang.conradwood.net/apis/apitest
   importname: ai_0
   varname   : client_ApiTestServiceClient_0
   clientname: ApiTestServiceClient
   servername: ApiTestServiceServer
   gscvname  : apitest.ApiTestService
   lockname  : lock_ApiTestServiceClient_0
   activename: active_ApiTestServiceClient_0
*/

package apitest

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ApiTestServiceClient_0 sync.Mutex
  client_ApiTestServiceClient_0 ApiTestServiceClient
)

func GetApiTestClient() ApiTestServiceClient { 
    if client_ApiTestServiceClient_0 != nil {
        return client_ApiTestServiceClient_0
    }

    lock_ApiTestServiceClient_0.Lock() 
    if client_ApiTestServiceClient_0 != nil {
       lock_ApiTestServiceClient_0.Unlock()
       return client_ApiTestServiceClient_0
    }

    client_ApiTestServiceClient_0 = NewApiTestServiceClient(client.Connect("apitest.ApiTestService"))
    lock_ApiTestServiceClient_0.Unlock()
    return client_ApiTestServiceClient_0
}

func GetApiTestServiceClient() ApiTestServiceClient { 
    if client_ApiTestServiceClient_0 != nil {
        return client_ApiTestServiceClient_0
    }

    lock_ApiTestServiceClient_0.Lock() 
    if client_ApiTestServiceClient_0 != nil {
       lock_ApiTestServiceClient_0.Unlock()
       return client_ApiTestServiceClient_0
    }

    client_ApiTestServiceClient_0 = NewApiTestServiceClient(client.Connect("apitest.ApiTestService"))
    lock_ApiTestServiceClient_0.Unlock()
    return client_ApiTestServiceClient_0
}

