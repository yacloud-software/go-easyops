// client create: DocumentProcessorClient
/* geninfo:
   filename  : conradwood.net/apis/documents/documents.proto
   gopackage : conradwood.net/apis/documents
   importname: ai_1
   varname   : client_DocumentProcessorClient_1
   clientname: DocumentProcessorClient
   servername: DocumentProcessorServer
   gscvname  : documents.DocumentProcessor
   lockname  : lock_DocumentProcessorClient_1
   activename: active_DocumentProcessorClient_1
*/

package documents

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_DocumentProcessorClient_1 sync.Mutex
  client_DocumentProcessorClient_1 DocumentProcessorClient
)

func GetDocumentProcessorClient() DocumentProcessorClient { 
    if client_DocumentProcessorClient_1 != nil {
        return client_DocumentProcessorClient_1
    }

    lock_DocumentProcessorClient_1.Lock() 
    if client_DocumentProcessorClient_1 != nil {
       lock_DocumentProcessorClient_1.Unlock()
       return client_DocumentProcessorClient_1
    }

    client_DocumentProcessorClient_1 = NewDocumentProcessorClient(client.Connect("documents.DocumentProcessor"))
    lock_DocumentProcessorClient_1.Unlock()
    return client_DocumentProcessorClient_1
}

