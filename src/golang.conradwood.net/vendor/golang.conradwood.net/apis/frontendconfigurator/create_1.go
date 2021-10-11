// client create: FrontEndConfiguratorServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/frontendconfigurator/frontendconfigurator.proto
   gopackage : golang.conradwood.net/apis/frontendconfigurator
   importname: ai_0
   varname   : client_FrontEndConfiguratorServiceClient_0
   clientname: FrontEndConfiguratorServiceClient
   servername: FrontEndConfiguratorServiceServer
   gscvname  : frontendconfigurator.FrontEndConfiguratorService
   lockname  : lock_FrontEndConfiguratorServiceClient_0
   activename: active_FrontEndConfiguratorServiceClient_0
*/

package frontendconfigurator

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_FrontEndConfiguratorServiceClient_0 sync.Mutex
  client_FrontEndConfiguratorServiceClient_0 FrontEndConfiguratorServiceClient
)

func GetFrontEndConfiguratorClient() FrontEndConfiguratorServiceClient { 
    if client_FrontEndConfiguratorServiceClient_0 != nil {
        return client_FrontEndConfiguratorServiceClient_0
    }

    lock_FrontEndConfiguratorServiceClient_0.Lock() 
    if client_FrontEndConfiguratorServiceClient_0 != nil {
       lock_FrontEndConfiguratorServiceClient_0.Unlock()
       return client_FrontEndConfiguratorServiceClient_0
    }

    client_FrontEndConfiguratorServiceClient_0 = NewFrontEndConfiguratorServiceClient(client.Connect("frontendconfigurator.FrontEndConfiguratorService"))
    lock_FrontEndConfiguratorServiceClient_0.Unlock()
    return client_FrontEndConfiguratorServiceClient_0
}

func GetFrontEndConfiguratorServiceClient() FrontEndConfiguratorServiceClient { 
    if client_FrontEndConfiguratorServiceClient_0 != nil {
        return client_FrontEndConfiguratorServiceClient_0
    }

    lock_FrontEndConfiguratorServiceClient_0.Lock() 
    if client_FrontEndConfiguratorServiceClient_0 != nil {
       lock_FrontEndConfiguratorServiceClient_0.Unlock()
       return client_FrontEndConfiguratorServiceClient_0
    }

    client_FrontEndConfiguratorServiceClient_0 = NewFrontEndConfiguratorServiceClient(client.Connect("frontendconfigurator.FrontEndConfiguratorService"))
    lock_FrontEndConfiguratorServiceClient_0.Unlock()
    return client_FrontEndConfiguratorServiceClient_0
}

