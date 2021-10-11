// client create: CertManagerClient
/* geninfo:
   filename  : golang.conradwood.net/apis/certmanager/certmanager.proto
   gopackage : golang.conradwood.net/apis/certmanager
   importname: ai_0
   varname   : client_CertManagerClient_0
   clientname: CertManagerClient
   servername: CertManagerServer
   gscvname  : certmanager.CertManager
   lockname  : lock_CertManagerClient_0
   activename: active_CertManagerClient_0
*/

package certmanager

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_CertManagerClient_0 sync.Mutex
  client_CertManagerClient_0 CertManagerClient
)

func GetCertManagerClient() CertManagerClient { 
    if client_CertManagerClient_0 != nil {
        return client_CertManagerClient_0
    }

    lock_CertManagerClient_0.Lock() 
    if client_CertManagerClient_0 != nil {
       lock_CertManagerClient_0.Unlock()
       return client_CertManagerClient_0
    }

    client_CertManagerClient_0 = NewCertManagerClient(client.Connect("certmanager.CertManager"))
    lock_CertManagerClient_0.Unlock()
    return client_CertManagerClient_0
}

