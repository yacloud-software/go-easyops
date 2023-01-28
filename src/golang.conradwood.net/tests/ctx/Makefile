.PHONY: server client
LDFLAGS=-ldflags "-X golang.conradwood.net/go-easyops/appinfo.LD_Number=56 -X golang.conradwood.net/go-easyops/appinfo.LD_Timestamp=78" 
all:
	go install ${LDFLAGS} test-ctx.go `ls -1 *.go|grep -v test-ctx.go`

