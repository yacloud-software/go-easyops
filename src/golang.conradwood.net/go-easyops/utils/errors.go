package utils

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang.conradwood.net/apis/common"
	fw "golang.conradwood.net/apis/framework"
	goe "golang.conradwood.net/apis/goeasyops"
	"google.golang.org/grpc/status"
	proto2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
	"strings"
)

// extracts the PRIVATE and possibly SENSITIVE debug error message from a string
// obsolete - use errors.ErrorString(err)
// the reason this is so convoluted with different types, is that different versions of grpc
// encapsulate status details in different messages.
func ErrorString(err error) string {
	st := status.Convert(err)
	s := "[STATUS] "
	deli := ""
	var cstatus *common.Status
	var gel *goe.GRPCErrorList
	for _, a := range st.Details() {
		unknown := true

		proto2m := a.(proto2.Message)
		msgname := proto2.MessageName(proto2m)
		//	msg := proto2m.ProtoReflect()
		pv1 := protoadapt.MessageV1Of(proto2m)
		//fmt.Printf("Proto2 (%s): %#v %v %v\n", msgname, proto2m, msg, pv1)
		if msgname == "goeasyops.GRPCErrorList" {
			xgel, ok := pv1.(*goe.GRPCErrorList)
			if ok {
				gel = xgel
				s = s + deli + ge2string(xgel)
				continue
			}
		} else if msgname == "common.Status" {
			st, ok := pv1.(*common.Status)
			if ok {
				cstatus = st
				continue
			}
		}

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

		x, ok := a.(goe.GRPCError)
		if unknown && ok {
			unknown = false
			s = s + deli + fmt.Sprintf("CALLTRACE: %v", x)
		}

		x2, ok := a.(*fw.CallTrace)
		if unknown && ok {
			unknown = false
			s = s + deli + fmt.Sprintf("CALLTRACE: %v", x2)
		}

		proto, ok := a.(proto.Message)
		if unknown && ok {
			unknown = false
			s = s + deli + "proto:" + proto.String()
		}

		deli = "->"

	}
	s = s + ": " + st.Message() + " [/STATUS]"
	if cstatus == nil || gel == nil {
		return s
	}
	s = fmt.Sprintf("%d(%s): ", cstatus.ErrorCode, cstatus.ErrorDescription) + ge2string(gel)
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
