// client create: BankHolidaysClient
/* geninfo:
   filename  : golang.conradwood.net/apis/bankholidays/bankholidays.proto
   gopackage : golang.conradwood.net/apis/bankholidays
   importname: ai_0
   varname   : client_BankHolidaysClient_0
   clientname: BankHolidaysClient
   servername: BankHolidaysServer
   gscvname  : bankholidays.BankHolidays
   lockname  : lock_BankHolidaysClient_0
   activename: active_BankHolidaysClient_0
*/

package bankholidays

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_BankHolidaysClient_0 sync.Mutex
  client_BankHolidaysClient_0 BankHolidaysClient
)

func GetBankHolidaysClient() BankHolidaysClient { 
    if client_BankHolidaysClient_0 != nil {
        return client_BankHolidaysClient_0
    }

    lock_BankHolidaysClient_0.Lock() 
    if client_BankHolidaysClient_0 != nil {
       lock_BankHolidaysClient_0.Unlock()
       return client_BankHolidaysClient_0
    }

    client_BankHolidaysClient_0 = NewBankHolidaysClient(client.Connect("bankholidays.BankHolidays"))
    lock_BankHolidaysClient_0.Unlock()
    return client_BankHolidaysClient_0
}

