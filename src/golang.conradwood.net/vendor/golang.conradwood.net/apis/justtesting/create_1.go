// client create: JustTestingClient
/* geninfo:
   filename  : golang.conradwood.net/apis/justtesting/justtesting.proto
   gopackage : golang.conradwood.net/apis/justtesting
   importname: ai_0
   varname   : client_JustTestingClient_0
   clientname: JustTestingClient
   servername: JustTestingServer
   gscvname  : justtesting.JustTesting
   lockname  : lock_JustTestingClient_0
   activename: active_JustTestingClient_0
*/

package justtesting

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_JustTestingClient_0 sync.Mutex
  client_JustTestingClient_0 JustTestingClient
)

func GetJustTestingClient() JustTestingClient { 
    if client_JustTestingClient_0 != nil {
        return client_JustTestingClient_0
    }

    lock_JustTestingClient_0.Lock() 
    if client_JustTestingClient_0 != nil {
       lock_JustTestingClient_0.Unlock()
       return client_JustTestingClient_0
    }

    client_JustTestingClient_0 = NewJustTestingClient(client.Connect("justtesting.JustTesting"))
    lock_JustTestingClient_0.Unlock()
    return client_JustTestingClient_0
}

