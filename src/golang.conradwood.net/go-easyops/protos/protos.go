/*
this package is imported by protos
*/
package protos

import (
	// take care with importing go-easyops here, it easily becomes circular in non-obvious ways
	_ "golang.conradwood.net/go-easyops/client"
)

// TODO: implement this at some point when all easyops are updated and we can cange the proto compile
func IsHealthy(name string) (bool, error) {
	return true, nil
}
