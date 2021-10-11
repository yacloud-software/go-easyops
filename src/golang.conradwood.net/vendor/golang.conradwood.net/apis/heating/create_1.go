// client create: HeatingServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/heating/heating.proto
   gopackage : golang.conradwood.net/apis/heating
   importname: ai_0
   varname   : client_HeatingServiceClient_0
   clientname: HeatingServiceClient
   servername: HeatingServiceServer
   gscvname  : heating.HeatingService
   lockname  : lock_HeatingServiceClient_0
   activename: active_HeatingServiceClient_0
*/

package heating

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_HeatingServiceClient_0 sync.Mutex
  client_HeatingServiceClient_0 HeatingServiceClient
)

func GetHeatingClient() HeatingServiceClient { 
    if client_HeatingServiceClient_0 != nil {
        return client_HeatingServiceClient_0
    }

    lock_HeatingServiceClient_0.Lock() 
    if client_HeatingServiceClient_0 != nil {
       lock_HeatingServiceClient_0.Unlock()
       return client_HeatingServiceClient_0
    }

    client_HeatingServiceClient_0 = NewHeatingServiceClient(client.Connect("heating.HeatingService"))
    lock_HeatingServiceClient_0.Unlock()
    return client_HeatingServiceClient_0
}

func GetHeatingServiceClient() HeatingServiceClient { 
    if client_HeatingServiceClient_0 != nil {
        return client_HeatingServiceClient_0
    }

    lock_HeatingServiceClient_0.Lock() 
    if client_HeatingServiceClient_0 != nil {
       lock_HeatingServiceClient_0.Unlock()
       return client_HeatingServiceClient_0
    }

    client_HeatingServiceClient_0 = NewHeatingServiceClient(client.Connect("heating.HeatingService"))
    lock_HeatingServiceClient_0.Unlock()
    return client_HeatingServiceClient_0
}

