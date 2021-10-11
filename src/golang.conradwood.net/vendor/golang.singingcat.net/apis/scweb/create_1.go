// client create: SCWebServiceClient
/* geninfo:
   filename  : golang.singingcat.net/apis/scweb/scweb.proto
   gopackage : golang.singingcat.net/apis/scweb
   importname: ai_0
   varname   : client_SCWebServiceClient_0
   clientname: SCWebServiceClient
   servername: SCWebServiceServer
   gscvname  : scweb.SCWebService
   lockname  : lock_SCWebServiceClient_0
   activename: active_SCWebServiceClient_0
*/

package scweb

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SCWebServiceClient_0 sync.Mutex
  client_SCWebServiceClient_0 SCWebServiceClient
)

func GetSCWebClient() SCWebServiceClient { 
    if client_SCWebServiceClient_0 != nil {
        return client_SCWebServiceClient_0
    }

    lock_SCWebServiceClient_0.Lock() 
    if client_SCWebServiceClient_0 != nil {
       lock_SCWebServiceClient_0.Unlock()
       return client_SCWebServiceClient_0
    }

    client_SCWebServiceClient_0 = NewSCWebServiceClient(client.Connect("scweb.SCWebService"))
    lock_SCWebServiceClient_0.Unlock()
    return client_SCWebServiceClient_0
}

func GetSCWebServiceClient() SCWebServiceClient { 
    if client_SCWebServiceClient_0 != nil {
        return client_SCWebServiceClient_0
    }

    lock_SCWebServiceClient_0.Lock() 
    if client_SCWebServiceClient_0 != nil {
       lock_SCWebServiceClient_0.Unlock()
       return client_SCWebServiceClient_0
    }

    client_SCWebServiceClient_0 = NewSCWebServiceClient(client.Connect("scweb.SCWebService"))
    lock_SCWebServiceClient_0.Unlock()
    return client_SCWebServiceClient_0
}

