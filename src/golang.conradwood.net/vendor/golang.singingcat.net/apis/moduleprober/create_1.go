// client create: ModuleProberClient
/* geninfo:
   filename  : golang.singingcat.net/apis/moduleprober/moduleprober.proto
   gopackage : golang.singingcat.net/apis/moduleprober
   importname: ai_0
   varname   : client_ModuleProberClient_0
   clientname: ModuleProberClient
   servername: ModuleProberServer
   gscvname  : moduleprober.ModuleProber
   lockname  : lock_ModuleProberClient_0
   activename: active_ModuleProberClient_0
*/

package moduleprober

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ModuleProberClient_0 sync.Mutex
  client_ModuleProberClient_0 ModuleProberClient
)

func GetModuleProberClient() ModuleProberClient { 
    if client_ModuleProberClient_0 != nil {
        return client_ModuleProberClient_0
    }

    lock_ModuleProberClient_0.Lock() 
    if client_ModuleProberClient_0 != nil {
       lock_ModuleProberClient_0.Unlock()
       return client_ModuleProberClient_0
    }

    client_ModuleProberClient_0 = NewModuleProberClient(client.Connect("moduleprober.ModuleProber"))
    lock_ModuleProberClient_0.Unlock()
    return client_ModuleProberClient_0
}

