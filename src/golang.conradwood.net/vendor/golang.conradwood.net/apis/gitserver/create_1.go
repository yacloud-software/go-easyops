// client create: GIT2Client
/* geninfo:
   filename  : golang.conradwood.net/apis/gitserver/gitserver.proto
   gopackage : golang.conradwood.net/apis/gitserver
   importname: ai_0
   varname   : client_GIT2Client_0
   clientname: GIT2Client
   servername: GIT2Server
   gscvname  : gitserver.GIT2
   lockname  : lock_GIT2Client_0
   activename: active_GIT2Client_0
*/

package gitserver

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_GIT2Client_0 sync.Mutex
  client_GIT2Client_0 GIT2Client
)

func GetGIT2Client() GIT2Client { 
    if client_GIT2Client_0 != nil {
        return client_GIT2Client_0
    }

    lock_GIT2Client_0.Lock() 
    if client_GIT2Client_0 != nil {
       lock_GIT2Client_0.Unlock()
       return client_GIT2Client_0
    }

    client_GIT2Client_0 = NewGIT2Client(client.Connect("gitserver.GIT2"))
    lock_GIT2Client_0.Unlock()
    return client_GIT2Client_0
}

