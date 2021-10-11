// client create: GDriveClient
/* geninfo:
   filename  : golang.conradwood.net/apis/gdrive/gdrive.proto
   gopackage : golang.conradwood.net/apis/gdrive
   importname: ai_0
   varname   : client_GDriveClient_0
   clientname: GDriveClient
   servername: GDriveServer
   gscvname  : gdrive.GDrive
   lockname  : lock_GDriveClient_0
   activename: active_GDriveClient_0
*/

package gdrive

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_GDriveClient_0 sync.Mutex
  client_GDriveClient_0 GDriveClient
)

func GetGDriveClient() GDriveClient { 
    if client_GDriveClient_0 != nil {
        return client_GDriveClient_0
    }

    lock_GDriveClient_0.Lock() 
    if client_GDriveClient_0 != nil {
       lock_GDriveClient_0.Unlock()
       return client_GDriveClient_0
    }

    client_GDriveClient_0 = NewGDriveClient(client.Connect("gdrive.GDrive"))
    lock_GDriveClient_0.Unlock()
    return client_GDriveClient_0
}

