// client create: RepoBuilderClient
/* geninfo:
   filename  : golang.conradwood.net/apis/repobuilder/repobuilder.proto
   gopackage : golang.conradwood.net/apis/repobuilder
   importname: ai_0
   varname   : client_RepoBuilderClient_0
   clientname: RepoBuilderClient
   servername: RepoBuilderServer
   gscvname  : repobuilder.RepoBuilder
   lockname  : lock_RepoBuilderClient_0
   activename: active_RepoBuilderClient_0
*/

package repobuilder

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_RepoBuilderClient_0 sync.Mutex
  client_RepoBuilderClient_0 RepoBuilderClient
)

func GetRepoBuilderClient() RepoBuilderClient { 
    if client_RepoBuilderClient_0 != nil {
        return client_RepoBuilderClient_0
    }

    lock_RepoBuilderClient_0.Lock() 
    if client_RepoBuilderClient_0 != nil {
       lock_RepoBuilderClient_0.Unlock()
       return client_RepoBuilderClient_0
    }

    client_RepoBuilderClient_0 = NewRepoBuilderClient(client.Connect("repobuilder.RepoBuilder"))
    lock_RepoBuilderClient_0.Unlock()
    return client_RepoBuilderClient_0
}

