// client create: SCGoLibClient
/* geninfo:
   filename  : golang.singingcat.net/apis/scgolib/scgolib.proto
   gopackage : golang.singingcat.net/apis/scgolib
   importname: ai_0
   varname   : client_SCGoLibClient_0
   clientname: SCGoLibClient
   servername: SCGoLibServer
   gscvname  : scgolib.SCGoLib
   lockname  : lock_SCGoLibClient_0
   activename: active_SCGoLibClient_0
*/

package scgolib

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SCGoLibClient_0 sync.Mutex
  client_SCGoLibClient_0 SCGoLibClient
)

func GetSCGoLibClient() SCGoLibClient { 
    if client_SCGoLibClient_0 != nil {
        return client_SCGoLibClient_0
    }

    lock_SCGoLibClient_0.Lock() 
    if client_SCGoLibClient_0 != nil {
       lock_SCGoLibClient_0.Unlock()
       return client_SCGoLibClient_0
    }

    client_SCGoLibClient_0 = NewSCGoLibClient(client.Connect("scgolib.SCGoLib"))
    lock_SCGoLibClient_0.Unlock()
    return client_SCGoLibClient_0
}

