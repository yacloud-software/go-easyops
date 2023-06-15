package auth

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"gopkg.in/yaml.v2"
	"sync"
)

type serviceToUserIDMap struct {
	mapping map[string]string // servicename->userid
}

const (
	service_mapping_filename = "/opt/yacloud/config/service_map.yaml"
)

var (
	service_mapping         *serviceToUserIDMap
	svcmaplock              sync.Mutex
	default_service_mapping = map[string]string{
		"h2gproxy.H2GProxyService":              "37",
		"jsonapimultiplexer.JSONApiMultiplexer": "59",
		"repobuilder.RepoBuilder":               "3539",
		"weblogin.Weblogin":                     "43",
		"artefact.ArtefactService":              "998",
	}
)

/*
sometimes a service needs to verify if it is being called by a specific service. This often implies permissions to access certain privileged bits of information. The assumption is, that service A authenticates a user and calls service B, either immediately after or some time later. In this case service B "trusts" service A. The security implication of this model is, that service B must be able to ensure service A really is who they say they are (the auth server signature should be used for this purpose) and service A has not been replaced with a different service of the same name.
For this purpose, in this case, the service to userid mappings are hardcoded into the file so to match the "yacloud" default. If someone wishes to run their own yacloud the mapping can be overriden. This then is not a programmatic option, but a configuration (administrator) option.
A file in /opt/yacloud/config/service_map.yaml, if exists, will be parsed on startup and used to provide this information.
Any lookup for a servicename that does not exist will lead to a panic() (because it is a fatal error!).
The intention of this function is to provide a means to create a common method of looking up this information, so that, in future, perhaps a good and secure way can be found to automatically map this through a combination of registry/auth-server lookups or similar.
*/
func GetServiceIDByName(servicename string) string {
	svc_to_user_load_mapping()
	uid, found := service_mapping.mapping[servicename]
	if found {
		return uid
	}
	panic(fmt.Sprintf("[go-easyops] Application requested service \"%s\", which is not mapped to a userid", servicename))
	return ""
}
func svc_to_user_load_mapping() {
	if service_mapping != nil {
		return
	}
	svcmaplock.Lock()
	defer svcmaplock.Unlock()
	if service_mapping != nil {
		return
	}
	res := &serviceToUserIDMap{
		mapping: default_service_mapping,
	}
	b, err := utils.ReadFile(service_mapping_filename)
	if err == nil {
		res = &serviceToUserIDMap{}
		err = yaml.Unmarshal(b, res)
		if err != nil {
			panic(fmt.Sprintf("File %s cannot be parsed: %s\n", service_mapping_filename, err))
		}
	}
	service_mapping = res
}
