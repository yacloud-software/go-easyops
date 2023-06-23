package authremote

import (
	"context"
	ge "golang.conradwood.net/apis/goeasyops"
	"time"
)

// TODO: add tags, timeout, etc. this is only used on "standalone" (thus not network environments)
func standalone_ContextWithTimeoutAndTags(t time.Duration, rt *ge.CTXRoutingTags) context.Context {
	return context.TODO()
}
