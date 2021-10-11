// client create: PingerClient
/* geninfo:
   filename  : golang.conradwood.net/apis/pinger/pinger.proto
   gopackage : golang.conradwood.net/apis/pinger
   importname: ai_0
   varname   : client_PingerClient_0
   clientname: PingerClient
   servername: PingerServer
   gscvname  : pinger.Pinger
   lockname  : lock_PingerClient_0
   activename: active_PingerClient_0
*/

package pinger

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_PingerClient_0 sync.Mutex
  client_PingerClient_0 PingerClient
)

func GetPingerClient() PingerClient { 
    if client_PingerClient_0 != nil {
        return client_PingerClient_0
    }

    lock_PingerClient_0.Lock() 
    if client_PingerClient_0 != nil {
       lock_PingerClient_0.Unlock()
       return client_PingerClient_0
    }

    client_PingerClient_0 = NewPingerClient(client.Connect("pinger.Pinger"))
    lock_PingerClient_0.Unlock()
    return client_PingerClient_0
}

