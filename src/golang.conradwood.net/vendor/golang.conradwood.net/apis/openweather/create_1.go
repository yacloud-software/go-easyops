// client create: OpenWeatherServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/openweather/openweather.proto
   gopackage : golang.conradwood.net/apis/openweather
   importname: ai_0
   varname   : client_OpenWeatherServiceClient_0
   clientname: OpenWeatherServiceClient
   servername: OpenWeatherServiceServer
   gscvname  : openweather.OpenWeatherService
   lockname  : lock_OpenWeatherServiceClient_0
   activename: active_OpenWeatherServiceClient_0
*/

package openweather

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_OpenWeatherServiceClient_0 sync.Mutex
  client_OpenWeatherServiceClient_0 OpenWeatherServiceClient
)

func GetOpenWeatherClient() OpenWeatherServiceClient { 
    if client_OpenWeatherServiceClient_0 != nil {
        return client_OpenWeatherServiceClient_0
    }

    lock_OpenWeatherServiceClient_0.Lock() 
    if client_OpenWeatherServiceClient_0 != nil {
       lock_OpenWeatherServiceClient_0.Unlock()
       return client_OpenWeatherServiceClient_0
    }

    client_OpenWeatherServiceClient_0 = NewOpenWeatherServiceClient(client.Connect("openweather.OpenWeatherService"))
    lock_OpenWeatherServiceClient_0.Unlock()
    return client_OpenWeatherServiceClient_0
}

func GetOpenWeatherServiceClient() OpenWeatherServiceClient { 
    if client_OpenWeatherServiceClient_0 != nil {
        return client_OpenWeatherServiceClient_0
    }

    lock_OpenWeatherServiceClient_0.Lock() 
    if client_OpenWeatherServiceClient_0 != nil {
       lock_OpenWeatherServiceClient_0.Unlock()
       return client_OpenWeatherServiceClient_0
    }

    client_OpenWeatherServiceClient_0 = NewOpenWeatherServiceClient(client.Connect("openweather.OpenWeatherService"))
    lock_OpenWeatherServiceClient_0.Unlock()
    return client_OpenWeatherServiceClient_0
}

