// client create: SQLArchiverJobServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/sqlarchiver/sqlarchiver.proto
   gopackage : golang.conradwood.net/apis/sqlarchiver
   importname: ai_0
   varname   : client_SQLArchiverJobServiceClient_0
   clientname: SQLArchiverJobServiceClient
   servername: SQLArchiverJobServiceServer
   gscvname  : sqlarchiver.SQLArchiverJobService
   lockname  : lock_SQLArchiverJobServiceClient_0
   activename: active_SQLArchiverJobServiceClient_0
*/

package sqlarchiver

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SQLArchiverJobServiceClient_0 sync.Mutex
  client_SQLArchiverJobServiceClient_0 SQLArchiverJobServiceClient
)

func GetSQLArchiverJobClient() SQLArchiverJobServiceClient { 
    if client_SQLArchiverJobServiceClient_0 != nil {
        return client_SQLArchiverJobServiceClient_0
    }

    lock_SQLArchiverJobServiceClient_0.Lock() 
    if client_SQLArchiverJobServiceClient_0 != nil {
       lock_SQLArchiverJobServiceClient_0.Unlock()
       return client_SQLArchiverJobServiceClient_0
    }

    client_SQLArchiverJobServiceClient_0 = NewSQLArchiverJobServiceClient(client.Connect("sqlarchiver.SQLArchiverJobService"))
    lock_SQLArchiverJobServiceClient_0.Unlock()
    return client_SQLArchiverJobServiceClient_0
}

func GetSQLArchiverJobServiceClient() SQLArchiverJobServiceClient { 
    if client_SQLArchiverJobServiceClient_0 != nil {
        return client_SQLArchiverJobServiceClient_0
    }

    lock_SQLArchiverJobServiceClient_0.Lock() 
    if client_SQLArchiverJobServiceClient_0 != nil {
       lock_SQLArchiverJobServiceClient_0.Unlock()
       return client_SQLArchiverJobServiceClient_0
    }

    client_SQLArchiverJobServiceClient_0 = NewSQLArchiverJobServiceClient(client.Connect("sqlarchiver.SQLArchiverJobService"))
    lock_SQLArchiverJobServiceClient_0.Unlock()
    return client_SQLArchiverJobServiceClient_0
}

