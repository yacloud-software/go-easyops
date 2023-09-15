
package getestservice

import (
    "sync"
    "golang.conradwood.net/go-easyops/protos"
)
var (
    package_lock sync.Mutex
    services []string
)

func AddService(name string) {
   package_lock.Lock()
   services = append(services,name)
   package_lock.Unlock()
}
func IsHealthy() ( bool,error) {
   for _,s:=range services {
      v,err:=protos.IsHealthy(s)
      if err != nil {
          return false,err
      }
      if !v {
          return false,nil
      }
   }
   return true,nil
}
