// client create: SCAutoUpdateClient
/* geninfo:
   filename  : golang.singingcat.net/apis/scautoupdate/scautoupdate.proto
   gopackage : golang.singingcat.net/apis/scautoupdate
   importname: ai_0
   varname   : client_SCAutoUpdateClient_0
   clientname: SCAutoUpdateClient
   servername: SCAutoUpdateServer
   gscvname  : scautoupdate.SCAutoUpdate
   lockname  : lock_SCAutoUpdateClient_0
   activename: active_SCAutoUpdateClient_0
*/

package scautoupdate

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SCAutoUpdateClient_0 sync.Mutex
  client_SCAutoUpdateClient_0 SCAutoUpdateClient
)

func GetSCAutoUpdateClient() SCAutoUpdateClient { 
    if client_SCAutoUpdateClient_0 != nil {
        return client_SCAutoUpdateClient_0
    }

    lock_SCAutoUpdateClient_0.Lock() 
    if client_SCAutoUpdateClient_0 != nil {
       lock_SCAutoUpdateClient_0.Unlock()
       return client_SCAutoUpdateClient_0
    }

    client_SCAutoUpdateClient_0 = NewSCAutoUpdateClient(client.Connect("scautoupdate.SCAutoUpdate"))
    lock_SCAutoUpdateClient_0.Unlock()
    return client_SCAutoUpdateClient_0
}

