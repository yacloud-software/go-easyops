// client create: KicadServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/kicad/kicad.proto
   gopackage : golang.conradwood.net/apis/kicad
   importname: ai_0
   varname   : client_KicadServiceClient_0
   clientname: KicadServiceClient
   servername: KicadServiceServer
   gscvname  : kicad.KicadService
   lockname  : lock_KicadServiceClient_0
   activename: active_KicadServiceClient_0
*/

package kicad

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_KicadServiceClient_0 sync.Mutex
  client_KicadServiceClient_0 KicadServiceClient
)

func GetKicadClient() KicadServiceClient { 
    if client_KicadServiceClient_0 != nil {
        return client_KicadServiceClient_0
    }

    lock_KicadServiceClient_0.Lock() 
    if client_KicadServiceClient_0 != nil {
       lock_KicadServiceClient_0.Unlock()
       return client_KicadServiceClient_0
    }

    client_KicadServiceClient_0 = NewKicadServiceClient(client.Connect("kicad.KicadService"))
    lock_KicadServiceClient_0.Unlock()
    return client_KicadServiceClient_0
}

func GetKicadServiceClient() KicadServiceClient { 
    if client_KicadServiceClient_0 != nil {
        return client_KicadServiceClient_0
    }

    lock_KicadServiceClient_0.Lock() 
    if client_KicadServiceClient_0 != nil {
       lock_KicadServiceClient_0.Unlock()
       return client_KicadServiceClient_0
    }

    client_KicadServiceClient_0 = NewKicadServiceClient(client.Connect("kicad.KicadService"))
    lock_KicadServiceClient_0.Unlock()
    return client_KicadServiceClient_0
}

