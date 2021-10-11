// client create: UnixAuthClient
/* geninfo:
   filename  : golang.conradwood.net/apis/unixauth/unixauth.proto
   gopackage : golang.conradwood.net/apis/unixauth
   importname: ai_0
   varname   : client_UnixAuthClient_0
   clientname: UnixAuthClient
   servername: UnixAuthServer
   gscvname  : unixauth.UnixAuth
   lockname  : lock_UnixAuthClient_0
   activename: active_UnixAuthClient_0
*/

package unixauth

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_UnixAuthClient_0 sync.Mutex
  client_UnixAuthClient_0 UnixAuthClient
)

func GetUnixAuthClient() UnixAuthClient { 
    if client_UnixAuthClient_0 != nil {
        return client_UnixAuthClient_0
    }

    lock_UnixAuthClient_0.Lock() 
    if client_UnixAuthClient_0 != nil {
       lock_UnixAuthClient_0.Unlock()
       return client_UnixAuthClient_0
    }

    client_UnixAuthClient_0 = NewUnixAuthClient(client.Connect("unixauth.UnixAuth"))
    lock_UnixAuthClient_0.Unlock()
    return client_UnixAuthClient_0
}

