package server

import (
	"context"
	"fmt"
	rc "golang.conradwood.net/apis/rpcinterceptor"
	"golang.conradwood.net/go-easyops/cache"
	"time"
)

var (
	serviceidcache = cache.NewResolvingCache("serviceidcache", time.Duration(180)*time.Second, 999)
)

func get_service_id(ctx context.Context, uid string) (*rc.Service, error) {
	fmt.Printf("[go-easyops] getting service by userid \"%s\"\n", uid)
	o, err := serviceidcache.RetrieveContext(ctx, uid, func(c context.Context, k string) (interface{}, error) {
		// we can only retrieve the serviceid from rpcaclapi/rpcinterceptor at the moment

		return rpcclient.GetServiceByUserID(c, &rc.ServiceByUserIDRequest{UserID: uid})
	})
	if err != nil {
		return nil, err
	}
	return o.(*rc.Service), nil

}
