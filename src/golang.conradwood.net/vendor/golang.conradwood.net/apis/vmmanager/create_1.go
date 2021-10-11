// client create: VMManagerClient
/* geninfo:
   filename  : golang.conradwood.net/apis/vmmanager/vmmanager.proto
   gopackage : golang.conradwood.net/apis/vmmanager
   importname: ai_0
   varname   : client_VMManagerClient_0
   clientname: VMManagerClient
   servername: VMManagerServer
   gscvname  : vmmanager.VMManager
   lockname  : lock_VMManagerClient_0
   activename: active_VMManagerClient_0
*/

package vmmanager

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_VMManagerClient_0 sync.Mutex
  client_VMManagerClient_0 VMManagerClient
)

func GetVMManagerClient() VMManagerClient { 
    if client_VMManagerClient_0 != nil {
        return client_VMManagerClient_0
    }

    lock_VMManagerClient_0.Lock() 
    if client_VMManagerClient_0 != nil {
       lock_VMManagerClient_0.Unlock()
       return client_VMManagerClient_0
    }

    client_VMManagerClient_0 = NewVMManagerClient(client.Connect("vmmanager.VMManager"))
    lock_VMManagerClient_0.Unlock()
    return client_VMManagerClient_0
}

