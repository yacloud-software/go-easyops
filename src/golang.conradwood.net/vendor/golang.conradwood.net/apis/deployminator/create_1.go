// client create: DeployminatorClient
/* geninfo:
   filename  : golang.conradwood.net/apis/deployminator/deployminator.proto
   gopackage : golang.conradwood.net/apis/deployminator
   importname: ai_0
   varname   : client_DeployminatorClient_0
   clientname: DeployminatorClient
   servername: DeployminatorServer
   gscvname  : deployminator.Deployminator
   lockname  : lock_DeployminatorClient_0
   activename: active_DeployminatorClient_0
*/

package deployminator

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_DeployminatorClient_0 sync.Mutex
  client_DeployminatorClient_0 DeployminatorClient
)

func GetDeployminatorClient() DeployminatorClient { 
    if client_DeployminatorClient_0 != nil {
        return client_DeployminatorClient_0
    }

    lock_DeployminatorClient_0.Lock() 
    if client_DeployminatorClient_0 != nil {
       lock_DeployminatorClient_0.Unlock()
       return client_DeployminatorClient_0
    }

    client_DeployminatorClient_0 = NewDeployminatorClient(client.Connect("deployminator.Deployminator"))
    lock_DeployminatorClient_0.Unlock()
    return client_DeployminatorClient_0
}

