// client create: EasyOpsTestClient
/*
  Created by /home/cnw/devel/go/yatools/src/golang.yacloud.eu/yatools/protoc-gen-cnw/protoc-gen-cnw.go
*/

/* geninfo:
   rendererv : 2
   filename  : golang.conradwood.net/apis/getestservice/getestservice.proto
   gopackage : golang.conradwood.net/apis/getestservice
   importname: ai_3
   clientfunc: GetEasyOpsTest
   serverfunc: NewEasyOpsTest
   lookupfunc: EasyOpsTestLookupID
   varname   : client_EasyOpsTestClient_3
   clientname: EasyOpsTestClient
   servername: EasyOpsTestServer
   gsvcname  : getestservice.EasyOpsTest
   lockname  : lock_EasyOpsTestClient_3
   activename: active_EasyOpsTestClient_3
*/

package getestservice

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_EasyOpsTestClient_3 sync.Mutex
  client_EasyOpsTestClient_3 EasyOpsTestClient
)

func GetEasyOpsTestClient() EasyOpsTestClient { 
    if client_EasyOpsTestClient_3 != nil {
        return client_EasyOpsTestClient_3
    }

    lock_EasyOpsTestClient_3.Lock() 
    if client_EasyOpsTestClient_3 != nil {
       lock_EasyOpsTestClient_3.Unlock()
       return client_EasyOpsTestClient_3
    }

    client_EasyOpsTestClient_3 = NewEasyOpsTestClient(client.Connect(EasyOpsTestLookupID()))
    lock_EasyOpsTestClient_3.Unlock()
    return client_EasyOpsTestClient_3
}

func EasyOpsTestLookupID() string { return "getestservice.EasyOpsTest" } // returns the ID suitable for lookup in the registry. treat as opaque, subject to change.

func init() {
   client.RegisterDependency("getestservice.EasyOpsTest")
   AddService("getestservice.EasyOpsTest")
}
