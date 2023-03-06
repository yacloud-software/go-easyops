// client create: EchoServiceClient
/*
  Created by /home/cnw/devel/go/go-tools/src/golang.conradwood.net/gotools/protoc-gen-cnw/protoc-gen-cnw.go
*/

/* geninfo:
   filename  : protos/golang.conradwood.net/apis/getestservice/getestservice.proto
   gopackage : golang.conradwood.net/apis/getestservice
   importname: ai_0
   clientfunc: GetEchoService
   serverfunc: NewEchoService
   lookupfunc: EchoServiceLookupID
   varname   : client_EchoServiceClient_0
   clientname: EchoServiceClient
   servername: EchoServiceServer
   gscvname  : getestservice.EchoService
   lockname  : lock_EchoServiceClient_0
   activename: active_EchoServiceClient_0
*/

package getestservice

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

    client_EchoServiceClient_0 = NewEchoServiceClient(client.Connect(EchoServiceLookupID()))
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

    client_EchoServiceClient_0 = NewEchoServiceClient(client.Connect(EchoServiceLookupID()))
    lock_EchoServiceClient_0.Unlock()
    return client_EchoServiceClient_0
}

func EchoServiceLookupID() string { return "getestservice.EchoService" } // returns the ID suitable for lookup in the registry. treat as opaque, subject to change.
