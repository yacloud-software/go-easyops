// client create: EasyOpsClient
/*
  Created by /home/cnw/devel/go/yatools/src/golang.yacloud.eu/yatools/protoc-gen-cnw/protoc-gen-cnw.go
*/

/* geninfo:
   filename  : protos/golang.conradwood.net/apis/getestservice/getestservice.proto
   gopackage : golang.conradwood.net/apis/getestservice
   importname: ai_2
   clientfunc: GetEasyOps
   serverfunc: NewEasyOps
   lookupfunc: EasyOpsLookupID
   varname   : client_EasyOpsClient_2
   clientname: EasyOpsClient
   servername: EasyOpsServer
   gsvcname  : getestservice.EasyOps
   lockname  : lock_EasyOpsClient_2
   activename: active_EasyOpsClient_2
*/

package getestservice

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_EasyOpsClient_2 sync.Mutex
  client_EasyOpsClient_2 EasyOpsClient
)

func GetEasyOpsClient() EasyOpsClient { 
    if client_EasyOpsClient_2 != nil {
        return client_EasyOpsClient_2
    }

    lock_EasyOpsClient_2.Lock() 
    if client_EasyOpsClient_2 != nil {
       lock_EasyOpsClient_2.Unlock()
       return client_EasyOpsClient_2
    }

    client_EasyOpsClient_2 = NewEasyOpsClient(client.Connect(EasyOpsLookupID()))
    lock_EasyOpsClient_2.Unlock()
    return client_EasyOpsClient_2
}

func EasyOpsLookupID() string { return "getestservice.EasyOps" } // returns the ID suitable for lookup in the registry. treat as opaque, subject to change.

func init() {
   client.RegisterDependency("getestservice.EasyOps")
}
