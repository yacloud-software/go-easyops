// client create: CalendarWrapperClient
/* geninfo:
   filename  : golang.conradwood.net/apis/calendarwrapper/calendarwrapper.proto
   gopackage : golang.conradwood.net/apis/calendarwrapper
   importname: ai_0
   varname   : client_CalendarWrapperClient_0
   clientname: CalendarWrapperClient
   servername: CalendarWrapperServer
   gscvname  : calendarwrapper.CalendarWrapper
   lockname  : lock_CalendarWrapperClient_0
   activename: active_CalendarWrapperClient_0
*/

package calendarwrapper

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_CalendarWrapperClient_0 sync.Mutex
  client_CalendarWrapperClient_0 CalendarWrapperClient
)

func GetCalendarWrapperClient() CalendarWrapperClient { 
    if client_CalendarWrapperClient_0 != nil {
        return client_CalendarWrapperClient_0
    }

    lock_CalendarWrapperClient_0.Lock() 
    if client_CalendarWrapperClient_0 != nil {
       lock_CalendarWrapperClient_0.Unlock()
       return client_CalendarWrapperClient_0
    }

    client_CalendarWrapperClient_0 = NewCalendarWrapperClient(client.Connect("calendarwrapper.CalendarWrapper"))
    lock_CalendarWrapperClient_0.Unlock()
    return client_CalendarWrapperClient_0
}

