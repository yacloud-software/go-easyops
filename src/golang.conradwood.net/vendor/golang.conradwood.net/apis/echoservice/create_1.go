// client create: EchoServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/echoservice/echoservice.proto
   gopackage : golang.conradwood.net/apis/echoservice
   importname: ai_0
   varname   : client_EchoServiceClient_0
   clientname: EchoServiceClient
   servername: EchoServiceServer
   gscvname  : echoservice.EchoService
   lockname  : lock_EchoServiceClient_0
   activename: active_EchoServiceClient_0
*/

package echoservice

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_EchoServiceClient_0 sync.Mutex
  client_EchoServiceClient_0 EchoServiceClient
)

func GetEchoClient() EchoServiceClient { 
    if client_EchoServiceClient_0 != nil {
        return client_EchoServiceClient_0
    }

    lock_EchoServiceClient_0.Lock() 
    if client_EchoServiceClient_0 != nil {
       lock_EchoServiceClient_0.Unlock()
       return client_EchoServiceClient_0
    }

    client_EchoServiceClient_0 = NewEchoServiceClient(client.Connect("echoservice.EchoService"))
    lock_EchoServiceClient_0.Unlock()
    return client_EchoServiceClient_0
}

func GetEchoServiceClient() EchoServiceClient { 
    if client_EchoServiceClient_0 != nil {
        return client_EchoServiceClient_0
    }

    lock_EchoServiceClient_0.Lock() 
    if client_EchoServiceClient_0 != nil {
       lock_EchoServiceClient_0.Unlock()
       return client_EchoServiceClient_0
    }

    client_EchoServiceClient_0 = NewEchoServiceClient(client.Connect("echoservice.EchoService"))
    lock_EchoServiceClient_0.Unlock()
    return client_EchoServiceClient_0
}

