module golang.conradwood.net/tests

go 1.23.0

toolchain go1.23.3

replace golang.conradwood.net/go-easyops => ../go-easyops

replace golang.conradwood.net/apis/goeasyops => ../apis/goeasyops

replace golang.conradwood.net/apis/getestservice => ../apis/getestservice

require (
	golang.conradwood.net/apis/apitest v1.1.3359
	golang.conradwood.net/apis/auth v1.1.4006
	golang.conradwood.net/apis/common v1.1.4006
	golang.conradwood.net/apis/getestservice v1.1.4006
	golang.conradwood.net/apis/gitserver v1.1.3995
	golang.conradwood.net/apis/goeasyops v1.1.4006
	golang.conradwood.net/apis/helloworld v1.1.4006
	golang.conradwood.net/apis/registry v1.1.4006
	golang.conradwood.net/go-easyops v0.1.34140
	google.golang.org/grpc v1.74.2
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/grafana/pyroscope-go v1.2.3 // indirect
	github.com/grafana/pyroscope-go/godeltaprof v0.1.8 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.23.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.65.0 // indirect
	github.com/prometheus/procfs v0.17.0 // indirect
	golang.conradwood.net/apis/autodeployer v1.1.4006 // indirect
	golang.conradwood.net/apis/certmanager v1.1.4006 // indirect
	golang.conradwood.net/apis/deploymonkey v1.1.4006 // indirect
	golang.conradwood.net/apis/errorlogger v1.1.4006 // indirect
	golang.conradwood.net/apis/framework v1.1.4006 // indirect
	golang.conradwood.net/apis/grafanadata v1.1.4006 // indirect
	golang.conradwood.net/apis/h2gproxy v1.1.4006 // indirect
	golang.conradwood.net/apis/objectstore v1.1.4006 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	golang.yacloud.eu/apis/autodeployer2 v1.1.4006 // indirect
	golang.yacloud.eu/apis/faultindicator v1.1.4006 // indirect
	golang.yacloud.eu/apis/fscache v1.1.4006 // indirect
	golang.yacloud.eu/apis/session v1.1.4006 // indirect
	golang.yacloud.eu/apis/unixipc v1.1.4006 // indirect
	golang.yacloud.eu/apis/urlcacher v1.1.4006 // indirect
	golang.yacloud.eu/unixipc v0.1.31725 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
