// client create: OTAModuleClient
/* geninfo:
   filename  : golang.conradwood.net/apis/otamodule/otamodule.proto
   gopackage : golang.conradwood.net/apis/otamodule
   importname: ai_0
   varname   : client_OTAModuleClient_0
   clientname: OTAModuleClient
   servername: OTAModuleServer
   gscvname  : otamodule.OTAModule
   lockname  : lock_OTAModuleClient_0
   activename: active_OTAModuleClient_0
*/

package otamodule

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_OTAModuleClient_0 sync.Mutex
  client_OTAModuleClient_0 OTAModuleClient
)

func GetOTAModuleClient() OTAModuleClient { 
    if client_OTAModuleClient_0 != nil {
        return client_OTAModuleClient_0
    }

    lock_OTAModuleClient_0.Lock() 
    if client_OTAModuleClient_0 != nil {
       lock_OTAModuleClient_0.Unlock()
       return client_OTAModuleClient_0
    }

    client_OTAModuleClient_0 = NewOTAModuleClient(client.Connect("otamodule.OTAModule"))
    lock_OTAModuleClient_0.Unlock()
    return client_OTAModuleClient_0
}

