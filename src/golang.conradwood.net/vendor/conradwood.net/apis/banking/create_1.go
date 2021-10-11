// client create: BankingClient
/* geninfo:
   filename  : conradwood.net/apis/banking/banking.proto
   gopackage : conradwood.net/apis/banking
   importname: ai_0
   varname   : client_BankingClient_0
   clientname: BankingClient
   servername: BankingServer
   gscvname  : banking.Banking
   lockname  : lock_BankingClient_0
   activename: active_BankingClient_0
*/

package banking

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_BankingClient_0 sync.Mutex
  client_BankingClient_0 BankingClient
)

func GetBankingClient() BankingClient { 
    if client_BankingClient_0 != nil {
        return client_BankingClient_0
    }

    lock_BankingClient_0.Lock() 
    if client_BankingClient_0 != nil {
       lock_BankingClient_0.Unlock()
       return client_BankingClient_0
    }

    client_BankingClient_0 = NewBankingClient(client.Connect("banking.Banking"))
    lock_BankingClient_0.Unlock()
    return client_BankingClient_0
}

