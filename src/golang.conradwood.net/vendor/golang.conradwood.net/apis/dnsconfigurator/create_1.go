// client create: DNSConfiguratorServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/dnsconfigurator/dnsconfigurator.proto
   gopackage : golang.conradwood.net/apis/dnsconfigurator
   importname: ai_0
   varname   : client_DNSConfiguratorServiceClient_0
   clientname: DNSConfiguratorServiceClient
   servername: DNSConfiguratorServiceServer
   gscvname  : dnsconfigurator.DNSConfiguratorService
   lockname  : lock_DNSConfiguratorServiceClient_0
   activename: active_DNSConfiguratorServiceClient_0
*/

package dnsconfigurator

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_DNSConfiguratorServiceClient_0 sync.Mutex
  client_DNSConfiguratorServiceClient_0 DNSConfiguratorServiceClient
)

func GetDNSConfiguratorClient() DNSConfiguratorServiceClient { 
    if client_DNSConfiguratorServiceClient_0 != nil {
        return client_DNSConfiguratorServiceClient_0
    }

    lock_DNSConfiguratorServiceClient_0.Lock() 
    if client_DNSConfiguratorServiceClient_0 != nil {
       lock_DNSConfiguratorServiceClient_0.Unlock()
       return client_DNSConfiguratorServiceClient_0
    }

    client_DNSConfiguratorServiceClient_0 = NewDNSConfiguratorServiceClient(client.Connect("dnsconfigurator.DNSConfiguratorService"))
    lock_DNSConfiguratorServiceClient_0.Unlock()
    return client_DNSConfiguratorServiceClient_0
}

func GetDNSConfiguratorServiceClient() DNSConfiguratorServiceClient { 
    if client_DNSConfiguratorServiceClient_0 != nil {
        return client_DNSConfiguratorServiceClient_0
    }

    lock_DNSConfiguratorServiceClient_0.Lock() 
    if client_DNSConfiguratorServiceClient_0 != nil {
       lock_DNSConfiguratorServiceClient_0.Unlock()
       return client_DNSConfiguratorServiceClient_0
    }

    client_DNSConfiguratorServiceClient_0 = NewDNSConfiguratorServiceClient(client.Connect("dnsconfigurator.DNSConfiguratorService"))
    lock_DNSConfiguratorServiceClient_0.Unlock()
    return client_DNSConfiguratorServiceClient_0
}

