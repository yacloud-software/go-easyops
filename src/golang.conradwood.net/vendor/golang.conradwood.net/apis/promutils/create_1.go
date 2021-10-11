// client create: PromUtilsClient
/* geninfo:
   filename  : golang.conradwood.net/apis/promutils/promutils.proto
   gopackage : golang.conradwood.net/apis/promutils
   importname: ai_0
   varname   : client_PromUtilsClient_0
   clientname: PromUtilsClient
   servername: PromUtilsServer
   gscvname  : promutils.PromUtils
   lockname  : lock_PromUtilsClient_0
   activename: active_PromUtilsClient_0
*/

package promutils

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_PromUtilsClient_0 sync.Mutex
  client_PromUtilsClient_0 PromUtilsClient
)

func GetPromUtilsClient() PromUtilsClient { 
    if client_PromUtilsClient_0 != nil {
        return client_PromUtilsClient_0
    }

    lock_PromUtilsClient_0.Lock() 
    if client_PromUtilsClient_0 != nil {
       lock_PromUtilsClient_0.Unlock()
       return client_PromUtilsClient_0
    }

    client_PromUtilsClient_0 = NewPromUtilsClient(client.Connect("promutils.PromUtils"))
    lock_PromUtilsClient_0.Unlock()
    return client_PromUtilsClient_0
}

