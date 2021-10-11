// client create: HelloWorldClient
/* geninfo:
   filename  : golang.conradwood.net/apis/helloworld/helloworld.proto
   gopackage : golang.conradwood.net/apis/helloworld
   importname: ai_0
   varname   : client_HelloWorldClient_0
   clientname: HelloWorldClient
   servername: HelloWorldServer
   gscvname  : helloworld.HelloWorld
   lockname  : lock_HelloWorldClient_0
   activename: active_HelloWorldClient_0
*/

package helloworld

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_HelloWorldClient_0 sync.Mutex
  client_HelloWorldClient_0 HelloWorldClient
)

func GetHelloWorldClient() HelloWorldClient { 
    if client_HelloWorldClient_0 != nil {
        return client_HelloWorldClient_0
    }

    lock_HelloWorldClient_0.Lock() 
    if client_HelloWorldClient_0 != nil {
       lock_HelloWorldClient_0.Unlock()
       return client_HelloWorldClient_0
    }

    client_HelloWorldClient_0 = NewHelloWorldClient(client.Connect("helloworld.HelloWorld"))
    lock_HelloWorldClient_0.Unlock()
    return client_HelloWorldClient_0
}

