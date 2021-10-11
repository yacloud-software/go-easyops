// client create: ApplyClient
/* geninfo:
   filename  : golang.singingcat.net/apis/scapply/scapply.proto
   gopackage : golang.singingcat.net/apis/scapply
   importname: ai_0
   varname   : client_ApplyClient_0
   clientname: ApplyClient
   servername: ApplyServer
   gscvname  : scapply.Apply
   lockname  : lock_ApplyClient_0
   activename: active_ApplyClient_0
*/

package scapply

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ApplyClient_0 sync.Mutex
  client_ApplyClient_0 ApplyClient
)

func GetApplyClient() ApplyClient { 
    if client_ApplyClient_0 != nil {
        return client_ApplyClient_0
    }

    lock_ApplyClient_0.Lock() 
    if client_ApplyClient_0 != nil {
       lock_ApplyClient_0.Unlock()
       return client_ApplyClient_0
    }

    client_ApplyClient_0 = NewApplyClient(client.Connect("scapply.Apply"))
    lock_ApplyClient_0.Unlock()
    return client_ApplyClient_0
}

