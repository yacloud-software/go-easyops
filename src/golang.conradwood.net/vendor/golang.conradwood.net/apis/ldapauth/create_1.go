// client create: LdapAuthClient
/* geninfo:
   filename  : golang.conradwood.net/apis/ldapauth/ldapauth.proto
   gopackage : golang.conradwood.net/apis/ldapauth
   importname: ai_0
   varname   : client_LdapAuthClient_0
   clientname: LdapAuthClient
   servername: LdapAuthServer
   gscvname  : ldapauth.LdapAuth
   lockname  : lock_LdapAuthClient_0
   activename: active_LdapAuthClient_0
*/

package ldapauth

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_LdapAuthClient_0 sync.Mutex
  client_LdapAuthClient_0 LdapAuthClient
)

func GetLdapAuthClient() LdapAuthClient { 
    if client_LdapAuthClient_0 != nil {
        return client_LdapAuthClient_0
    }

    lock_LdapAuthClient_0.Lock() 
    if client_LdapAuthClient_0 != nil {
       lock_LdapAuthClient_0.Unlock()
       return client_LdapAuthClient_0
    }

    client_LdapAuthClient_0 = NewLdapAuthClient(client.Connect("ldapauth.LdapAuth"))
    lock_LdapAuthClient_0.Unlock()
    return client_LdapAuthClient_0
}

