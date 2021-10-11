// client create: AntiDOSClient
/* geninfo:
   filename  : golang.conradwood.net/apis/antidos/antidos.proto
   gopackage : golang.conradwood.net/apis/antidos
   importname: ai_0
   varname   : client_AntiDOSClient_0
   clientname: AntiDOSClient
   servername: AntiDOSServer
   gscvname  : antidos.AntiDOS
   lockname  : lock_AntiDOSClient_0
   activename: active_AntiDOSClient_0
*/

package antidos

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_AntiDOSClient_0 sync.Mutex
  client_AntiDOSClient_0 AntiDOSClient
)

func GetAntiDOSClient() AntiDOSClient { 
    if client_AntiDOSClient_0 != nil {
        return client_AntiDOSClient_0
    }

    lock_AntiDOSClient_0.Lock() 
    if client_AntiDOSClient_0 != nil {
       lock_AntiDOSClient_0.Unlock()
       return client_AntiDOSClient_0
    }

    client_AntiDOSClient_0 = NewAntiDOSClient(client.Connect("antidos.AntiDOS"))
    lock_AntiDOSClient_0.Unlock()
    return client_AntiDOSClient_0
}

