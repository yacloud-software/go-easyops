// client create: AM43ControllerClient
/* geninfo:
   filename  : conradwood.net/apis/am43controller/am43controller.proto
   gopackage : conradwood.net/apis/am43controller
   importname: ai_0
   varname   : client_AM43ControllerClient_0
   clientname: AM43ControllerClient
   servername: AM43ControllerServer
   gscvname  : am43controller.AM43Controller
   lockname  : lock_AM43ControllerClient_0
   activename: active_AM43ControllerClient_0
*/

package am43controller

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_AM43ControllerClient_0 sync.Mutex
  client_AM43ControllerClient_0 AM43ControllerClient
)

func GetAM43ControllerClient() AM43ControllerClient { 
    if client_AM43ControllerClient_0 != nil {
        return client_AM43ControllerClient_0
    }

    lock_AM43ControllerClient_0.Lock() 
    if client_AM43ControllerClient_0 != nil {
       lock_AM43ControllerClient_0.Unlock()
       return client_AM43ControllerClient_0
    }

    client_AM43ControllerClient_0 = NewAM43ControllerClient(client.Connect("am43controller.AM43Controller"))
    lock_AM43ControllerClient_0.Unlock()
    return client_AM43ControllerClient_0
}

