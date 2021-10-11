// client create: EasyOpsClient
/* geninfo:
   filename  : golang.conradwood.net/apis/goeasyops/goeasyops.proto
   gopackage : golang.conradwood.net/apis/goeasyops
   importname: ai_0
   varname   : client_EasyOpsClient_0
   clientname: EasyOpsClient
   servername: EasyOpsServer
   gscvname  : goeasyops.EasyOps
   lockname  : lock_EasyOpsClient_0
   activename: active_EasyOpsClient_0
*/

package goeasyops

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_EasyOpsClient_0 sync.Mutex
  client_EasyOpsClient_0 EasyOpsClient
)

func GetEasyOpsClient() EasyOpsClient { 
    if client_EasyOpsClient_0 != nil {
        return client_EasyOpsClient_0
    }

    lock_EasyOpsClient_0.Lock() 
    if client_EasyOpsClient_0 != nil {
       lock_EasyOpsClient_0.Unlock()
       return client_EasyOpsClient_0
    }

    client_EasyOpsClient_0 = NewEasyOpsClient(client.Connect("goeasyops.EasyOps"))
    lock_EasyOpsClient_0.Unlock()
    return client_EasyOpsClient_0
}

