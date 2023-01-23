module golang.conradwood.net/go-easyops

go 1.19

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/go-sql-driver/mysql v1.7.0
	github.com/golang/protobuf v1.5.2
	github.com/lib/pq v1.10.7
	github.com/prometheus/client_golang v1.14.0
	github.com/prometheus/client_model v0.3.0
	golang.conradwood.net/apis/auth v1.1.2073
	golang.conradwood.net/apis/common v1.1.2073
	golang.conradwood.net/apis/echoservice v1.1.2073
	golang.conradwood.net/apis/errorlogger v1.1.2121
	golang.conradwood.net/apis/framework v1.1.2073
	golang.conradwood.net/apis/goeasyops v0.0.0-00010101000000-000000000000
	golang.conradwood.net/apis/logservice v1.1.2073
	golang.conradwood.net/apis/objectstore v1.1.2073
	golang.conradwood.net/apis/registry v1.1.2073
	golang.conradwood.net/apis/rpcinterceptor v1.1.2073
	golang.org/x/net v0.4.0
	golang.org/x/sys v0.3.0
	golang.yacloud.eu/apis/urlcacher v1.1.2073
	google.golang.org/grpc v1.51.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/common v0.39.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	golang.conradwood.net/apis/autodeployer v1.1.2073 // indirect
	golang.conradwood.net/apis/deploymonkey v1.1.2073 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace golang.conradwood.net/apis/goeasyops => ../apis/goeasyops
