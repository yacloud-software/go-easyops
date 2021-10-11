// client create: HTMLServerServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/htmlserver/htmlserver.proto
   gopackage : golang.conradwood.net/apis/htmlserver
   importname: ai_0
   varname   : client_HTMLServerServiceClient_0
   clientname: HTMLServerServiceClient
   servername: HTMLServerServiceServer
   gscvname  : htmlserver.HTMLServerService
   lockname  : lock_HTMLServerServiceClient_0
   activename: active_HTMLServerServiceClient_0
*/

package htmlserver

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_HTMLServerServiceClient_0 sync.Mutex
  client_HTMLServerServiceClient_0 HTMLServerServiceClient
)

func GetHTMLServerClient() HTMLServerServiceClient { 
    if client_HTMLServerServiceClient_0 != nil {
        return client_HTMLServerServiceClient_0
    }

    lock_HTMLServerServiceClient_0.Lock() 
    if client_HTMLServerServiceClient_0 != nil {
       lock_HTMLServerServiceClient_0.Unlock()
       return client_HTMLServerServiceClient_0
    }

    client_HTMLServerServiceClient_0 = NewHTMLServerServiceClient(client.Connect("htmlserver.HTMLServerService"))
    lock_HTMLServerServiceClient_0.Unlock()
    return client_HTMLServerServiceClient_0
}

func GetHTMLServerServiceClient() HTMLServerServiceClient { 
    if client_HTMLServerServiceClient_0 != nil {
        return client_HTMLServerServiceClient_0
    }

    lock_HTMLServerServiceClient_0.Lock() 
    if client_HTMLServerServiceClient_0 != nil {
       lock_HTMLServerServiceClient_0.Unlock()
       return client_HTMLServerServiceClient_0
    }

    client_HTMLServerServiceClient_0 = NewHTMLServerServiceClient(client.Connect("htmlserver.HTMLServerService"))
    lock_HTMLServerServiceClient_0.Unlock()
    return client_HTMLServerServiceClient_0
}

