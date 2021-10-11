// client create: PublicProxyClient
/* geninfo:
   filename  : golang.conradwood.net/apis/publicproxy/publicproxy.proto
   gopackage : golang.conradwood.net/apis/publicproxy
   importname: ai_0
   varname   : client_PublicProxyClient_0
   clientname: PublicProxyClient
   servername: PublicProxyServer
   gscvname  : publicproxy.PublicProxy
   lockname  : lock_PublicProxyClient_0
   activename: active_PublicProxyClient_0
*/

package publicproxy

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_PublicProxyClient_0 sync.Mutex
  client_PublicProxyClient_0 PublicProxyClient
)

func GetPublicProxyClient() PublicProxyClient { 
    if client_PublicProxyClient_0 != nil {
        return client_PublicProxyClient_0
    }

    lock_PublicProxyClient_0.Lock() 
    if client_PublicProxyClient_0 != nil {
       lock_PublicProxyClient_0.Unlock()
       return client_PublicProxyClient_0
    }

    client_PublicProxyClient_0 = NewPublicProxyClient(client.Connect("publicproxy.PublicProxy"))
    lock_PublicProxyClient_0.Unlock()
    return client_PublicProxyClient_0
}

