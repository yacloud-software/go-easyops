// client create: WebClient
/* geninfo:
   filename  : golang.singingcat.net/apis/web/web.proto
   gopackage : golang.singingcat.net/apis/web
   importname: ai_0
   varname   : client_WebClient_0
   clientname: WebClient
   servername: WebServer
   gscvname  : web.Web
   lockname  : lock_WebClient_0
   activename: active_WebClient_0
*/

package web

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_WebClient_0 sync.Mutex
  client_WebClient_0 WebClient
)

func GetWebClient() WebClient { 
    if client_WebClient_0 != nil {
        return client_WebClient_0
    }

    lock_WebClient_0.Lock() 
    if client_WebClient_0 != nil {
       lock_WebClient_0.Unlock()
       return client_WebClient_0
    }

    client_WebClient_0 = NewWebClient(client.Connect("web.Web"))
    lock_WebClient_0.Unlock()
    return client_WebClient_0
}

