// client create: HTTPKPITrackerClient
/* geninfo:
   filename  : golang.conradwood.net/apis/httpkpi/httpkpi.proto
   gopackage : golang.conradwood.net/apis/httpkpi
   importname: ai_0
   varname   : client_HTTPKPITrackerClient_0
   clientname: HTTPKPITrackerClient
   servername: HTTPKPITrackerServer
   gscvname  : httpkpi.HTTPKPITracker
   lockname  : lock_HTTPKPITrackerClient_0
   activename: active_HTTPKPITrackerClient_0
*/

package httpkpi

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_HTTPKPITrackerClient_0 sync.Mutex
  client_HTTPKPITrackerClient_0 HTTPKPITrackerClient
)

func GetHTTPKPITrackerClient() HTTPKPITrackerClient { 
    if client_HTTPKPITrackerClient_0 != nil {
        return client_HTTPKPITrackerClient_0
    }

    lock_HTTPKPITrackerClient_0.Lock() 
    if client_HTTPKPITrackerClient_0 != nil {
       lock_HTTPKPITrackerClient_0.Unlock()
       return client_HTTPKPITrackerClient_0
    }

    client_HTTPKPITrackerClient_0 = NewHTTPKPITrackerClient(client.Connect("httpkpi.HTTPKPITracker"))
    lock_HTTPKPITrackerClient_0.Unlock()
    return client_HTTPKPITrackerClient_0
}

