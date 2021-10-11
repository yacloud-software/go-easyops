// client create: EasyOpsTestClient
/* geninfo:
   filename  : golang.conradwood.net/apis/goeasyops/goeasyops.proto
   gopackage : golang.conradwood.net/apis/goeasyops
   importname: ai_1
   varname   : client_EasyOpsTestClient_1
   clientname: EasyOpsTestClient
   servername: EasyOpsTestServer
   gscvname  : goeasyops.EasyOpsTest
   lockname  : lock_EasyOpsTestClient_1
   activename: active_EasyOpsTestClient_1
*/

package goeasyops

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_EasyOpsTestClient_1 sync.Mutex
  client_EasyOpsTestClient_1 EasyOpsTestClient
)

func GetEasyOpsTestClient() EasyOpsTestClient { 
    if client_EasyOpsTestClient_1 != nil {
        return client_EasyOpsTestClient_1
    }

    lock_EasyOpsTestClient_1.Lock() 
    if client_EasyOpsTestClient_1 != nil {
       lock_EasyOpsTestClient_1.Unlock()
       return client_EasyOpsTestClient_1
    }

    client_EasyOpsTestClient_1 = NewEasyOpsTestClient(client.Connect("goeasyops.EasyOpsTest"))
    lock_EasyOpsTestClient_1.Unlock()
    return client_EasyOpsTestClient_1
}

