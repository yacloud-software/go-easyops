// client create: UIObjectsClient
/* geninfo:
   filename  : golang.yacloud.eu/apis/uiobjects/uiobjects.proto
   gopackage : golang.yacloud.eu/apis/uiobjects
   importname: ai_0
   varname   : client_UIObjectsClient_0
   clientname: UIObjectsClient
   servername: UIObjectsServer
   gscvname  : uiobjects.UIObjects
   lockname  : lock_UIObjectsClient_0
   activename: active_UIObjectsClient_0
*/

package uiobjects

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_UIObjectsClient_0 sync.Mutex
  client_UIObjectsClient_0 UIObjectsClient
)

func GetUIObjectsClient() UIObjectsClient { 
    if client_UIObjectsClient_0 != nil {
        return client_UIObjectsClient_0
    }

    lock_UIObjectsClient_0.Lock() 
    if client_UIObjectsClient_0 != nil {
       lock_UIObjectsClient_0.Unlock()
       return client_UIObjectsClient_0
    }

    client_UIObjectsClient_0 = NewUIObjectsClient(client.Connect("uiobjects.UIObjects"))
    lock_UIObjectsClient_0.Unlock()
    return client_UIObjectsClient_0
}

