// client create: JavaRepoClient
/* geninfo:
   filename  : golang.conradwood.net/apis/javarepo/javarepo.proto
   gopackage : golang.conradwood.net/apis/javarepo
   importname: ai_0
   varname   : client_JavaRepoClient_0
   clientname: JavaRepoClient
   servername: JavaRepoServer
   gscvname  : javarepo.JavaRepo
   lockname  : lock_JavaRepoClient_0
   activename: active_JavaRepoClient_0
*/

package javarepo

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_JavaRepoClient_0 sync.Mutex
  client_JavaRepoClient_0 JavaRepoClient
)

func GetJavaRepoClient() JavaRepoClient { 
    if client_JavaRepoClient_0 != nil {
        return client_JavaRepoClient_0
    }

    lock_JavaRepoClient_0.Lock() 
    if client_JavaRepoClient_0 != nil {
       lock_JavaRepoClient_0.Unlock()
       return client_JavaRepoClient_0
    }

    client_JavaRepoClient_0 = NewJavaRepoClient(client.Connect("javarepo.JavaRepo"))
    lock_JavaRepoClient_0.Unlock()
    return client_JavaRepoClient_0
}

