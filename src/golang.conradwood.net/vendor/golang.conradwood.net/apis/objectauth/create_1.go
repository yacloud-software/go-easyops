// client create: ObjectAuthServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/objectauth/objectauth.proto
   gopackage : golang.conradwood.net/apis/objectauth
   importname: ai_0
   varname   : client_ObjectAuthServiceClient_0
   clientname: ObjectAuthServiceClient
   servername: ObjectAuthServiceServer
   gscvname  : objectauth.ObjectAuthService
   lockname  : lock_ObjectAuthServiceClient_0
   activename: active_ObjectAuthServiceClient_0
*/

package objectauth

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ObjectAuthServiceClient_0 sync.Mutex
  client_ObjectAuthServiceClient_0 ObjectAuthServiceClient
)

func GetObjectAuthClient() ObjectAuthServiceClient { 
    if client_ObjectAuthServiceClient_0 != nil {
        return client_ObjectAuthServiceClient_0
    }

    lock_ObjectAuthServiceClient_0.Lock() 
    if client_ObjectAuthServiceClient_0 != nil {
       lock_ObjectAuthServiceClient_0.Unlock()
       return client_ObjectAuthServiceClient_0
    }

    client_ObjectAuthServiceClient_0 = NewObjectAuthServiceClient(client.Connect("objectauth.ObjectAuthService"))
    lock_ObjectAuthServiceClient_0.Unlock()
    return client_ObjectAuthServiceClient_0
}

func GetObjectAuthServiceClient() ObjectAuthServiceClient { 
    if client_ObjectAuthServiceClient_0 != nil {
        return client_ObjectAuthServiceClient_0
    }

    lock_ObjectAuthServiceClient_0.Lock() 
    if client_ObjectAuthServiceClient_0 != nil {
       lock_ObjectAuthServiceClient_0.Unlock()
       return client_ObjectAuthServiceClient_0
    }

    client_ObjectAuthServiceClient_0 = NewObjectAuthServiceClient(client.Connect("objectauth.ObjectAuthService"))
    lock_ObjectAuthServiceClient_0.Unlock()
    return client_ObjectAuthServiceClient_0
}

