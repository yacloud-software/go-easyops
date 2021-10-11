// client create: EchoClient
/* geninfo:
   filename  : golang.lakeparime.com/apis/echoservice/echoservice.proto
   gopackage : golang.lakeparime.com/apis/echoservice
   importname: ai_0
   varname   : client_EchoClient_0
   clientname: EchoClient
   servername: EchoServer
   gscvname  : echo.Echo
   lockname  : lock_EchoClient_0
   activename: active_EchoClient_0
*/

package echo

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_EchoClient_0 sync.Mutex
  client_EchoClient_0 EchoClient
)

func GetEchoClient() EchoClient { 
    if client_EchoClient_0 != nil {
        return client_EchoClient_0
    }

    lock_EchoClient_0.Lock() 
    if client_EchoClient_0 != nil {
       lock_EchoClient_0.Unlock()
       return client_EchoClient_0
    }

    client_EchoClient_0 = NewEchoClient(client.Connect("echo.Echo"))
    lock_EchoClient_0.Unlock()
    return client_EchoClient_0
}

