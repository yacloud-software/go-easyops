// client create: MKDBClient
/* geninfo:
   filename  : golang.conradwood.net/apis/mkdb/mkdb.proto
   gopackage : golang.conradwood.net/apis/mkdb
   importname: ai_0
   varname   : client_MKDBClient_0
   clientname: MKDBClient
   servername: MKDBServer
   gscvname  : mkdb.MKDB
   lockname  : lock_MKDBClient_0
   activename: active_MKDBClient_0
*/

package mkdb

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_MKDBClient_0 sync.Mutex
  client_MKDBClient_0 MKDBClient
)

func GetMKDBClient() MKDBClient { 
    if client_MKDBClient_0 != nil {
        return client_MKDBClient_0
    }

    lock_MKDBClient_0.Lock() 
    if client_MKDBClient_0 != nil {
       lock_MKDBClient_0.Unlock()
       return client_MKDBClient_0
    }

    client_MKDBClient_0 = NewMKDBClient(client.Connect("mkdb.MKDB"))
    lock_MKDBClient_0.Unlock()
    return client_MKDBClient_0
}

