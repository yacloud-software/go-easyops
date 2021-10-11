// client create: SensorServerClient
/* geninfo:
   filename  : golang.singingcat.net/apis/sensors/sensors.proto
   gopackage : golang.singingcat.net/apis/sensors
   importname: ai_0
   varname   : client_SensorServerClient_0
   clientname: SensorServerClient
   servername: SensorServerServer
   gscvname  : sensors.SensorServer
   lockname  : lock_SensorServerClient_0
   activename: active_SensorServerClient_0
*/

package sensors

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SensorServerClient_0 sync.Mutex
  client_SensorServerClient_0 SensorServerClient
)

func GetSensorServerClient() SensorServerClient { 
    if client_SensorServerClient_0 != nil {
        return client_SensorServerClient_0
    }

    lock_SensorServerClient_0.Lock() 
    if client_SensorServerClient_0 != nil {
       lock_SensorServerClient_0.Unlock()
       return client_SensorServerClient_0
    }

    client_SensorServerClient_0 = NewSensorServerClient(client.Connect("sensors.SensorServer"))
    lock_SensorServerClient_0.Unlock()
    return client_SensorServerClient_0
}

