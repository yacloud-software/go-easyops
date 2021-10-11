// client create: SCModCommsServiceClient
/* geninfo:
   filename  : golang.singingcat.net/apis/scmodcomms/scmodcomms.proto
   gopackage : golang.singingcat.net/apis/scmodcomms
   importname: ai_0
   varname   : client_SCModCommsServiceClient_0
   clientname: SCModCommsServiceClient
   servername: SCModCommsServiceServer
   gscvname  : scmodcomms.SCModCommsService
   lockname  : lock_SCModCommsServiceClient_0
   activename: active_SCModCommsServiceClient_0
*/

package scmodcomms

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SCModCommsServiceClient_0 sync.Mutex
  client_SCModCommsServiceClient_0 SCModCommsServiceClient
)

func GetSCModCommsClient() SCModCommsServiceClient { 
    if client_SCModCommsServiceClient_0 != nil {
        return client_SCModCommsServiceClient_0
    }

    lock_SCModCommsServiceClient_0.Lock() 
    if client_SCModCommsServiceClient_0 != nil {
       lock_SCModCommsServiceClient_0.Unlock()
       return client_SCModCommsServiceClient_0
    }

    client_SCModCommsServiceClient_0 = NewSCModCommsServiceClient(client.Connect("scmodcomms.SCModCommsService"))
    lock_SCModCommsServiceClient_0.Unlock()
    return client_SCModCommsServiceClient_0
}

func GetSCModCommsServiceClient() SCModCommsServiceClient { 
    if client_SCModCommsServiceClient_0 != nil {
        return client_SCModCommsServiceClient_0
    }

    lock_SCModCommsServiceClient_0.Lock() 
    if client_SCModCommsServiceClient_0 != nil {
       lock_SCModCommsServiceClient_0.Unlock()
       return client_SCModCommsServiceClient_0
    }

    client_SCModCommsServiceClient_0 = NewSCModCommsServiceClient(client.Connect("scmodcomms.SCModCommsService"))
    lock_SCModCommsServiceClient_0.Unlock()
    return client_SCModCommsServiceClient_0
}

