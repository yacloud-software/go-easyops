// client create: ExchnageServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/exchnage/exchnage.proto
   gopackage : golang.conradwood.net/apis/exchnage
   importname: ai_0
   varname   : client_ExchnageServiceClient_0
   clientname: ExchnageServiceClient
   servername: ExchnageServiceServer
   gscvname  : exchnage.ExchnageService
   lockname  : lock_ExchnageServiceClient_0
   activename: active_ExchnageServiceClient_0
*/

package exchnage

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ExchnageServiceClient_0 sync.Mutex
  client_ExchnageServiceClient_0 ExchnageServiceClient
)

func GetExchnageClient() ExchnageServiceClient { 
    if client_ExchnageServiceClient_0 != nil {
        return client_ExchnageServiceClient_0
    }

    lock_ExchnageServiceClient_0.Lock() 
    if client_ExchnageServiceClient_0 != nil {
       lock_ExchnageServiceClient_0.Unlock()
       return client_ExchnageServiceClient_0
    }

    client_ExchnageServiceClient_0 = NewExchnageServiceClient(client.Connect("exchnage.ExchnageService"))
    lock_ExchnageServiceClient_0.Unlock()
    return client_ExchnageServiceClient_0
}

func GetExchnageServiceClient() ExchnageServiceClient { 
    if client_ExchnageServiceClient_0 != nil {
        return client_ExchnageServiceClient_0
    }

    lock_ExchnageServiceClient_0.Lock() 
    if client_ExchnageServiceClient_0 != nil {
       lock_ExchnageServiceClient_0.Unlock()
       return client_ExchnageServiceClient_0
    }

    client_ExchnageServiceClient_0 = NewExchnageServiceClient(client.Connect("exchnage.ExchnageService"))
    lock_ExchnageServiceClient_0.Unlock()
    return client_ExchnageServiceClient_0
}

