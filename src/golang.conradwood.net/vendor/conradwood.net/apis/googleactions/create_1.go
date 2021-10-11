// client create: GoogleActionsClient
/* geninfo:
   filename  : conradwood.net/apis/googleactions/googleactions.proto
   gopackage : conradwood.net/apis/googleactions
   importname: ai_0
   varname   : client_GoogleActionsClient_0
   clientname: GoogleActionsClient
   servername: GoogleActionsServer
   gscvname  : googleactions.GoogleActions
   lockname  : lock_GoogleActionsClient_0
   activename: active_GoogleActionsClient_0
*/

package googleactions

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_GoogleActionsClient_0 sync.Mutex
  client_GoogleActionsClient_0 GoogleActionsClient
)

func GetGoogleActionsClient() GoogleActionsClient { 
    if client_GoogleActionsClient_0 != nil {
        return client_GoogleActionsClient_0
    }

    lock_GoogleActionsClient_0.Lock() 
    if client_GoogleActionsClient_0 != nil {
       lock_GoogleActionsClient_0.Unlock()
       return client_GoogleActionsClient_0
    }

    client_GoogleActionsClient_0 = NewGoogleActionsClient(client.Connect("googleactions.GoogleActions"))
    lock_GoogleActionsClient_0.Unlock()
    return client_GoogleActionsClient_0
}

