// client create: FirewallMgrClient
/* geninfo:
   filename  : golang.conradwood.net/apis/firewallmgr/firewallmgr.proto
   gopackage : golang.conradwood.net/apis/firewallmgr
   importname: ai_0
   varname   : client_FirewallMgrClient_0
   clientname: FirewallMgrClient
   servername: FirewallMgrServer
   gscvname  : firewallmgr.FirewallMgr
   lockname  : lock_FirewallMgrClient_0
   activename: active_FirewallMgrClient_0
*/

package firewallmgr

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_FirewallMgrClient_0 sync.Mutex
  client_FirewallMgrClient_0 FirewallMgrClient
)

func GetFirewallMgrClient() FirewallMgrClient { 
    if client_FirewallMgrClient_0 != nil {
        return client_FirewallMgrClient_0
    }

    lock_FirewallMgrClient_0.Lock() 
    if client_FirewallMgrClient_0 != nil {
       lock_FirewallMgrClient_0.Unlock()
       return client_FirewallMgrClient_0
    }

    client_FirewallMgrClient_0 = NewFirewallMgrClient(client.Connect("firewallmgr.FirewallMgr"))
    lock_FirewallMgrClient_0.Unlock()
    return client_FirewallMgrClient_0
}

