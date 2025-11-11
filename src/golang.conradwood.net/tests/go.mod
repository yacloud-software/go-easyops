module golang.conradwood.net/tests

go 1.24.0

replace golang.conradwood.net/go-easyops => ../go-easyops

replace golang.conradwood.net/apis/goeasyops => ../apis/goeasyops

replace golang.conradwood.net/apis/getestservice => ../apis/getestservice

require (
	golang.conradwood.net/apis/apitest v1.1.3359
	golang.conradwood.net/apis/auth v1.1.4186
	golang.conradwood.net/apis/common v1.1.4186
	golang.conradwood.net/apis/getestservice v1.1.4186
	golang.conradwood.net/apis/gitserver v1.1.4186
	golang.conradwood.net/apis/goeasyops v1.1.4186
	golang.conradwood.net/apis/helloworld v1.1.4186
	golang.conradwood.net/apis/registry v1.1.4186
	golang.conradwood.net/go-easyops v0.1.34140
	google.golang.org/grpc v1.76.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/grafana/pyroscope-go v1.2.7 // indirect
	github.com/grafana/pyroscope-go/godeltaprof v0.1.9 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.2 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	golang.conradwood.net/apis/autodeployer v1.1.4186 // indirect
	golang.conradwood.net/apis/certmanager v1.1.4186 // indirect
	golang.conradwood.net/apis/deploymonkey v1.1.4186 // indirect
	golang.conradwood.net/apis/errorlogger v1.1.4186 // indirect
	golang.conradwood.net/apis/framework v1.1.4186 // indirect
	golang.conradwood.net/apis/grafanadata v1.1.4186 // indirect
	golang.conradwood.net/apis/h2gproxy v1.1.4186 // indirect
	golang.conradwood.net/apis/objectstore v1.1.4186 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	golang.yacloud.eu/apis/autodeployer2 v1.1.4186 // indirect
	golang.yacloud.eu/apis/faultindicator v1.1.4186 // indirect
	golang.yacloud.eu/apis/fscache v1.1.4186 // indirect
	golang.yacloud.eu/apis/session v1.1.4186 // indirect
	golang.yacloud.eu/apis/unixipc v1.1.4186 // indirect
	golang.yacloud.eu/apis/urlcacher v1.1.4186 // indirect
	golang.yacloud.eu/unixipc v0.1.31725 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
