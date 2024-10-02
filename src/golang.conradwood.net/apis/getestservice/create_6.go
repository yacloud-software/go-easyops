// client create: Ctx2TestClient
/*
  Created by /home/cnw/devel/go/yatools/src/golang.yacloud.eu/yatools/protoc-gen-cnw/protoc-gen-cnw.go
*/

/* geninfo:
   filename  : protos/golang.conradwood.net/apis/getestservice/getestservice.proto
   gopackage : golang.conradwood.net/apis/getestservice
   importname: ai_5
   clientfunc: GetCtx2Test
   serverfunc: NewCtx2Test
   lookupfunc: Ctx2TestLookupID
   varname   : client_Ctx2TestClient_5
   clientname: Ctx2TestClient
   servername: Ctx2TestServer
   gsvcname  : getestservice.Ctx2Test
   lockname  : lock_Ctx2TestClient_5
   activename: active_Ctx2TestClient_5
*/

package getestservice

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_Ctx2TestClient_5 sync.Mutex
  client_Ctx2TestClient_5 Ctx2TestClient
)

func GetCtx2TestClient() Ctx2TestClient { 
    if client_Ctx2TestClient_5 != nil {
        return client_Ctx2TestClient_5
    }

    lock_Ctx2TestClient_5.Lock() 
    if client_Ctx2TestClient_5 != nil {
       lock_Ctx2TestClient_5.Unlock()
       return client_Ctx2TestClient_5
    }

    client_Ctx2TestClient_5 = NewCtx2TestClient(client.Connect(Ctx2TestLookupID()))
    lock_Ctx2TestClient_5.Unlock()
    return client_Ctx2TestClient_5
}

func Ctx2TestLookupID() string { return "getestservice.Ctx2Test" } // returns the ID suitable for lookup in the registry. treat as opaque, subject to change.

func init() {
   client.RegisterDependency("getestservice.Ctx2Test")
   AddService("getestservice.Ctx2Test")
}
