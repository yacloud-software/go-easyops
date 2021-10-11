// client create: DHCPConfiguratorClient
/* geninfo:
   filename  : golang.conradwood.net/apis/dhcpconfigurator/dhcpconfigurator.proto
   gopackage : golang.conradwood.net/apis/dhcpconfigurator
   importname: ai_0
   varname   : client_DHCPConfiguratorClient_0
   clientname: DHCPConfiguratorClient
   servername: DHCPConfiguratorServer
   gscvname  : dhcpconfigurator.DHCPConfigurator
   lockname  : lock_DHCPConfiguratorClient_0
   activename: active_DHCPConfiguratorClient_0
*/

package dhcpconfigurator

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_DHCPConfiguratorClient_0 sync.Mutex
  client_DHCPConfiguratorClient_0 DHCPConfiguratorClient
)

func GetDHCPConfiguratorClient() DHCPConfiguratorClient { 
    if client_DHCPConfiguratorClient_0 != nil {
        return client_DHCPConfiguratorClient_0
    }

    lock_DHCPConfiguratorClient_0.Lock() 
    if client_DHCPConfiguratorClient_0 != nil {
       lock_DHCPConfiguratorClient_0.Unlock()
       return client_DHCPConfiguratorClient_0
    }

    client_DHCPConfiguratorClient_0 = NewDHCPConfiguratorClient(client.Connect("dhcpconfigurator.DHCPConfigurator"))
    lock_DHCPConfiguratorClient_0.Unlock()
    return client_DHCPConfiguratorClient_0
}

