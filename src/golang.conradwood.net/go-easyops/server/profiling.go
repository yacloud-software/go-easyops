package server

import (
	"flag"
	"fmt"

	"github.com/grafana/pyroscope-go"
	"golang.conradwood.net/go-easyops/appinfo"
)

const (
	pyroscope_port = 3900
)

var (
	start_pyroscope = flag.Bool("ge_start_pyroscope", false, "use ge-pyroscope:3900 to record hotspots")
	pyroscope_host  = flag.String("ge_pyroscope_host", "ge-pyroscope", "use `hostname`:3900 to record hotspots")
)

func start_profiling(sd *serverDef) {
	if *start_pyroscope {
		pyro_url := fmt.Sprintf("http://%s:%d", *pyroscope_host, pyroscope_port)
		pyroscope.Start(pyroscope.Config{
			ApplicationName: sd.name,
			ServerAddress:   pyro_url,
			//			Logger:          nil,
			Logger: pyroscope.StandardLogger,
			Tags: map[string]string{
				"service_user_id": fmt.Sprintf("%s", sd.service_user_id),
				"artefactid":      fmt.Sprintf("%d", appinfo.AppInfo().ArtefactID),
				"version":         fmt.Sprintf("%d", appinfo.AppInfo().Number),
			},
		})
		fmt.Printf("[go-easyops] pyroscope started, connecting to host %s\n", pyro_url)
	}

}
