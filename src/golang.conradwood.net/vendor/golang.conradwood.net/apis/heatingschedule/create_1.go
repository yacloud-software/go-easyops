// client create: HeatingScheduleServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/heatingschedule/heatingschedule.proto
   gopackage : golang.conradwood.net/apis/heatingschedule
   importname: ai_0
   varname   : client_HeatingScheduleServiceClient_0
   clientname: HeatingScheduleServiceClient
   servername: HeatingScheduleServiceServer
   gscvname  : heatingschedule.HeatingScheduleService
   lockname  : lock_HeatingScheduleServiceClient_0
   activename: active_HeatingScheduleServiceClient_0
*/

package heatingschedule

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_HeatingScheduleServiceClient_0 sync.Mutex
  client_HeatingScheduleServiceClient_0 HeatingScheduleServiceClient
)

func GetHeatingScheduleClient() HeatingScheduleServiceClient { 
    if client_HeatingScheduleServiceClient_0 != nil {
        return client_HeatingScheduleServiceClient_0
    }

    lock_HeatingScheduleServiceClient_0.Lock() 
    if client_HeatingScheduleServiceClient_0 != nil {
       lock_HeatingScheduleServiceClient_0.Unlock()
       return client_HeatingScheduleServiceClient_0
    }

    client_HeatingScheduleServiceClient_0 = NewHeatingScheduleServiceClient(client.Connect("heatingschedule.HeatingScheduleService"))
    lock_HeatingScheduleServiceClient_0.Unlock()
    return client_HeatingScheduleServiceClient_0
}

func GetHeatingScheduleServiceClient() HeatingScheduleServiceClient { 
    if client_HeatingScheduleServiceClient_0 != nil {
        return client_HeatingScheduleServiceClient_0
    }

    lock_HeatingScheduleServiceClient_0.Lock() 
    if client_HeatingScheduleServiceClient_0 != nil {
       lock_HeatingScheduleServiceClient_0.Unlock()
       return client_HeatingScheduleServiceClient_0
    }

    client_HeatingScheduleServiceClient_0 = NewHeatingScheduleServiceClient(client.Connect("heatingschedule.HeatingScheduleService"))
    lock_HeatingScheduleServiceClient_0.Unlock()
    return client_HeatingScheduleServiceClient_0
}

