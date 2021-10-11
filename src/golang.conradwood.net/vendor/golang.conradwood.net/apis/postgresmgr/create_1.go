// client create: PostgresMgrClient
/* geninfo:
   filename  : golang.conradwood.net/apis/postgresmgr/postgresmgr.proto
   gopackage : golang.conradwood.net/apis/postgresmgr
   importname: ai_0
   varname   : client_PostgresMgrClient_0
   clientname: PostgresMgrClient
   servername: PostgresMgrServer
   gscvname  : postgresmgr.PostgresMgr
   lockname  : lock_PostgresMgrClient_0
   activename: active_PostgresMgrClient_0
*/

package postgresmgr

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_PostgresMgrClient_0 sync.Mutex
  client_PostgresMgrClient_0 PostgresMgrClient
)

func GetPostgresMgrClient() PostgresMgrClient { 
    if client_PostgresMgrClient_0 != nil {
        return client_PostgresMgrClient_0
    }

    lock_PostgresMgrClient_0.Lock() 
    if client_PostgresMgrClient_0 != nil {
       lock_PostgresMgrClient_0.Unlock()
       return client_PostgresMgrClient_0
    }

    client_PostgresMgrClient_0 = NewPostgresMgrClient(client.Connect("postgresmgr.PostgresMgr"))
    lock_PostgresMgrClient_0.Unlock()
    return client_PostgresMgrClient_0
}

