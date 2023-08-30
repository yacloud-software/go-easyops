// client create: CtxTestClient
/*
  Created by /home/cnw/devel/go/yatools/src/golang.yacloud.eu/yatools/protoc-gen-cnw/protoc-gen-cnw.go
*/

/* geninfo:
   filename  : protos/golang.conradwood.net/apis/getestservice/getestservice.proto
   gopackage : golang.conradwood.net/apis/getestservice
   importname: ai_4
   clientfunc: GetCtxTest
   serverfunc: NewCtxTest
   lookupfunc: CtxTestLookupID
   varname   : client_CtxTestClient_4
   clientname: CtxTestClient
   servername: CtxTestServer
   gsvcname  : getestservice.CtxTest
   lockname  : lock_CtxTestClient_4
   activename: active_CtxTestClient_4
*/

package getestservice

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_CtxTestClient_4 sync.Mutex
  client_CtxTestClient_4 CtxTestClient
)

func GetCtxTestClient() CtxTestClient { 
    if client_CtxTestClient_4 != nil {
        return client_CtxTestClient_4
    }

    lock_CtxTestClient_4.Lock() 
    if client_CtxTestClient_4 != nil {
       lock_CtxTestClient_4.Unlock()
       return client_CtxTestClient_4
    }

    client_CtxTestClient_4 = NewCtxTestClient(client.Connect(CtxTestLookupID()))
    lock_CtxTestClient_4.Unlock()
    return client_CtxTestClient_4
}

func CtxTestLookupID() string { return "getestservice.CtxTest" } // returns the ID suitable for lookup in the registry. treat as opaque, subject to change.

func init() {
   client.RegisterDependency("getestservice.CtxTest")
}
