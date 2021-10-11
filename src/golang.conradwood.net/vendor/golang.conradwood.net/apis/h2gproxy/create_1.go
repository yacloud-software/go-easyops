// client create: DownloadStreamerClient
/* geninfo:
   filename  : golang.conradwood.net/apis/h2gproxy/h2gproxy.proto
   gopackage : golang.conradwood.net/apis/h2gproxy
   importname: ai_0
   varname   : client_DownloadStreamerClient_0
   clientname: DownloadStreamerClient
   servername: DownloadStreamerServer
   gscvname  : h2gproxy.DownloadStreamer
   lockname  : lock_DownloadStreamerClient_0
   activename: active_DownloadStreamerClient_0
*/

package h2gproxy

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_DownloadStreamerClient_0 sync.Mutex
  client_DownloadStreamerClient_0 DownloadStreamerClient
)

func GetDownloadStreamerClient() DownloadStreamerClient { 
    if client_DownloadStreamerClient_0 != nil {
        return client_DownloadStreamerClient_0
    }

    lock_DownloadStreamerClient_0.Lock() 
    if client_DownloadStreamerClient_0 != nil {
       lock_DownloadStreamerClient_0.Unlock()
       return client_DownloadStreamerClient_0
    }

    client_DownloadStreamerClient_0 = NewDownloadStreamerClient(client.Connect("h2gproxy.DownloadStreamer"))
    lock_DownloadStreamerClient_0.Unlock()
    return client_DownloadStreamerClient_0
}

