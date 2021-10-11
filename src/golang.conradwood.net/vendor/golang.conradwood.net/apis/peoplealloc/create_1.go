// client create: PeopleAllocateClient
/* geninfo:
   filename  : golang.conradwood.net/apis/peoplealloc/peoplealloc.proto
   gopackage : golang.conradwood.net/apis/peoplealloc
   importname: ai_0
   varname   : client_PeopleAllocateClient_0
   clientname: PeopleAllocateClient
   servername: PeopleAllocateServer
   gscvname  : peoplealloc.PeopleAllocate
   lockname  : lock_PeopleAllocateClient_0
   activename: active_PeopleAllocateClient_0
*/

package peoplealloc

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_PeopleAllocateClient_0 sync.Mutex
  client_PeopleAllocateClient_0 PeopleAllocateClient
)

func GetPeopleAllocateClient() PeopleAllocateClient { 
    if client_PeopleAllocateClient_0 != nil {
        return client_PeopleAllocateClient_0
    }

    lock_PeopleAllocateClient_0.Lock() 
    if client_PeopleAllocateClient_0 != nil {
       lock_PeopleAllocateClient_0.Unlock()
       return client_PeopleAllocateClient_0
    }

    client_PeopleAllocateClient_0 = NewPeopleAllocateClient(client.Connect("peoplealloc.PeopleAllocate"))
    lock_PeopleAllocateClient_0.Unlock()
    return client_PeopleAllocateClient_0
}

