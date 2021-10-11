// client create: JSONApiMultiplexerClient
/* geninfo:
   filename  : golang.conradwood.net/apis/jsonapimultiplexer/jsonapimultiplexer.proto
   gopackage : golang.conradwood.net/apis/jsonapimultiplexer
   importname: ai_0
   varname   : client_JSONApiMultiplexerClient_0
   clientname: JSONApiMultiplexerClient
   servername: JSONApiMultiplexerServer
   gscvname  : jsonapimultiplexer.JSONApiMultiplexer
   lockname  : lock_JSONApiMultiplexerClient_0
   activename: active_JSONApiMultiplexerClient_0
*/

package jsonapimultiplexer

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_JSONApiMultiplexerClient_0 sync.Mutex
  client_JSONApiMultiplexerClient_0 JSONApiMultiplexerClient
)

func GetJSONApiMultiplexerClient() JSONApiMultiplexerClient { 
    if client_JSONApiMultiplexerClient_0 != nil {
        return client_JSONApiMultiplexerClient_0
    }

    lock_JSONApiMultiplexerClient_0.Lock() 
    if client_JSONApiMultiplexerClient_0 != nil {
       lock_JSONApiMultiplexerClient_0.Unlock()
       return client_JSONApiMultiplexerClient_0
    }

    client_JSONApiMultiplexerClient_0 = NewJSONApiMultiplexerClient(client.Connect("jsonapimultiplexer.JSONApiMultiplexer"))
    lock_JSONApiMultiplexerClient_0.Unlock()
    return client_JSONApiMultiplexerClient_0
}

