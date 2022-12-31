// client create: EasyOpsTestClient
/* geninfo:
   filename  : protos/golang.conradwood.net/apis/getestservice/getestservice.proto
   gopackage : golang.conradwood.net/apis/getestservice
   importname: ai_3
   varname   : client_EasyOpsTestClient_3
   clientname: EasyOpsTestClient
   servername: EasyOpsTestServer
   gscvname  : getestservice.EasyOpsTest
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

    client_EasyOpsTestClient_3 = NewEasyOpsTestClient(client.Connect("getestservice.EasyOpsTest"))
    lock_EasyOpsTestClient_3.Unlock()
    return client_EasyOpsTestClient_3
}

