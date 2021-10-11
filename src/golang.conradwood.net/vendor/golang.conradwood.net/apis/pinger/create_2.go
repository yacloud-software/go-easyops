// client create: PingerListClient
/* geninfo:
   filename  : golang.conradwood.net/apis/pinger/pinger.proto
   gopackage : golang.conradwood.net/apis/pinger
   importname: ai_1
   varname   : client_PingerListClient_1
   clientname: PingerListClient
   servername: PingerListServer
   gscvname  : pinger.PingerList
   lockname  : lock_PingerListClient_1
   activename: active_PingerListClient_1
*/

package pinger

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_PingerListClient_1 sync.Mutex
  client_PingerListClient_1 PingerListClient
)

func GetPingerListClient() PingerListClient { 
    if client_PingerListClient_1 != nil {
        return client_PingerListClient_1
    }

    lock_PingerListClient_1.Lock() 
    if client_PingerListClient_1 != nil {
       lock_PingerListClient_1.Unlock()
       return client_PingerListClient_1
    }

    client_PingerListClient_1 = NewPingerListClient(client.Connect("pinger.PingerList"))
    lock_PingerListClient_1.Unlock()
    return client_PingerListClient_1
}

