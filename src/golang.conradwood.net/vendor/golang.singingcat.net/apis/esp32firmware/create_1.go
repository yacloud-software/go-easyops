// client create: ESP32Client
/* geninfo:
   filename  : golang.singingcat.net/apis/esp32firmware/esp32-firmware.proto
   gopackage : golang.singingcat.net/apis/esp32firmware
   importname: ai_0
   varname   : client_ESP32Client_0
   clientname: ESP32Client
   servername: ESP32Server
   gscvname  : esp32firmware.ESP32
   lockname  : lock_ESP32Client_0
   activename: active_ESP32Client_0
*/

package esp32firmware

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_ESP32Client_0 sync.Mutex
  client_ESP32Client_0 ESP32Client
)

func GetESP32Client() ESP32Client { 
    if client_ESP32Client_0 != nil {
        return client_ESP32Client_0
    }

    lock_ESP32Client_0.Lock() 
    if client_ESP32Client_0 != nil {
       lock_ESP32Client_0.Unlock()
       return client_ESP32Client_0
    }

    client_ESP32Client_0 = NewESP32Client(client.Connect("esp32firmware.ESP32"))
    lock_ESP32Client_0.Unlock()
    return client_ESP32Client_0
}

