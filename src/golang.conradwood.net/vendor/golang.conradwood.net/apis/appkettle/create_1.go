// client create: AppKettleServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/appkettle/appkettle.proto
   gopackage : golang.conradwood.net/apis/appkettle
   importname: ai_0
   varname   : client_AppKettleServiceClient_0
   clientname: AppKettleServiceClient
   servername: AppKettleServiceServer
   gscvname  : appkettle.AppKettleService
   lockname  : lock_AppKettleServiceClient_0
   activename: active_AppKettleServiceClient_0
*/

package appkettle

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_AppKettleServiceClient_0 sync.Mutex
  client_AppKettleServiceClient_0 AppKettleServiceClient
)

func GetAppKettleClient() AppKettleServiceClient { 
    if client_AppKettleServiceClient_0 != nil {
        return client_AppKettleServiceClient_0
    }

    lock_AppKettleServiceClient_0.Lock() 
    if client_AppKettleServiceClient_0 != nil {
       lock_AppKettleServiceClient_0.Unlock()
       return client_AppKettleServiceClient_0
    }

    client_AppKettleServiceClient_0 = NewAppKettleServiceClient(client.Connect("appkettle.AppKettleService"))
    lock_AppKettleServiceClient_0.Unlock()
    return client_AppKettleServiceClient_0
}

func GetAppKettleServiceClient() AppKettleServiceClient { 
    if client_AppKettleServiceClient_0 != nil {
        return client_AppKettleServiceClient_0
    }

    lock_AppKettleServiceClient_0.Lock() 
    if client_AppKettleServiceClient_0 != nil {
       lock_AppKettleServiceClient_0.Unlock()
       return client_AppKettleServiceClient_0
    }

    client_AppKettleServiceClient_0 = NewAppKettleServiceClient(client.Connect("appkettle.AppKettleService"))
    lock_AppKettleServiceClient_0.Unlock()
    return client_AppKettleServiceClient_0
}

