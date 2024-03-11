package server

import (
	"fmt"
	"github.com/grafana/pyroscope-go"
	"golang.conradwood.net/go-easyops/appinfo"
)

func start_profiling(sd *serverDef) {
	pyroscope.Start(pyroscope.Config{
		ApplicationName: sd.name,
		ServerAddress:   "http://ge-pyroscope:3900",
		Logger:          nil,
		Tags: map[string]string{
			"service_user_id": fmt.Sprintf("%s", sd.service_user_id),
			"artefactid":      fmt.Sprintf("%d", appinfo.AppInfo().ArtefactID),
			"version":         fmt.Sprintf("%d", appinfo.AppInfo().Number),
		},
	})

}
