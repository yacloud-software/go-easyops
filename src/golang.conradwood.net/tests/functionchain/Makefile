BINARY=$(shell basename `/usr/bin/pwd`)
.PHONY: all
all:
	echo "makefile for ${BINARY} binary"
	go install ${LDFLAGS} ${BINARY}.go `ls -1 *.go|grep -v ${BINARY}.go`
