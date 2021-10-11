// client create: ModLogClient
/* geninfo:
   filename  : golang.singingcat.net/apis/modlog/modlog.proto
   gopackage : golang.singingcat.net/apis/modlog
   importname: ai_0
   varname   : client_ModLogClient_0
   clientname: ModLogClient
   servername: ModLogServer
   gscvname  : modlog.ModLog
   lockname  : lock_ModLogClient_0
   activename: active_ModLogClient_0
*/

package modlog

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ModLogClient_0 sync.Mutex
  client_ModLogClient_0 ModLogClient
)

func GetModLogClient() ModLogClient { 
    if client_ModLogClient_0 != nil {
        return client_ModLogClient_0
    }

    lock_ModLogClient_0.Lock() 
    if client_ModLogClient_0 != nil {
       lock_ModLogClient_0.Unlock()
       return client_ModLogClient_0
    }

    client_ModLogClient_0 = NewModLogClient(client.Connect("modlog.ModLog"))
    lock_ModLogClient_0.Unlock()
    return client_ModLogClient_0
}

