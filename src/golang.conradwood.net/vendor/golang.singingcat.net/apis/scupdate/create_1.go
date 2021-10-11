// client create: SCUpdateServiceClient
/* geninfo:
   filename  : golang.singingcat.net/apis/scupdate/scupdate.proto
   gopackage : golang.singingcat.net/apis/scupdate
   importname: ai_0
   varname   : client_SCUpdateServiceClient_0
   clientname: SCUpdateServiceClient
   servername: SCUpdateServiceServer
   gscvname  : scupdate.SCUpdateService
   lockname  : lock_SCUpdateServiceClient_0
   activename: active_SCUpdateServiceClient_0
*/

package scupdate

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SCUpdateServiceClient_0 sync.Mutex
  client_SCUpdateServiceClient_0 SCUpdateServiceClient
)

func GetSCUpdateClient() SCUpdateServiceClient { 
    if client_SCUpdateServiceClient_0 != nil {
        return client_SCUpdateServiceClient_0
    }

    lock_SCUpdateServiceClient_0.Lock() 
    if client_SCUpdateServiceClient_0 != nil {
       lock_SCUpdateServiceClient_0.Unlock()
       return client_SCUpdateServiceClient_0
    }

    client_SCUpdateServiceClient_0 = NewSCUpdateServiceClient(client.Connect("scupdate.SCUpdateService"))
    lock_SCUpdateServiceClient_0.Unlock()
    return client_SCUpdateServiceClient_0
}

func GetSCUpdateServiceClient() SCUpdateServiceClient { 
    if client_SCUpdateServiceClient_0 != nil {
        return client_SCUpdateServiceClient_0
    }

    lock_SCUpdateServiceClient_0.Lock() 
    if client_SCUpdateServiceClient_0 != nil {
       lock_SCUpdateServiceClient_0.Unlock()
       return client_SCUpdateServiceClient_0
    }

    client_SCUpdateServiceClient_0 = NewSCUpdateServiceClient(client.Connect("scupdate.SCUpdateService"))
    lock_SCUpdateServiceClient_0.Unlock()
    return client_SCUpdateServiceClient_0
}

