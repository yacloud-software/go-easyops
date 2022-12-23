// client create: EasyOpsClient
/* geninfo:
   filename  : protos/golang.conradwood.net/apis/getestservice/echoservice.proto
   gopackage : golang.conradwood.net/apis/getestservice
   importname: ai_2
   varname   : client_EasyOpsClient_2
   clientname: EasyOpsClient
   servername: EasyOpsServer
   gscvname  : getestservice.EasyOps
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

    client_EasyOpsClient_2 = NewEasyOpsClient(client.Connect("getestservice.EasyOps"))
    lock_EasyOpsClient_2.Unlock()
    return client_EasyOpsClient_2
}

