// client create: DNSConfigServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/dnsconfig/dnsconfig.proto
   gopackage : golang.conradwood.net/apis/dnsconfig
   importname: ai_0
   varname   : client_DNSConfigServiceClient_0
   clientname: DNSConfigServiceClient
   servername: DNSConfigServiceServer
   gscvname  : dnsconfig.DNSConfigService
   lockname  : lock_DNSConfigServiceClient_0
   activename: active_DNSConfigServiceClient_0
*/

package dnsconfig

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_DNSConfigServiceClient_0 sync.Mutex
  client_DNSConfigServiceClient_0 DNSConfigServiceClient
)

func GetDNSConfigClient() DNSConfigServiceClient { 
    if client_DNSConfigServiceClient_0 != nil {
        return client_DNSConfigServiceClient_0
    }

    lock_DNSConfigServiceClient_0.Lock() 
    if client_DNSConfigServiceClient_0 != nil {
       lock_DNSConfigServiceClient_0.Unlock()
       return client_DNSConfigServiceClient_0
    }

    client_DNSConfigServiceClient_0 = NewDNSConfigServiceClient(client.Connect("dnsconfig.DNSConfigService"))
    lock_DNSConfigServiceClient_0.Unlock()
    return client_DNSConfigServiceClient_0
}

func GetDNSConfigServiceClient() DNSConfigServiceClient { 
    if client_DNSConfigServiceClient_0 != nil {
        return client_DNSConfigServiceClient_0
    }

    lock_DNSConfigServiceClient_0.Lock() 
    if client_DNSConfigServiceClient_0 != nil {
       lock_DNSConfigServiceClient_0.Unlock()
       return client_DNSConfigServiceClient_0
    }

    client_DNSConfigServiceClient_0 = NewDNSConfigServiceClient(client.Connect("dnsconfig.DNSConfigService"))
    lock_DNSConfigServiceClient_0.Unlock()
    return client_DNSConfigServiceClient_0
}

