// client create: FormRendererClient
/* geninfo:
   filename  : golang.conradwood.net/apis/workflow/workflow.proto
   gopackage : golang.conradwood.net/apis/workflow
   importname: ai_1
   varname   : client_FormRendererClient_1
   clientname: FormRendererClient
   servername: FormRendererServer
   gscvname  : workflow.FormRenderer
   lockname  : lock_FormRendererClient_1
   activename: active_FormRendererClient_1
*/

package workflow

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_FormRendererClient_1 sync.Mutex
  client_FormRendererClient_1 FormRendererClient
)

func GetFormRendererClient() FormRendererClient { 
    if client_FormRendererClient_1 != nil {
        return client_FormRendererClient_1
    }

    lock_FormRendererClient_1.Lock() 
    if client_FormRendererClient_1 != nil {
       lock_FormRendererClient_1.Unlock()
       return client_FormRendererClient_1
    }

    client_FormRendererClient_1 = NewFormRendererClient(client.Connect("workflow.FormRenderer"))
    lock_FormRendererClient_1.Unlock()
    return client_FormRendererClient_1
}

