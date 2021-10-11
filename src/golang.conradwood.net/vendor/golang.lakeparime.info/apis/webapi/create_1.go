// client create: WebApiClient
/* geninfo:
   filename  : golang.lakeparime.info/apis/webapi/webapi.proto
   gopackage : golang.lakeparime.info/apis/webapi
   importname: ai_0
   varname   : client_WebApiClient_0
   clientname: WebApiClient
   servername: WebApiServer
   gscvname  : webapi.WebApi
   lockname  : lock_WebApiClient_0
   activename: active_WebApiClient_0
*/

package webapi

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_WebApiClient_0 sync.Mutex
  client_WebApiClient_0 WebApiClient
)

func GetWebApiClient() WebApiClient { 
    if client_WebApiClient_0 != nil {
        return client_WebApiClient_0
    }

    lock_WebApiClient_0.Lock() 
    if client_WebApiClient_0 != nil {
       lock_WebApiClient_0.Unlock()
       return client_WebApiClient_0
    }

    client_WebApiClient_0 = NewWebApiClient(client.Connect("webapi.WebApi"))
    lock_WebApiClient_0.Unlock()
    return client_WebApiClient_0
}

