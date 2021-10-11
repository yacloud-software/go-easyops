// client create: IFTTTServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/ifttt/ifttt.proto
   gopackage : golang.conradwood.net/apis/ifttt
   importname: ai_0
   varname   : client_IFTTTServiceClient_0
   clientname: IFTTTServiceClient
   servername: IFTTTServiceServer
   gscvname  : ifttt.IFTTTService
   lockname  : lock_IFTTTServiceClient_0
   activename: active_IFTTTServiceClient_0
*/

package ifttt

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_IFTTTServiceClient_0 sync.Mutex
  client_IFTTTServiceClient_0 IFTTTServiceClient
)

func GetIFTTTClient() IFTTTServiceClient { 
    if client_IFTTTServiceClient_0 != nil {
        return client_IFTTTServiceClient_0
    }

    lock_IFTTTServiceClient_0.Lock() 
    if client_IFTTTServiceClient_0 != nil {
       lock_IFTTTServiceClient_0.Unlock()
       return client_IFTTTServiceClient_0
    }

    client_IFTTTServiceClient_0 = NewIFTTTServiceClient(client.Connect("ifttt.IFTTTService"))
    lock_IFTTTServiceClient_0.Unlock()
    return client_IFTTTServiceClient_0
}

func GetIFTTTServiceClient() IFTTTServiceClient { 
    if client_IFTTTServiceClient_0 != nil {
        return client_IFTTTServiceClient_0
    }

    lock_IFTTTServiceClient_0.Lock() 
    if client_IFTTTServiceClient_0 != nil {
       lock_IFTTTServiceClient_0.Unlock()
       return client_IFTTTServiceClient_0
    }

    client_IFTTTServiceClient_0 = NewIFTTTServiceClient(client.Connect("ifttt.IFTTTService"))
    lock_IFTTTServiceClient_0.Unlock()
    return client_IFTTTServiceClient_0
}

