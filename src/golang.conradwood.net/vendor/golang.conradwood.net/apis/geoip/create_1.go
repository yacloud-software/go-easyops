// client create: GeoIPServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/geoip/geoip.proto
   gopackage : golang.conradwood.net/apis/geoip
   importname: ai_0
   varname   : client_GeoIPServiceClient_0
   clientname: GeoIPServiceClient
   servername: GeoIPServiceServer
   gscvname  : geoip.GeoIPService
   lockname  : lock_GeoIPServiceClient_0
   activename: active_GeoIPServiceClient_0
*/

package geoip

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_GeoIPServiceClient_0 sync.Mutex
  client_GeoIPServiceClient_0 GeoIPServiceClient
)

func GetGeoIPClient() GeoIPServiceClient { 
    if client_GeoIPServiceClient_0 != nil {
        return client_GeoIPServiceClient_0
    }

    lock_GeoIPServiceClient_0.Lock() 
    if client_GeoIPServiceClient_0 != nil {
       lock_GeoIPServiceClient_0.Unlock()
       return client_GeoIPServiceClient_0
    }

    client_GeoIPServiceClient_0 = NewGeoIPServiceClient(client.Connect("geoip.GeoIPService"))
    lock_GeoIPServiceClient_0.Unlock()
    return client_GeoIPServiceClient_0
}

func GetGeoIPServiceClient() GeoIPServiceClient { 
    if client_GeoIPServiceClient_0 != nil {
        return client_GeoIPServiceClient_0
    }

    lock_GeoIPServiceClient_0.Lock() 
    if client_GeoIPServiceClient_0 != nil {
       lock_GeoIPServiceClient_0.Unlock()
       return client_GeoIPServiceClient_0
    }

    client_GeoIPServiceClient_0 = NewGeoIPServiceClient(client.Connect("geoip.GeoIPService"))
    lock_GeoIPServiceClient_0.Unlock()
    return client_GeoIPServiceClient_0
}

