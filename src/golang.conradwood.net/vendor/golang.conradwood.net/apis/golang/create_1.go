// client create: GolangServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/golang/golang.proto
   gopackage : golang.conradwood.net/apis/golang
   importname: ai_0
   varname   : client_GolangServiceClient_0
   clientname: GolangServiceClient
   servername: GolangServiceServer
   gscvname  : golang.GolangService
   lockname  : lock_GolangServiceClient_0
   activename: active_GolangServiceClient_0
*/

package golang

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_GolangServiceClient_0 sync.Mutex
  client_GolangServiceClient_0 GolangServiceClient
)

func GetGolangClient() GolangServiceClient { 
    if client_GolangServiceClient_0 != nil {
        return client_GolangServiceClient_0
    }

    lock_GolangServiceClient_0.Lock() 
    if client_GolangServiceClient_0 != nil {
       lock_GolangServiceClient_0.Unlock()
       return client_GolangServiceClient_0
    }

    client_GolangServiceClient_0 = NewGolangServiceClient(client.Connect("golang.GolangService"))
    lock_GolangServiceClient_0.Unlock()
    return client_GolangServiceClient_0
}

func GetGolangServiceClient() GolangServiceClient { 
    if client_GolangServiceClient_0 != nil {
        return client_GolangServiceClient_0
    }

    lock_GolangServiceClient_0.Lock() 
    if client_GolangServiceClient_0 != nil {
       lock_GolangServiceClient_0.Unlock()
       return client_GolangServiceClient_0
    }

    client_GolangServiceClient_0 = NewGolangServiceClient(client.Connect("golang.GolangService"))
    lock_GolangServiceClient_0.Unlock()
    return client_GolangServiceClient_0
}

