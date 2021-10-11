// client create: CallGraphServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/callgraph/callgraph.proto
   gopackage : golang.conradwood.net/apis/callgraph
   importname: ai_0
   varname   : client_CallGraphServiceClient_0
   clientname: CallGraphServiceClient
   servername: CallGraphServiceServer
   gscvname  : callgraph.CallGraphService
   lockname  : lock_CallGraphServiceClient_0
   activename: active_CallGraphServiceClient_0
*/

package callgraph

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_CallGraphServiceClient_0 sync.Mutex
  client_CallGraphServiceClient_0 CallGraphServiceClient
)

func GetCallGraphClient() CallGraphServiceClient { 
    if client_CallGraphServiceClient_0 != nil {
        return client_CallGraphServiceClient_0
    }

    lock_CallGraphServiceClient_0.Lock() 
    if client_CallGraphServiceClient_0 != nil {
       lock_CallGraphServiceClient_0.Unlock()
       return client_CallGraphServiceClient_0
    }

    client_CallGraphServiceClient_0 = NewCallGraphServiceClient(client.Connect("callgraph.CallGraphService"))
    lock_CallGraphServiceClient_0.Unlock()
    return client_CallGraphServiceClient_0
}

func GetCallGraphServiceClient() CallGraphServiceClient { 
    if client_CallGraphServiceClient_0 != nil {
        return client_CallGraphServiceClient_0
    }

    lock_CallGraphServiceClient_0.Lock() 
    if client_CallGraphServiceClient_0 != nil {
       lock_CallGraphServiceClient_0.Unlock()
       return client_CallGraphServiceClient_0
    }

    client_CallGraphServiceClient_0 = NewCallGraphServiceClient(client.Connect("callgraph.CallGraphService"))
    lock_CallGraphServiceClient_0.Unlock()
    return client_CallGraphServiceClient_0
}

