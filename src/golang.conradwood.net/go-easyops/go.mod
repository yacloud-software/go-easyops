module golang.conradwood.net/go-easyops

go 1.21.1

require (
	github.com/dustin/go-humanize v1.0.1
	github.com/go-sql-driver/mysql v1.8.0
	github.com/golang/protobuf v1.5.4
	github.com/grafana/pyroscope-go v1.1.1
	github.com/lib/pq v1.10.9
	github.com/prometheus/client_golang v1.19.0
	github.com/prometheus/client_model v0.6.0
	golang.conradwood.net/apis/auth v1.1.2878
	golang.conradwood.net/apis/common v1.1.2878
	golang.conradwood.net/apis/echoservice v1.1.2878
	golang.conradwood.net/apis/errorlogger v1.1.2878
	golang.conradwood.net/apis/framework v1.1.2878
	golang.conradwood.net/apis/goeasyops v1.1.2878
	golang.conradwood.net/apis/objectstore v1.1.2878
	golang.conradwood.net/apis/registry v1.1.2878
	golang.org/x/net v0.22.0
	golang.org/x/sys v0.18.0
	golang.yacloud.eu/apis/session v1.1.2878
	golang.yacloud.eu/apis/urlcacher v1.1.2878
	golang.yacloud.eu/unixipc v0.1.26120
	google.golang.org/grpc v1.62.1
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/grafana/pyroscope-go/godeltaprof v0.1.7 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/prometheus/common v0.50.0 // indirect
	github.com/prometheus/procfs v0.13.0 // indirect
	golang.conradwood.net/apis/autodeployer v1.1.2878 // indirect
	golang.conradwood.net/apis/deploymonkey v1.1.2878 // indirect
	golang.conradwood.net/apis/grafanadata v1.1.2878 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.yacloud.eu/apis/fscache v1.1.2878 // indirect
	golang.yacloud.eu/apis/unixipc v1.1.2878 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

//replace golang.conradwood.net/apis/goeasyops => ../apis/goeasyops
