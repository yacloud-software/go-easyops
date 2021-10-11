// client create: ThemesClient
/* geninfo:
   filename  : golang.conradwood.net/apis/themes/themes.proto
   gopackage : golang.conradwood.net/apis/themes
   importname: ai_0
   varname   : client_ThemesClient_0
   clientname: ThemesClient
   servername: ThemesServer
   gscvname  : themes.Themes
   lockname  : lock_ThemesClient_0
   activename: active_ThemesClient_0
*/

package themes

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ThemesClient_0 sync.Mutex
  client_ThemesClient_0 ThemesClient
)

func GetThemesClient() ThemesClient { 
    if client_ThemesClient_0 != nil {
        return client_ThemesClient_0
    }

    lock_ThemesClient_0.Lock() 
    if client_ThemesClient_0 != nil {
       lock_ThemesClient_0.Unlock()
       return client_ThemesClient_0
    }

    client_ThemesClient_0 = NewThemesClient(client.Connect("themes.Themes"))
    lock_ThemesClient_0.Unlock()
    return client_ThemesClient_0
}

