// client create: RegistryMultiplexerServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/registrymultiplexer/registrymultiplexer.proto
   gopackage : golang.conradwood.net/apis/registrymultiplexer
   importname: ai_0
   varname   : client_RegistryMultiplexerServiceClient_0
   clientname: RegistryMultiplexerServiceClient
   servername: RegistryMultiplexerServiceServer
   gscvname  : registrymultiplexer.RegistryMultiplexerService
   lockname  : lock_RegistryMultiplexerServiceClient_0
   activename: active_RegistryMultiplexerServiceClient_0
*/

package registrymultiplexer

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_RegistryMultiplexerServiceClient_0 sync.Mutex
  client_RegistryMultiplexerServiceClient_0 RegistryMultiplexerServiceClient
)

func GetRegistryMultiplexerClient() RegistryMultiplexerServiceClient { 
    if client_RegistryMultiplexerServiceClient_0 != nil {
        return client_RegistryMultiplexerServiceClient_0
    }

    lock_RegistryMultiplexerServiceClient_0.Lock() 
    if client_RegistryMultiplexerServiceClient_0 != nil {
       lock_RegistryMultiplexerServiceClient_0.Unlock()
       return client_RegistryMultiplexerServiceClient_0
    }

    client_RegistryMultiplexerServiceClient_0 = NewRegistryMultiplexerServiceClient(client.Connect("registrymultiplexer.RegistryMultiplexerService"))
    lock_RegistryMultiplexerServiceClient_0.Unlock()
    return client_RegistryMultiplexerServiceClient_0
}

func GetRegistryMultiplexerServiceClient() RegistryMultiplexerServiceClient { 
    if client_RegistryMultiplexerServiceClient_0 != nil {
        return client_RegistryMultiplexerServiceClient_0
    }

    lock_RegistryMultiplexerServiceClient_0.Lock() 
    if client_RegistryMultiplexerServiceClient_0 != nil {
       lock_RegistryMultiplexerServiceClient_0.Unlock()
       return client_RegistryMultiplexerServiceClient_0
    }

    client_RegistryMultiplexerServiceClient_0 = NewRegistryMultiplexerServiceClient(client.Connect("registrymultiplexer.RegistryMultiplexerService"))
    lock_RegistryMultiplexerServiceClient_0.Unlock()
    return client_RegistryMultiplexerServiceClient_0
}

