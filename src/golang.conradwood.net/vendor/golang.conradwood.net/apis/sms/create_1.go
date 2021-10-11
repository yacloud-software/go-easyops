// client create: SMSServiceClient
/* geninfo:
   filename  : golang.conradwood.net/apis/sms/sms.proto
   gopackage : golang.conradwood.net/apis/sms
   importname: ai_0
   varname   : client_SMSServiceClient_0
   clientname: SMSServiceClient
   servername: SMSServiceServer
   gscvname  : sms.SMSService
   lockname  : lock_SMSServiceClient_0
   activename: active_SMSServiceClient_0
*/

package sms

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_SMSServiceClient_0 sync.Mutex
  client_SMSServiceClient_0 SMSServiceClient
)

func GetSMSClient() SMSServiceClient { 
    if client_SMSServiceClient_0 != nil {
        return client_SMSServiceClient_0
    }

    lock_SMSServiceClient_0.Lock() 
    if client_SMSServiceClient_0 != nil {
       lock_SMSServiceClient_0.Unlock()
       return client_SMSServiceClient_0
    }

    client_SMSServiceClient_0 = NewSMSServiceClient(client.Connect("sms.SMSService"))
    lock_SMSServiceClient_0.Unlock()
    return client_SMSServiceClient_0
}

func GetSMSServiceClient() SMSServiceClient { 
    if client_SMSServiceClient_0 != nil {
        return client_SMSServiceClient_0
    }

    lock_SMSServiceClient_0.Lock() 
    if client_SMSServiceClient_0 != nil {
       lock_SMSServiceClient_0.Unlock()
       return client_SMSServiceClient_0
    }

    client_SMSServiceClient_0 = NewSMSServiceClient(client.Connect("sms.SMSService"))
    lock_SMSServiceClient_0.Unlock()
    return client_SMSServiceClient_0
}

