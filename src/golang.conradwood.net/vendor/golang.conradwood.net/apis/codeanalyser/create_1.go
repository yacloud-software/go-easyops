// client create: CodeAnalyserServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/codeanalyser/codeanalyser.proto
   gopackage : golang.conradwood.net/apis/codeanalyser
   importname: ai_0
   varname   : client_CodeAnalyserServiceClient_0
   clientname: CodeAnalyserServiceClient
   servername: CodeAnalyserServiceServer
   gscvname  : codeanalyser.CodeAnalyserService
   lockname  : lock_CodeAnalyserServiceClient_0
   activename: active_CodeAnalyserServiceClient_0
*/

package codeanalyser

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_CodeAnalyserServiceClient_0 sync.Mutex
  client_CodeAnalyserServiceClient_0 CodeAnalyserServiceClient
)

func GetCodeAnalyserClient() CodeAnalyserServiceClient { 
    if client_CodeAnalyserServiceClient_0 != nil {
        return client_CodeAnalyserServiceClient_0
    }

    lock_CodeAnalyserServiceClient_0.Lock() 
    if client_CodeAnalyserServiceClient_0 != nil {
       lock_CodeAnalyserServiceClient_0.Unlock()
       return client_CodeAnalyserServiceClient_0
    }

    client_CodeAnalyserServiceClient_0 = NewCodeAnalyserServiceClient(client.Connect("codeanalyser.CodeAnalyserService"))
    lock_CodeAnalyserServiceClient_0.Unlock()
    return client_CodeAnalyserServiceClient_0
}

func GetCodeAnalyserServiceClient() CodeAnalyserServiceClient { 
    if client_CodeAnalyserServiceClient_0 != nil {
        return client_CodeAnalyserServiceClient_0
    }

    lock_CodeAnalyserServiceClient_0.Lock() 
    if client_CodeAnalyserServiceClient_0 != nil {
       lock_CodeAnalyserServiceClient_0.Unlock()
       return client_CodeAnalyserServiceClient_0
    }

    client_CodeAnalyserServiceClient_0 = NewCodeAnalyserServiceClient(client.Connect("codeanalyser.CodeAnalyserService"))
    lock_CodeAnalyserServiceClient_0.Unlock()
    return client_CodeAnalyserServiceClient_0
}

