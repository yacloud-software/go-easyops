// client create: TestServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/test/test.proto
   gopackage : golang.conradwood.net/apis/test
   importname: ai_0
   varname   : client_TestServiceClient_0
   clientname: TestServiceClient
   servername: TestServiceServer
   gscvname  : test.TestService
   lockname  : lock_TestServiceClient_0
   activename: active_TestServiceClient_0
*/

package test

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_TestServiceClient_0 sync.Mutex
  client_TestServiceClient_0 TestServiceClient
)

func GetTestClient() TestServiceClient { 
    if client_TestServiceClient_0 != nil {
        return client_TestServiceClient_0
    }

    lock_TestServiceClient_0.Lock() 
    if client_TestServiceClient_0 != nil {
       lock_TestServiceClient_0.Unlock()
       return client_TestServiceClient_0
    }

    client_TestServiceClient_0 = NewTestServiceClient(client.Connect("test.TestService"))
    lock_TestServiceClient_0.Unlock()
    return client_TestServiceClient_0
}

func GetTestServiceClient() TestServiceClient { 
    if client_TestServiceClient_0 != nil {
        return client_TestServiceClient_0
    }

    lock_TestServiceClient_0.Lock() 
    if client_TestServiceClient_0 != nil {
       lock_TestServiceClient_0.Unlock()
       return client_TestServiceClient_0
    }

    client_TestServiceClient_0 = NewTestServiceClient(client.Connect("test.TestService"))
    lock_TestServiceClient_0.Unlock()
    return client_TestServiceClient_0
}

