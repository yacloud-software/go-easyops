// client create: CtxTestClient
/* geninfo:
   filename  : protos/golang.conradwood.net/apis/getestservice/getestservice.proto
   gopackage : golang.conradwood.net/apis/getestservice
   importname: ai_4
   varname   : client_CtxTestClient_4
   clientname: CtxTestClient
   servername: CtxTestServer
   gscvname  : getestservice.CtxTest
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

    client_CtxTestClient_4 = NewCtxTestClient(client.Connect("getestservice.CtxTest"))
    lock_CtxTestClient_4.Unlock()
    return client_CtxTestClient_4
}

