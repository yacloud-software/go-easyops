// client create: SecureArgsServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/secureargs/secureargs.proto
   gopackage : golang.conradwood.net/apis/secureargs
   importname: ai_0
   varname   : client_SecureArgsServiceClient_0
   clientname: SecureArgsServiceClient
   servername: SecureArgsServiceServer
   gscvname  : secureargs.SecureArgsService
   lockname  : lock_SecureArgsServiceClient_0
   activename: active_SecureArgsServiceClient_0
*/

package secureargs

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SecureArgsServiceClient_0 sync.Mutex
  client_SecureArgsServiceClient_0 SecureArgsServiceClient
)

func GetSecureArgsClient() SecureArgsServiceClient { 
    if client_SecureArgsServiceClient_0 != nil {
        return client_SecureArgsServiceClient_0
    }

    lock_SecureArgsServiceClient_0.Lock() 
    if client_SecureArgsServiceClient_0 != nil {
       lock_SecureArgsServiceClient_0.Unlock()
       return client_SecureArgsServiceClient_0
    }

    client_SecureArgsServiceClient_0 = NewSecureArgsServiceClient(client.Connect("secureargs.SecureArgsService"))
    lock_SecureArgsServiceClient_0.Unlock()
    return client_SecureArgsServiceClient_0
}

func GetSecureArgsServiceClient() SecureArgsServiceClient { 
    if client_SecureArgsServiceClient_0 != nil {
        return client_SecureArgsServiceClient_0
    }

    lock_SecureArgsServiceClient_0.Lock() 
    if client_SecureArgsServiceClient_0 != nil {
       lock_SecureArgsServiceClient_0.Unlock()
       return client_SecureArgsServiceClient_0
    }

    client_SecureArgsServiceClient_0 = NewSecureArgsServiceClient(client.Connect("secureargs.SecureArgsService"))
    lock_SecureArgsServiceClient_0.Unlock()
    return client_SecureArgsServiceClient_0
}

