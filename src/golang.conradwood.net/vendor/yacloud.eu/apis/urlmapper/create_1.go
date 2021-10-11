// client create: URLMapperClient
/* geninfo:
   filename  : yacloud.eu/apis/urlmapper/urlmapper.proto
   gopackage : yacloud.eu/apis/urlmapper
   importname: ai_0
   varname   : client_URLMapperClient_0
   clientname: URLMapperClient
   servername: URLMapperServer
   gscvname  : urlmapper.URLMapper
   lockname  : lock_URLMapperClient_0
   activename: active_URLMapperClient_0
*/

package urlmapper

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_URLMapperClient_0 sync.Mutex
  client_URLMapperClient_0 URLMapperClient
)

func GetURLMapperClient() URLMapperClient { 
    if client_URLMapperClient_0 != nil {
        return client_URLMapperClient_0
    }

    lock_URLMapperClient_0.Lock() 
    if client_URLMapperClient_0 != nil {
       lock_URLMapperClient_0.Unlock()
       return client_URLMapperClient_0
    }

    client_URLMapperClient_0 = NewURLMapperClient(client.Connect("urlmapper.URLMapper"))
    lock_URLMapperClient_0.Unlock()
    return client_URLMapperClient_0
}

