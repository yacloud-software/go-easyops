// client create: UserCommandServiceClient
/* geninfo:
   filename  : golang.singingcat.net/apis/usercommand/usercommand.proto
   gopackage : golang.singingcat.net/apis/usercommand
   importname: ai_0
   varname   : client_UserCommandServiceClient_0
   clientname: UserCommandServiceClient
   servername: UserCommandServiceServer
   gscvname  : usercommand.UserCommandService
   lockname  : lock_UserCommandServiceClient_0
   activename: active_UserCommandServiceClient_0
*/

package usercommand

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_UserCommandServiceClient_0 sync.Mutex
  client_UserCommandServiceClient_0 UserCommandServiceClient
)

func GetUserCommandClient() UserCommandServiceClient { 
    if client_UserCommandServiceClient_0 != nil {
        return client_UserCommandServiceClient_0
    }

    lock_UserCommandServiceClient_0.Lock() 
    if client_UserCommandServiceClient_0 != nil {
       lock_UserCommandServiceClient_0.Unlock()
       return client_UserCommandServiceClient_0
    }

    client_UserCommandServiceClient_0 = NewUserCommandServiceClient(client.Connect("usercommand.UserCommandService"))
    lock_UserCommandServiceClient_0.Unlock()
    return client_UserCommandServiceClient_0
}

func GetUserCommandServiceClient() UserCommandServiceClient { 
    if client_UserCommandServiceClient_0 != nil {
        return client_UserCommandServiceClient_0
    }

    lock_UserCommandServiceClient_0.Lock() 
    if client_UserCommandServiceClient_0 != nil {
       lock_UserCommandServiceClient_0.Unlock()
       return client_UserCommandServiceClient_0
    }

    client_UserCommandServiceClient_0 = NewUserCommandServiceClient(client.Connect("usercommand.UserCommandService"))
    lock_UserCommandServiceClient_0.Unlock()
    return client_UserCommandServiceClient_0
}

