// client create: LogServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/logservice/logservice.proto
   gopackage : golang.conradwood.net/apis/logservice
   importname: ai_0
   varname   : client_LogServiceClient_0
   clientname: LogServiceClient
   servername: LogServiceServer
   gscvname  : logservice.LogService
   lockname  : lock_LogServiceClient_0
   activename: active_LogServiceClient_0
*/

package logservice

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_LogServiceClient_0 sync.Mutex
  client_LogServiceClient_0 LogServiceClient
)

func GetLogClient() LogServiceClient { 
    if client_LogServiceClient_0 != nil {
        return client_LogServiceClient_0
    }

    lock_LogServiceClient_0.Lock() 
    if client_LogServiceClient_0 != nil {
       lock_LogServiceClient_0.Unlock()
       return client_LogServiceClient_0
    }

    client_LogServiceClient_0 = NewLogServiceClient(client.Connect("logservice.LogService"))
    lock_LogServiceClient_0.Unlock()
    return client_LogServiceClient_0
}

func GetLogServiceClient() LogServiceClient { 
    if client_LogServiceClient_0 != nil {
        return client_LogServiceClient_0
    }

    lock_LogServiceClient_0.Lock() 
    if client_LogServiceClient_0 != nil {
       lock_LogServiceClient_0.Unlock()
       return client_LogServiceClient_0
    }

    client_LogServiceClient_0 = NewLogServiceClient(client.Connect("logservice.LogService"))
    lock_LogServiceClient_0.Unlock()
    return client_LogServiceClient_0
}

