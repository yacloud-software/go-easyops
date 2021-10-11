// client create: GoAsteriskServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/goasterisk/goasterisk.proto
   gopackage : golang.conradwood.net/apis/goasterisk
   importname: ai_0
   varname   : client_GoAsteriskServiceClient_0
   clientname: GoAsteriskServiceClient
   servername: GoAsteriskServiceServer
   gscvname  : goasterisk.GoAsteriskService
   lockname  : lock_GoAsteriskServiceClient_0
   activename: active_GoAsteriskServiceClient_0
*/

package goasterisk

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_GoAsteriskServiceClient_0 sync.Mutex
  client_GoAsteriskServiceClient_0 GoAsteriskServiceClient
)

func GetGoAsteriskClient() GoAsteriskServiceClient { 
    if client_GoAsteriskServiceClient_0 != nil {
        return client_GoAsteriskServiceClient_0
    }

    lock_GoAsteriskServiceClient_0.Lock() 
    if client_GoAsteriskServiceClient_0 != nil {
       lock_GoAsteriskServiceClient_0.Unlock()
       return client_GoAsteriskServiceClient_0
    }

    client_GoAsteriskServiceClient_0 = NewGoAsteriskServiceClient(client.Connect("goasterisk.GoAsteriskService"))
    lock_GoAsteriskServiceClient_0.Unlock()
    return client_GoAsteriskServiceClient_0
}

func GetGoAsteriskServiceClient() GoAsteriskServiceClient { 
    if client_GoAsteriskServiceClient_0 != nil {
        return client_GoAsteriskServiceClient_0
    }

    lock_GoAsteriskServiceClient_0.Lock() 
    if client_GoAsteriskServiceClient_0 != nil {
       lock_GoAsteriskServiceClient_0.Unlock()
       return client_GoAsteriskServiceClient_0
    }

    client_GoAsteriskServiceClient_0 = NewGoAsteriskServiceClient(client.Connect("goasterisk.GoAsteriskService"))
    lock_GoAsteriskServiceClient_0.Unlock()
    return client_GoAsteriskServiceClient_0
}

