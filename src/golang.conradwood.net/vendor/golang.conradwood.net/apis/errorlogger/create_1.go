// client create: ErrorLoggerClient
/* geninfo:
   filename  : golang.conradwood.net/apis/errorlogger/errorlogger.proto
   gopackage : golang.conradwood.net/apis/errorlogger
   importname: ai_0
   varname   : client_ErrorLoggerClient_0
   clientname: ErrorLoggerClient
   servername: ErrorLoggerServer
   gscvname  : errorlogger.ErrorLogger
   lockname  : lock_ErrorLoggerClient_0
   activename: active_ErrorLoggerClient_0
*/

package errorlogger

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ErrorLoggerClient_0 sync.Mutex
  client_ErrorLoggerClient_0 ErrorLoggerClient
)

func GetErrorLoggerClient() ErrorLoggerClient { 
    if client_ErrorLoggerClient_0 != nil {
        return client_ErrorLoggerClient_0
    }

    lock_ErrorLoggerClient_0.Lock() 
    if client_ErrorLoggerClient_0 != nil {
       lock_ErrorLoggerClient_0.Unlock()
       return client_ErrorLoggerClient_0
    }

    client_ErrorLoggerClient_0 = NewErrorLoggerClient(client.Connect("errorlogger.ErrorLogger"))
    lock_ErrorLoggerClient_0.Unlock()
    return client_ErrorLoggerClient_0
}

