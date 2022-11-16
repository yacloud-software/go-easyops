package authremote

import (
	"context"
	rc "golang.conradwood.net/apis/rpcinterceptor"
	"time"
)

// TODO: add tags, timeout, etc. this is only used on "standalone" (thus not network environments)
func standalone_ContextWithTimeoutAndTags(t time.Duration, rt *rc.CTXRoutingTags) context.Context {
	return context.TODO()
}
