package utils

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	fw "golang.conradwood.net/apis/framework"
	goe "golang.conradwood.net/apis/goeasyops"
	"google.golang.org/grpc/status"
	"reflect"
	"strings"
)

// extracts the PRIVATE and possibly SENSITIVE debug error message from a string
// obsolete - use errors.ErrorString(err)
func ErrorString(err error) string {
	st := status.Convert(err)
	s := "[STATUS] "
	deli := ""
	for _, a := range st.Details() {

		unknown := true

		fmd, ok := a.(*fw.FrameworkMessageDetail)
		if ok {
			unknown = false
			s = s + deli + fmd2string(fmd)
		}

		ge, ok := a.(*goe.GRPCErrorList)
		if unknown && ok {
			unknown = false
			s = s + deli + ge2string(ge)
		}

		ge2, ok := a.(goe.GRPCErrorList)
		if unknown && ok {
			unknown = false
			s = s + deli + ge2string(&ge2)
		}

		proto, ok := a.(proto.Message)
		if unknown && ok {
			unknown = false
			s = s + deli + "proto:" + proto.String()
		}

		if unknown {
			s = s + fmt.Sprintf("%s", reflect.TypeOf(a))
			s = s + fmt.Sprintf("\"%v\" ", a)
		}
		deli = "->"

	}
	s = s + ": " + st.Message() + " [/STATUS]"
	return s

}

func fmd2string(fmd *fw.FrameworkMessageDetail) string {
	s := ""
	for _, ct := range fmd.CallTraces {
		if ct.Service != "" {
			spl := strings.SplitN(ct.Service, ".", 2)
			sn := ct.Service
			if len(spl) == 2 {
				sn = spl[1]
			}
			s = fmt.Sprintf("%s.%s", sn, ct.Method)
		} else {
			s = fmt.Sprintf("%s", ct.Message)
		}
	}
	return s
}

func ge2string(fmd *goe.GRPCErrorList) string {
	s := ""
	for _, ct := range fmd.Errors {
		if ct.ServiceName != "" {
			spl := strings.SplitN(ct.ServiceName, ".", 2)
			sn := ct.ServiceName
			if len(spl) == 2 {
				sn = spl[1]
			}
			s = fmt.Sprintf("%s.%s", sn, ct.MethodName)
		} else {
			s = fmt.Sprintf("%s", ct.LogMessage)
		}
	}
	return s
}
