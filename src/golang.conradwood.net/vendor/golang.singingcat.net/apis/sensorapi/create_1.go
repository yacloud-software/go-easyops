// client create: SensorAPIServiceClient
/* geninfo:
   filename  : golang.singingcat.net/apis/sensorapi/sensorapi.proto
   gopackage : golang.singingcat.net/apis/sensorapi
   importname: ai_0
   varname   : client_SensorAPIServiceClient_0
   clientname: SensorAPIServiceClient
   servername: SensorAPIServiceServer
   gscvname  : sensorapi.SensorAPIService
   lockname  : lock_SensorAPIServiceClient_0
   activename: active_SensorAPIServiceClient_0
*/

package sensorapi

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SensorAPIServiceClient_0 sync.Mutex
  client_SensorAPIServiceClient_0 SensorAPIServiceClient
)

func GetSensorAPIClient() SensorAPIServiceClient { 
    if client_SensorAPIServiceClient_0 != nil {
        return client_SensorAPIServiceClient_0
    }

    lock_SensorAPIServiceClient_0.Lock() 
    if client_SensorAPIServiceClient_0 != nil {
       lock_SensorAPIServiceClient_0.Unlock()
       return client_SensorAPIServiceClient_0
    }

    client_SensorAPIServiceClient_0 = NewSensorAPIServiceClient(client.Connect("sensorapi.SensorAPIService"))
    lock_SensorAPIServiceClient_0.Unlock()
    return client_SensorAPIServiceClient_0
}

func GetSensorAPIServiceClient() SensorAPIServiceClient { 
    if client_SensorAPIServiceClient_0 != nil {
        return client_SensorAPIServiceClient_0
    }

    lock_SensorAPIServiceClient_0.Lock() 
    if client_SensorAPIServiceClient_0 != nil {
       lock_SensorAPIServiceClient_0.Unlock()
       return client_SensorAPIServiceClient_0
    }

    client_SensorAPIServiceClient_0 = NewSensorAPIServiceClient(client.Connect("sensorapi.SensorAPIService"))
    lock_SensorAPIServiceClient_0.Unlock()
    return client_SensorAPIServiceClient_0
}

