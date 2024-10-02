package shared

import (
	"fmt"
	"strings"

	ge "golang.conradwood.net/apis/goeasyops"
)

func LocalState2string(ls LocalState) string {
	if isNil(ls) {
		return "[no localstate]"
	}
	s := "LocalState:\n"
	s = s + fmt.Sprintf("  User            : %s\n", UserIDString(ls.User()))
	s = s + fmt.Sprintf("  Sudo-User       : %s\n", UserIDString(ls.SudoUser()))
	s = s + fmt.Sprintf("  CreatorService  : %s\n", UserIDString(ls.CreatorService()))
	s = s + fmt.Sprintf("  CallingService  : %s\n", UserIDString(ls.CallingService()))
	s = s + fmt.Sprintf("  Services        : %s\n", comma_delimited(ls.Services(), func(x *ge.ServiceTrace) string { return x.ID }))
	s = s + fmt.Sprintf("  Experiments     : %s\n", comma_delimited(ls.Experiments(), func(x *ge.Experiment) string { return x.Name }))
	s = s + fmt.Sprintf("  Debug           : %v\n", ls.Debug())
	s = s + fmt.Sprintf("  Trace           : %v\n", ls.Trace())
	s = s + fmt.Sprintf("  RequestID       : %s\n", ls.RequestID())
	s = s + fmt.Sprintf("  RoutingTags     : %v\n", ls.RoutingTags())

	return s

}
func ContextProto2string(prefix string, x *ge.InContext) string {
	s := "InContext:\n"
	s = s + "  Immutable:\n"
	s = s + Imctx2string("    ", x.ImCtx)
	s = s + "\n  Mutable:\n"
	s = s + Mctx2string("    ", x.MCtx)
	return AddPrefixToLines(prefix, s)
}
func Imctx2string(prefix string, x *ge.ImmutableContext) string {
	s := ""
	s = s + fmt.Sprintf("RequestID      : %s\n", x.RequestID)
	s = s + fmt.Sprintf("User           : %s\n", UserIDString(x.User))
	s = s + fmt.Sprintf("SudoUser       : %s\n", UserIDString(x.SudoUser))
	s = s + fmt.Sprintf("CreatorService : %s", UserIDString(x.CreatorService))
	return AddPrefixToLines(prefix, s)
}
func Mctx2string(prefix string, x *ge.MutableContext) string {
	s := ""
	s = s + fmt.Sprintf("CallingService : %s\n", UserIDString(x.CallingService))
	s = s + fmt.Sprintf("Serviceids     : %d\n", len(x.ServiceIDs))
	s = s + fmt.Sprintf("Debug          : %v\n", x.Debug)
	s = s + fmt.Sprintf("Trace          : %v\n", x.Trace)
	tag_s := ""
	if x.Tags != nil {
		tag_s = "Got tags"
	}
	s = s + fmt.Sprintf("Tags           : %v\n", tag_s)
	s = s + fmt.Sprintf("Experiments    : %s", comma_delimited(x.Experiments, func(x *ge.Experiment) string { return x.Name }))
	return AddPrefixToLines(prefix, s)
}

func comma_delimited[K interface{}](objects []K, f func(K) string) string {
	deli := ""
	xx := ""
	if len(objects) > 0 {
		xx = ": "
	}
	res := fmt.Sprintf("%d%s", len(objects), xx)
	for _, x := range objects {
		xs := f(x)
		res = res + deli + fmt.Sprintf("\"%s\"", xs)
		deli = ", "
	}
	return res
}

func AddPrefixToLines(prefix, txt string) string {
	res := ""
	for _, line := range strings.Split(txt, "\n") {
		nl := prefix + line + "\n"
		res = res + nl
	}
	return strings.TrimSuffix(res, "\n")
}
