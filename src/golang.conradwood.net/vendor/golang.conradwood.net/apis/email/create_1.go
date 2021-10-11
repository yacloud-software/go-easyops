// client create: EmailServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/email/email.proto
   gopackage : golang.conradwood.net/apis/email
   importname: ai_0
   varname   : client_EmailServiceClient_0
   clientname: EmailServiceClient
   servername: EmailServiceServer
   gscvname  : email.EmailService
   lockname  : lock_EmailServiceClient_0
   activename: active_EmailServiceClient_0
*/

package email

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_EmailServiceClient_0 sync.Mutex
  client_EmailServiceClient_0 EmailServiceClient
)

func GetEmailClient() EmailServiceClient { 
    if client_EmailServiceClient_0 != nil {
        return client_EmailServiceClient_0
    }

    lock_EmailServiceClient_0.Lock() 
    if client_EmailServiceClient_0 != nil {
       lock_EmailServiceClient_0.Unlock()
       return client_EmailServiceClient_0
    }

    client_EmailServiceClient_0 = NewEmailServiceClient(client.Connect("email.EmailService"))
    lock_EmailServiceClient_0.Unlock()
    return client_EmailServiceClient_0
}

func GetEmailServiceClient() EmailServiceClient { 
    if client_EmailServiceClient_0 != nil {
        return client_EmailServiceClient_0
    }

    lock_EmailServiceClient_0.Lock() 
    if client_EmailServiceClient_0 != nil {
       lock_EmailServiceClient_0.Unlock()
       return client_EmailServiceClient_0
    }

    client_EmailServiceClient_0 = NewEmailServiceClient(client.Connect("email.EmailService"))
    lock_EmailServiceClient_0.Unlock()
    return client_EmailServiceClient_0
}

