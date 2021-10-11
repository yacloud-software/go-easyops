// client create: DocumentsClient
/* geninfo:
   filename  : conradwood.net/apis/documents/documents.proto
   gopackage : conradwood.net/apis/documents
   importname: ai_0
   varname   : client_DocumentsClient_0
   clientname: DocumentsClient
   servername: DocumentsServer
   gscvname  : documents.Documents
   lockname  : lock_DocumentsClient_0
   activename: active_DocumentsClient_0
*/

package documents

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_DocumentsClient_0 sync.Mutex
  client_DocumentsClient_0 DocumentsClient
)

func GetDocumentsClient() DocumentsClient { 
    if client_DocumentsClient_0 != nil {
        return client_DocumentsClient_0
    }

    lock_DocumentsClient_0.Lock() 
    if client_DocumentsClient_0 != nil {
       lock_DocumentsClient_0.Unlock()
       return client_DocumentsClient_0
    }

    client_DocumentsClient_0 = NewDocumentsClient(client.Connect("documents.Documents"))
    lock_DocumentsClient_0.Unlock()
    return client_DocumentsClient_0
}

