// client create: PhoneBookServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/phonebook/phonebook.proto
   gopackage : golang.conradwood.net/apis/phonebook
   importname: ai_0
   varname   : client_PhoneBookServiceClient_0
   clientname: PhoneBookServiceClient
   servername: PhoneBookServiceServer
   gscvname  : phonebook.PhoneBookService
   lockname  : lock_PhoneBookServiceClient_0
   activename: active_PhoneBookServiceClient_0
*/

package phonebook

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_PhoneBookServiceClient_0 sync.Mutex
  client_PhoneBookServiceClient_0 PhoneBookServiceClient
)

func GetPhoneBookClient() PhoneBookServiceClient { 
    if client_PhoneBookServiceClient_0 != nil {
        return client_PhoneBookServiceClient_0
    }

    lock_PhoneBookServiceClient_0.Lock() 
    if client_PhoneBookServiceClient_0 != nil {
       lock_PhoneBookServiceClient_0.Unlock()
       return client_PhoneBookServiceClient_0
    }

    client_PhoneBookServiceClient_0 = NewPhoneBookServiceClient(client.Connect("phonebook.PhoneBookService"))
    lock_PhoneBookServiceClient_0.Unlock()
    return client_PhoneBookServiceClient_0
}

func GetPhoneBookServiceClient() PhoneBookServiceClient { 
    if client_PhoneBookServiceClient_0 != nil {
        return client_PhoneBookServiceClient_0
    }

    lock_PhoneBookServiceClient_0.Lock() 
    if client_PhoneBookServiceClient_0 != nil {
       lock_PhoneBookServiceClient_0.Unlock()
       return client_PhoneBookServiceClient_0
    }

    client_PhoneBookServiceClient_0 = NewPhoneBookServiceClient(client.Connect("phonebook.PhoneBookService"))
    lock_PhoneBookServiceClient_0.Unlock()
    return client_PhoneBookServiceClient_0
}

