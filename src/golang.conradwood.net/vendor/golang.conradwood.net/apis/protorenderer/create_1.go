// client create: ProtoRendererServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/protorenderer/protorenderer.proto
   gopackage : golang.conradwood.net/apis/protorenderer
   importname: ai_0
   varname   : client_ProtoRendererServiceClient_0
   clientname: ProtoRendererServiceClient
   servername: ProtoRendererServiceServer
   gscvname  : protorenderer.ProtoRendererService
   lockname  : lock_ProtoRendererServiceClient_0
   activename: active_ProtoRendererServiceClient_0
*/

package protorenderer

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ProtoRendererServiceClient_0 sync.Mutex
  client_ProtoRendererServiceClient_0 ProtoRendererServiceClient
)

func GetProtoRendererClient() ProtoRendererServiceClient { 
    if client_ProtoRendererServiceClient_0 != nil {
        return client_ProtoRendererServiceClient_0
    }

    lock_ProtoRendererServiceClient_0.Lock() 
    if client_ProtoRendererServiceClient_0 != nil {
       lock_ProtoRendererServiceClient_0.Unlock()
       return client_ProtoRendererServiceClient_0
    }

    client_ProtoRendererServiceClient_0 = NewProtoRendererServiceClient(client.Connect("protorenderer.ProtoRendererService"))
    lock_ProtoRendererServiceClient_0.Unlock()
    return client_ProtoRendererServiceClient_0
}

func GetProtoRendererServiceClient() ProtoRendererServiceClient { 
    if client_ProtoRendererServiceClient_0 != nil {
        return client_ProtoRendererServiceClient_0
    }

    lock_ProtoRendererServiceClient_0.Lock() 
    if client_ProtoRendererServiceClient_0 != nil {
       lock_ProtoRendererServiceClient_0.Unlock()
       return client_ProtoRendererServiceClient_0
    }

    client_ProtoRendererServiceClient_0 = NewProtoRendererServiceClient(client.Connect("protorenderer.ProtoRendererService"))
    lock_ProtoRendererServiceClient_0.Unlock()
    return client_ProtoRendererServiceClient_0
}

