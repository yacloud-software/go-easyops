// client create: ArtefactServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/artefact/artefact.proto
   gopackage : golang.conradwood.net/apis/artefact
   importname: ai_0
   varname   : client_ArtefactServiceClient_0
   clientname: ArtefactServiceClient
   servername: ArtefactServiceServer
   gscvname  : artefact.ArtefactService
   lockname  : lock_ArtefactServiceClient_0
   activename: active_ArtefactServiceClient_0
*/

package artefact

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ArtefactServiceClient_0 sync.Mutex
  client_ArtefactServiceClient_0 ArtefactServiceClient
)

func GetArtefactClient() ArtefactServiceClient { 
    if client_ArtefactServiceClient_0 != nil {
        return client_ArtefactServiceClient_0
    }

    lock_ArtefactServiceClient_0.Lock() 
    if client_ArtefactServiceClient_0 != nil {
       lock_ArtefactServiceClient_0.Unlock()
       return client_ArtefactServiceClient_0
    }

    client_ArtefactServiceClient_0 = NewArtefactServiceClient(client.Connect("artefact.ArtefactService"))
    lock_ArtefactServiceClient_0.Unlock()
    return client_ArtefactServiceClient_0
}

func GetArtefactServiceClient() ArtefactServiceClient { 
    if client_ArtefactServiceClient_0 != nil {
        return client_ArtefactServiceClient_0
    }

    lock_ArtefactServiceClient_0.Lock() 
    if client_ArtefactServiceClient_0 != nil {
       lock_ArtefactServiceClient_0.Unlock()
       return client_ArtefactServiceClient_0
    }

    client_ArtefactServiceClient_0 = NewArtefactServiceClient(client.Connect("artefact.ArtefactService"))
    lock_ArtefactServiceClient_0.Unlock()
    return client_ArtefactServiceClient_0
}

