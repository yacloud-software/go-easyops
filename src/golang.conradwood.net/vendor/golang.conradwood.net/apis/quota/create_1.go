// client create: QuotaServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/quota/quota.proto
   gopackage : golang.conradwood.net/apis/quota
   importname: ai_0
   varname   : client_QuotaServiceClient_0
   clientname: QuotaServiceClient
   servername: QuotaServiceServer
   gscvname  : quota.QuotaService
   lockname  : lock_QuotaServiceClient_0
   activename: active_QuotaServiceClient_0
*/

package quota

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_QuotaServiceClient_0 sync.Mutex
  client_QuotaServiceClient_0 QuotaServiceClient
)

func GetQuotaClient() QuotaServiceClient { 
    if client_QuotaServiceClient_0 != nil {
        return client_QuotaServiceClient_0
    }

    lock_QuotaServiceClient_0.Lock() 
    if client_QuotaServiceClient_0 != nil {
       lock_QuotaServiceClient_0.Unlock()
       return client_QuotaServiceClient_0
    }

    client_QuotaServiceClient_0 = NewQuotaServiceClient(client.Connect("quota.QuotaService"))
    lock_QuotaServiceClient_0.Unlock()
    return client_QuotaServiceClient_0
}

func GetQuotaServiceClient() QuotaServiceClient { 
    if client_QuotaServiceClient_0 != nil {
        return client_QuotaServiceClient_0
    }

    lock_QuotaServiceClient_0.Lock() 
    if client_QuotaServiceClient_0 != nil {
       lock_QuotaServiceClient_0.Unlock()
       return client_QuotaServiceClient_0
    }

    client_QuotaServiceClient_0 = NewQuotaServiceClient(client.Connect("quota.QuotaService"))
    lock_QuotaServiceClient_0.Unlock()
    return client_QuotaServiceClient_0
}

