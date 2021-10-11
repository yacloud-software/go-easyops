// client create: CSFServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/csf/csf.proto
   gopackage : golang.conradwood.net/apis/csf
   importname: ai_0
   varname   : client_CSFServiceClient_0
   clientname: CSFServiceClient
   servername: CSFServiceServer
   gscvname  : csf.CSFService
   lockname  : lock_CSFServiceClient_0
   activename: active_CSFServiceClient_0
*/

package csf

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_CSFServiceClient_0 sync.Mutex
  client_CSFServiceClient_0 CSFServiceClient
)

func GetCSFClient() CSFServiceClient { 
    if client_CSFServiceClient_0 != nil {
        return client_CSFServiceClient_0
    }

    lock_CSFServiceClient_0.Lock() 
    if client_CSFServiceClient_0 != nil {
       lock_CSFServiceClient_0.Unlock()
       return client_CSFServiceClient_0
    }

    client_CSFServiceClient_0 = NewCSFServiceClient(client.Connect("csf.CSFService"))
    lock_CSFServiceClient_0.Unlock()
    return client_CSFServiceClient_0
}

func GetCSFServiceClient() CSFServiceClient { 
    if client_CSFServiceClient_0 != nil {
        return client_CSFServiceClient_0
    }

    lock_CSFServiceClient_0.Lock() 
    if client_CSFServiceClient_0 != nil {
       lock_CSFServiceClient_0.Unlock()
       return client_CSFServiceClient_0
    }

    client_CSFServiceClient_0 = NewCSFServiceClient(client.Connect("csf.CSFService"))
    lock_CSFServiceClient_0.Unlock()
    return client_CSFServiceClient_0
}

