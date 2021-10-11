// client create: IPExporterServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/ipexporter/ipexporter.proto
   gopackage : golang.conradwood.net/apis/ipexporter
   importname: ai_0
   varname   : client_IPExporterServiceClient_0
   clientname: IPExporterServiceClient
   servername: IPExporterServiceServer
   gscvname  : ipexporter.IPExporterService
   lockname  : lock_IPExporterServiceClient_0
   activename: active_IPExporterServiceClient_0
*/

package ipexporter

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_IPExporterServiceClient_0 sync.Mutex
  client_IPExporterServiceClient_0 IPExporterServiceClient
)

func GetIPExporterClient() IPExporterServiceClient { 
    if client_IPExporterServiceClient_0 != nil {
        return client_IPExporterServiceClient_0
    }

    lock_IPExporterServiceClient_0.Lock() 
    if client_IPExporterServiceClient_0 != nil {
       lock_IPExporterServiceClient_0.Unlock()
       return client_IPExporterServiceClient_0
    }

    client_IPExporterServiceClient_0 = NewIPExporterServiceClient(client.Connect("ipexporter.IPExporterService"))
    lock_IPExporterServiceClient_0.Unlock()
    return client_IPExporterServiceClient_0
}

func GetIPExporterServiceClient() IPExporterServiceClient { 
    if client_IPExporterServiceClient_0 != nil {
        return client_IPExporterServiceClient_0
    }

    lock_IPExporterServiceClient_0.Lock() 
    if client_IPExporterServiceClient_0 != nil {
       lock_IPExporterServiceClient_0.Unlock()
       return client_IPExporterServiceClient_0
    }

    client_IPExporterServiceClient_0 = NewIPExporterServiceClient(client.Connect("ipexporter.IPExporterService"))
    lock_IPExporterServiceClient_0.Unlock()
    return client_IPExporterServiceClient_0
}

