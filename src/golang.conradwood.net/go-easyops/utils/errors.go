package utils

import (
	"fmt"
	fw "golang.conradwood.net/apis/framework"
	"google.golang.org/grpc/status"
	"strings"
)

// extracts the PRIVATE and possibly SENSITIVE debug error message from a string
func ErrorString(err error) string {
	st := status.Convert(err)
	s := "[STATUS] "
	deli := ""
	for _, a := range st.Details() {
		fmd, ok := a.(*fw.FrameworkMessageDetail)
		if !ok {
			s = s + fmt.Sprintf("\"%v\" ", a)
			continue
		}
		for _, ct := range fmd.CallTraces {
			if ct.Service != "" {
				spl := strings.SplitN(ct.Service, ".", 2)
				sn := ct.Service
				if len(spl) == 2 {
					sn = spl[1]
				}
				s = s + deli + fmt.Sprintf("%s.%s", sn, ct.Method)
			} else {
				s = s + deli + fmt.Sprintf("%s", ct.Message)
			}
			deli = "->"
		}

	}
	s = s + ": " + st.Message() + " [/STATUS]"
	return s

}
