module golang.conradwood.net/go-easyops

go 1.19

require (
	github.com/dustin/go-humanize v1.0.1
	github.com/go-sql-driver/mysql v1.7.1
	github.com/golang/protobuf v1.5.3
	github.com/lib/pq v1.10.9
	github.com/prometheus/client_golang v1.16.0
	github.com/prometheus/client_model v0.4.0
	golang.conradwood.net/apis/auth v1.1.2643
	golang.conradwood.net/apis/common v1.1.2643
	golang.conradwood.net/apis/echoservice v1.1.2495
	golang.conradwood.net/apis/errorlogger v1.1.2495
	golang.conradwood.net/apis/framework v1.1.2503
	golang.conradwood.net/apis/goeasyops v1.1.2643
	golang.conradwood.net/apis/objectstore v1.1.2503
	golang.conradwood.net/apis/registry v1.1.2503
	golang.org/x/net v0.17.0
	golang.org/x/sys v0.13.0
	golang.yacloud.eu/apis/session v1.1.2643
	golang.yacloud.eu/apis/urlcacher v1.1.2495
	google.golang.org/grpc v1.58.3
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	golang.conradwood.net/apis/autodeployer v1.1.2503 // indirect
	golang.conradwood.net/apis/commondeploy v1.1.2503 // indirect
	golang.conradwood.net/apis/deploymonkey v1.1.2503 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230711160842-782d3b101e98 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

//replace golang.conradwood.net/apis/goeasyops => ../apis/goeasyops
