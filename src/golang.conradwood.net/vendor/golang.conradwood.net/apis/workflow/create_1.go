// client create: WorkFlowClient
/* geninfo:
   filename  : golang.conradwood.net/apis/workflow/workflow.proto
   gopackage : golang.conradwood.net/apis/workflow
   importname: ai_0
   varname   : client_WorkFlowClient_0
   clientname: WorkFlowClient
   servername: WorkFlowServer
   gscvname  : workflow.WorkFlow
   lockname  : lock_WorkFlowClient_0
   activename: active_WorkFlowClient_0
*/

package workflow

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_WorkFlowClient_0 sync.Mutex
  client_WorkFlowClient_0 WorkFlowClient
)

func GetWorkFlowClient() WorkFlowClient { 
    if client_WorkFlowClient_0 != nil {
        return client_WorkFlowClient_0
    }

    lock_WorkFlowClient_0.Lock() 
    if client_WorkFlowClient_0 != nil {
       lock_WorkFlowClient_0.Unlock()
       return client_WorkFlowClient_0
    }

    client_WorkFlowClient_0 = NewWorkFlowClient(client.Connect("workflow.WorkFlow"))
    lock_WorkFlowClient_0.Unlock()
    return client_WorkFlowClient_0
}

