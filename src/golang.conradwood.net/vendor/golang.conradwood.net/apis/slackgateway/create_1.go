// client create: SlackGatewayClient
/* geninfo:
   filename  : golang.conradwood.net/apis/slackgateway/slackgateway.proto
   gopackage : golang.conradwood.net/apis/slackgateway
   importname: ai_0
   varname   : client_SlackGatewayClient_0
   clientname: SlackGatewayClient
   servername: SlackGatewayServer
   gscvname  : slackgateway.SlackGateway
   lockname  : lock_SlackGatewayClient_0
   activename: active_SlackGatewayClient_0
*/

package slackgateway

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SlackGatewayClient_0 sync.Mutex
  client_SlackGatewayClient_0 SlackGatewayClient
)

func GetSlackGatewayClient() SlackGatewayClient { 
    if client_SlackGatewayClient_0 != nil {
        return client_SlackGatewayClient_0
    }

    lock_SlackGatewayClient_0.Lock() 
    if client_SlackGatewayClient_0 != nil {
       lock_SlackGatewayClient_0.Unlock()
       return client_SlackGatewayClient_0
    }

    client_SlackGatewayClient_0 = NewSlackGatewayClient(client.Connect("slackgateway.SlackGateway"))
    lock_SlackGatewayClient_0.Unlock()
    return client_SlackGatewayClient_0
}

