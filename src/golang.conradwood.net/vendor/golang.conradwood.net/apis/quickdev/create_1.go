// client create: QuickDevServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/quickdev/quickdev.proto
   gopackage : golang.conradwood.net/apis/quickdev
   importname: ai_0
   varname   : client_QuickDevServiceClient_0
   clientname: QuickDevServiceClient
   servername: QuickDevServiceServer
   gscvname  : quickdev.QuickDevService
   lockname  : lock_QuickDevServiceClient_0
   activename: active_QuickDevServiceClient_0
*/

package quickdev

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_QuickDevServiceClient_0 sync.Mutex
  client_QuickDevServiceClient_0 QuickDevServiceClient
)

func GetQuickDevClient() QuickDevServiceClient { 
    if client_QuickDevServiceClient_0 != nil {
        return client_QuickDevServiceClient_0
    }

    lock_QuickDevServiceClient_0.Lock() 
    if client_QuickDevServiceClient_0 != nil {
       lock_QuickDevServiceClient_0.Unlock()
       return client_QuickDevServiceClient_0
    }

    client_QuickDevServiceClient_0 = NewQuickDevServiceClient(client.Connect("quickdev.QuickDevService"))
    lock_QuickDevServiceClient_0.Unlock()
    return client_QuickDevServiceClient_0
}

func GetQuickDevServiceClient() QuickDevServiceClient { 
    if client_QuickDevServiceClient_0 != nil {
        return client_QuickDevServiceClient_0
    }

    lock_QuickDevServiceClient_0.Lock() 
    if client_QuickDevServiceClient_0 != nil {
       lock_QuickDevServiceClient_0.Unlock()
       return client_QuickDevServiceClient_0
    }

    client_QuickDevServiceClient_0 = NewQuickDevServiceClient(client.Connect("quickdev.QuickDevService"))
    lock_QuickDevServiceClient_0.Unlock()
    return client_QuickDevServiceClient_0
}

