// client create: SpeakToMeServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/speaktome/speaktome.proto
   gopackage : golang.conradwood.net/apis/speaktome
   importname: ai_0
   varname   : client_SpeakToMeServiceClient_0
   clientname: SpeakToMeServiceClient
   servername: SpeakToMeServiceServer
   gscvname  : speaktome.SpeakToMeService
   lockname  : lock_SpeakToMeServiceClient_0
   activename: active_SpeakToMeServiceClient_0
*/

package speaktome

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SpeakToMeServiceClient_0 sync.Mutex
  client_SpeakToMeServiceClient_0 SpeakToMeServiceClient
)

func GetSpeakToMeClient() SpeakToMeServiceClient { 
    if client_SpeakToMeServiceClient_0 != nil {
        return client_SpeakToMeServiceClient_0
    }

    lock_SpeakToMeServiceClient_0.Lock() 
    if client_SpeakToMeServiceClient_0 != nil {
       lock_SpeakToMeServiceClient_0.Unlock()
       return client_SpeakToMeServiceClient_0
    }

    client_SpeakToMeServiceClient_0 = NewSpeakToMeServiceClient(client.Connect("speaktome.SpeakToMeService"))
    lock_SpeakToMeServiceClient_0.Unlock()
    return client_SpeakToMeServiceClient_0
}

func GetSpeakToMeServiceClient() SpeakToMeServiceClient { 
    if client_SpeakToMeServiceClient_0 != nil {
        return client_SpeakToMeServiceClient_0
    }

    lock_SpeakToMeServiceClient_0.Lock() 
    if client_SpeakToMeServiceClient_0 != nil {
       lock_SpeakToMeServiceClient_0.Unlock()
       return client_SpeakToMeServiceClient_0
    }

    client_SpeakToMeServiceClient_0 = NewSpeakToMeServiceClient(client.Connect("speaktome.SpeakToMeService"))
    lock_SpeakToMeServiceClient_0.Unlock()
    return client_SpeakToMeServiceClient_0
}

