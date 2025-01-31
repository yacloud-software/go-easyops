// client create: EchoStreamServiceClient
/*
  Created by /home/cnw/devel/go/yatools/src/golang.yacloud.eu/yatools/protoc-gen-cnw/protoc-gen-cnw.go
*/

/* geninfo:
   rendererv : 2
   filename  : golang.conradwood.net/apis/getestservice/getestservice.proto
   gopackage : golang.conradwood.net/apis/getestservice
   importname: ai_1
   clientfunc: GetEchoStreamService
   serverfunc: NewEchoStreamService
   lookupfunc: EchoStreamServiceLookupID
   varname   : client_EchoStreamServiceClient_1
   clientname: EchoStreamServiceClient
   servername: EchoStreamServiceServer
   gsvcname  : getestservice.EchoStreamService
   lockname  : lock_EchoStreamServiceClient_1
   activename: active_EchoStreamServiceClient_1
*/

package getestservice

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_EchoStreamServiceClient_1 sync.Mutex
  client_EchoStreamServiceClient_1 EchoStreamServiceClient
)

func GetEchoStreamClient() EchoStreamServiceClient { 
    if client_EchoStreamServiceClient_1 != nil {
        return client_EchoStreamServiceClient_1
    }

    lock_EchoStreamServiceClient_1.Lock() 
    if client_EchoStreamServiceClient_1 != nil {
       lock_EchoStreamServiceClient_1.Unlock()
       return client_EchoStreamServiceClient_1
    }

    client_EchoStreamServiceClient_1 = NewEchoStreamServiceClient(client.Connect(EchoStreamServiceLookupID()))
    lock_EchoStreamServiceClient_1.Unlock()
    return client_EchoStreamServiceClient_1
}

func GetEchoStreamServiceClient() EchoStreamServiceClient { 
    if client_EchoStreamServiceClient_1 != nil {
        return client_EchoStreamServiceClient_1
    }

    lock_EchoStreamServiceClient_1.Lock() 
    if client_EchoStreamServiceClient_1 != nil {
       lock_EchoStreamServiceClient_1.Unlock()
       return client_EchoStreamServiceClient_1
    }

    client_EchoStreamServiceClient_1 = NewEchoStreamServiceClient(client.Connect(EchoStreamServiceLookupID()))
    lock_EchoStreamServiceClient_1.Unlock()
    return client_EchoStreamServiceClient_1
}

func EchoStreamServiceLookupID() string { return "getestservice.EchoStreamService" } // returns the ID suitable for lookup in the registry. treat as opaque, subject to change.

func init() {
   client.RegisterDependency("getestservice.EchoStreamService")
   AddService("getestservice.EchoStreamService")
}
