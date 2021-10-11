// client create: Covid19ServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/covid19/covid19.proto
   gopackage : golang.conradwood.net/apis/covid19
   importname: ai_0
   varname   : client_Covid19ServiceClient_0
   clientname: Covid19ServiceClient
   servername: Covid19ServiceServer
   gscvname  : covid19.Covid19Service
   lockname  : lock_Covid19ServiceClient_0
   activename: active_Covid19ServiceClient_0
*/

package covid19

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_Covid19ServiceClient_0 sync.Mutex
  client_Covid19ServiceClient_0 Covid19ServiceClient
)

func GetCovid19Client() Covid19ServiceClient { 
    if client_Covid19ServiceClient_0 != nil {
        return client_Covid19ServiceClient_0
    }

    lock_Covid19ServiceClient_0.Lock() 
    if client_Covid19ServiceClient_0 != nil {
       lock_Covid19ServiceClient_0.Unlock()
       return client_Covid19ServiceClient_0
    }

    client_Covid19ServiceClient_0 = NewCovid19ServiceClient(client.Connect("covid19.Covid19Service"))
    lock_Covid19ServiceClient_0.Unlock()
    return client_Covid19ServiceClient_0
}

func GetCovid19ServiceClient() Covid19ServiceClient { 
    if client_Covid19ServiceClient_0 != nil {
        return client_Covid19ServiceClient_0
    }

    lock_Covid19ServiceClient_0.Lock() 
    if client_Covid19ServiceClient_0 != nil {
       lock_Covid19ServiceClient_0.Unlock()
       return client_Covid19ServiceClient_0
    }

    client_Covid19ServiceClient_0 = NewCovid19ServiceClient(client.Connect("covid19.Covid19Service"))
    lock_Covid19ServiceClient_0.Unlock()
    return client_Covid19ServiceClient_0
}

