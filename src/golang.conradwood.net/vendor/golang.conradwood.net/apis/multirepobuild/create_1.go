// client create: MultiRepoBuildClient
/* geninfo:
   filename  : golang.conradwood.net/apis/multirepobuild/multirepobuild.proto
   gopackage : golang.conradwood.net/apis/multirepobuild
   importname: ai_0
   varname   : client_MultiRepoBuildClient_0
   clientname: MultiRepoBuildClient
   servername: MultiRepoBuildServer
   gscvname  : multirepobuild.MultiRepoBuild
   lockname  : lock_MultiRepoBuildClient_0
   activename: active_MultiRepoBuildClient_0
*/

package multirepobuild

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_MultiRepoBuildClient_0 sync.Mutex
  client_MultiRepoBuildClient_0 MultiRepoBuildClient
)

func GetMultiRepoBuildClient() MultiRepoBuildClient { 
    if client_MultiRepoBuildClient_0 != nil {
        return client_MultiRepoBuildClient_0
    }

    lock_MultiRepoBuildClient_0.Lock() 
    if client_MultiRepoBuildClient_0 != nil {
       lock_MultiRepoBuildClient_0.Unlock()
       return client_MultiRepoBuildClient_0
    }

    client_MultiRepoBuildClient_0 = NewMultiRepoBuildClient(client.Connect("multirepobuild.MultiRepoBuild"))
    lock_MultiRepoBuildClient_0.Unlock()
    return client_MultiRepoBuildClient_0
}

