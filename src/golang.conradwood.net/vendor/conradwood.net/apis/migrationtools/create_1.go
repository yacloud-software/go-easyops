// client create: MigratorClient
/* geninfo:
   filename  : conradwood.net/apis/migrationtools/migrationtools.proto
   gopackage : conradwood.net/apis/migrationtools
   importname: ai_0
   varname   : client_MigratorClient_0
   clientname: MigratorClient
   servername: MigratorServer
   gscvname  : migrationtools.Migrator
   lockname  : lock_MigratorClient_0
   activename: active_MigratorClient_0
*/

package migrationtools

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_MigratorClient_0 sync.Mutex
  client_MigratorClient_0 MigratorClient
)

func GetMigratorClient() MigratorClient { 
    if client_MigratorClient_0 != nil {
        return client_MigratorClient_0
    }

    lock_MigratorClient_0.Lock() 
    if client_MigratorClient_0 != nil {
       lock_MigratorClient_0.Unlock()
       return client_MigratorClient_0
    }

    client_MigratorClient_0 = NewMigratorClient(client.Connect("migrationtools.Migrator"))
    lock_MigratorClient_0.Unlock()
    return client_MigratorClient_0
}

