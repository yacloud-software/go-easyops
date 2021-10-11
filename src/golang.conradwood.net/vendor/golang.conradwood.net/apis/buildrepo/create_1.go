// client create: BuildRepoManagerClient
/* geninfo:
   filename  : golang.conradwood.net/apis/buildrepo/buildrepo.proto
   gopackage : golang.conradwood.net/apis/buildrepo
   importname: ai_0
   varname   : client_BuildRepoManagerClient_0
   clientname: BuildRepoManagerClient
   servername: BuildRepoManagerServer
   gscvname  : buildrepo.BuildRepoManager
   lockname  : lock_BuildRepoManagerClient_0
   activename: active_BuildRepoManagerClient_0
*/

package buildrepo

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_BuildRepoManagerClient_0 sync.Mutex
  client_BuildRepoManagerClient_0 BuildRepoManagerClient
)

func GetBuildRepoManagerClient() BuildRepoManagerClient { 
    if client_BuildRepoManagerClient_0 != nil {
        return client_BuildRepoManagerClient_0
    }

    lock_BuildRepoManagerClient_0.Lock() 
    if client_BuildRepoManagerClient_0 != nil {
       lock_BuildRepoManagerClient_0.Unlock()
       return client_BuildRepoManagerClient_0
    }

    client_BuildRepoManagerClient_0 = NewBuildRepoManagerClient(client.Connect("buildrepo.BuildRepoManager"))
    lock_BuildRepoManagerClient_0.Unlock()
    return client_BuildRepoManagerClient_0
}

