// client create: WebloginClient
/* geninfo:
   filename  : golang.conradwood.net/apis/weblogin/weblogin.proto
   gopackage : golang.conradwood.net/apis/weblogin
   importname: ai_0
   varname   : client_WebloginClient_0
   clientname: WebloginClient
   servername: WebloginServer
   gscvname  : weblogin.Weblogin
   lockname  : lock_WebloginClient_0
   activename: active_WebloginClient_0
*/

package weblogin

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_WebloginClient_0 sync.Mutex
  client_WebloginClient_0 WebloginClient
)

func GetWebloginClient() WebloginClient { 
    if client_WebloginClient_0 != nil {
        return client_WebloginClient_0
    }

    lock_WebloginClient_0.Lock() 
    if client_WebloginClient_0 != nil {
       lock_WebloginClient_0.Unlock()
       return client_WebloginClient_0
    }

    client_WebloginClient_0 = NewWebloginClient(client.Connect("weblogin.Weblogin"))
    lock_WebloginClient_0.Unlock()
    return client_WebloginClient_0
}

