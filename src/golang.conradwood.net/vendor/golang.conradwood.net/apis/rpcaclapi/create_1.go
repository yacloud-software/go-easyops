// client create: RPCACLServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/rpcaclapi/rpcaclapi.proto
   gopackage : golang.conradwood.net/apis/rpcaclapi
   importname: ai_0
   varname   : client_RPCACLServiceClient_0
   clientname: RPCACLServiceClient
   servername: RPCACLServiceServer
   gscvname  : rpcaclapi.RPCACLService
   lockname  : lock_RPCACLServiceClient_0
   activename: active_RPCACLServiceClient_0
*/

package rpcaclapi

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_RPCACLServiceClient_0 sync.Mutex
  client_RPCACLServiceClient_0 RPCACLServiceClient
)

func GetRPCACLClient() RPCACLServiceClient { 
    if client_RPCACLServiceClient_0 != nil {
        return client_RPCACLServiceClient_0
    }

    lock_RPCACLServiceClient_0.Lock() 
    if client_RPCACLServiceClient_0 != nil {
       lock_RPCACLServiceClient_0.Unlock()
       return client_RPCACLServiceClient_0
    }

    client_RPCACLServiceClient_0 = NewRPCACLServiceClient(client.Connect("rpcaclapi.RPCACLService"))
    lock_RPCACLServiceClient_0.Unlock()
    return client_RPCACLServiceClient_0
}

func GetRPCACLServiceClient() RPCACLServiceClient { 
    if client_RPCACLServiceClient_0 != nil {
        return client_RPCACLServiceClient_0
    }

    lock_RPCACLServiceClient_0.Lock() 
    if client_RPCACLServiceClient_0 != nil {
       lock_RPCACLServiceClient_0.Unlock()
       return client_RPCACLServiceClient_0
    }

    client_RPCACLServiceClient_0 = NewRPCACLServiceClient(client.Connect("rpcaclapi.RPCACLService"))
    lock_RPCACLServiceClient_0.Unlock()
    return client_RPCACLServiceClient_0
}

