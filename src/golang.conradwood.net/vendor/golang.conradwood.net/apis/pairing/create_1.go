// client create: PairingServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/pairing/pairing.proto
   gopackage : golang.conradwood.net/apis/pairing
   importname: ai_0
   varname   : client_PairingServiceClient_0
   clientname: PairingServiceClient
   servername: PairingServiceServer
   gscvname  : pairing.PairingService
   lockname  : lock_PairingServiceClient_0
   activename: active_PairingServiceClient_0
*/

package pairing

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_PairingServiceClient_0 sync.Mutex
  client_PairingServiceClient_0 PairingServiceClient
)

func GetPairingClient() PairingServiceClient { 
    if client_PairingServiceClient_0 != nil {
        return client_PairingServiceClient_0
    }

    lock_PairingServiceClient_0.Lock() 
    if client_PairingServiceClient_0 != nil {
       lock_PairingServiceClient_0.Unlock()
       return client_PairingServiceClient_0
    }

    client_PairingServiceClient_0 = NewPairingServiceClient(client.Connect("pairing.PairingService"))
    lock_PairingServiceClient_0.Unlock()
    return client_PairingServiceClient_0
}

func GetPairingServiceClient() PairingServiceClient { 
    if client_PairingServiceClient_0 != nil {
        return client_PairingServiceClient_0
    }

    lock_PairingServiceClient_0.Lock() 
    if client_PairingServiceClient_0 != nil {
       lock_PairingServiceClient_0.Unlock()
       return client_PairingServiceClient_0
    }

    client_PairingServiceClient_0 = NewPairingServiceClient(client.Connect("pairing.PairingService"))
    lock_PairingServiceClient_0.Unlock()
    return client_PairingServiceClient_0
}

