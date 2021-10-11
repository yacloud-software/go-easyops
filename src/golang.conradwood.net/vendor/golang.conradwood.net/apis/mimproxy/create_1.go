// client create: MiMProxyServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/mimproxy/mimproxy.proto
   gopackage : golang.conradwood.net/apis/mimproxy
   importname: ai_0
   varname   : client_MiMProxyServiceClient_0
   clientname: MiMProxyServiceClient
   servername: MiMProxyServiceServer
   gscvname  : mimproxy.MiMProxyService
   lockname  : lock_MiMProxyServiceClient_0
   activename: active_MiMProxyServiceClient_0
*/

package mimproxy

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_MiMProxyServiceClient_0 sync.Mutex
  client_MiMProxyServiceClient_0 MiMProxyServiceClient
)

func GetMiMProxyClient() MiMProxyServiceClient { 
    if client_MiMProxyServiceClient_0 != nil {
        return client_MiMProxyServiceClient_0
    }

    lock_MiMProxyServiceClient_0.Lock() 
    if client_MiMProxyServiceClient_0 != nil {
       lock_MiMProxyServiceClient_0.Unlock()
       return client_MiMProxyServiceClient_0
    }

    client_MiMProxyServiceClient_0 = NewMiMProxyServiceClient(client.Connect("mimproxy.MiMProxyService"))
    lock_MiMProxyServiceClient_0.Unlock()
    return client_MiMProxyServiceClient_0
}

func GetMiMProxyServiceClient() MiMProxyServiceClient { 
    if client_MiMProxyServiceClient_0 != nil {
        return client_MiMProxyServiceClient_0
    }

    lock_MiMProxyServiceClient_0.Lock() 
    if client_MiMProxyServiceClient_0 != nil {
       lock_MiMProxyServiceClient_0.Unlock()
       return client_MiMProxyServiceClient_0
    }

    client_MiMProxyServiceClient_0 = NewMiMProxyServiceClient(client.Connect("mimproxy.MiMProxyService"))
    lock_MiMProxyServiceClient_0.Unlock()
    return client_MiMProxyServiceClient_0
}

