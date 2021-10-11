package rpc

import (
	"context"
	"flag"
	"fmt"
	apb "golang.conradwood.net/apis/auth"
	"golang.conradwood.net/apis/goeasyops"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
)

var (
	easyopsClient        goeasyops.EasyOpsClient
	easyopsClientAddress = flag.String("ge_service", "localhost:5002", "go easyops service ip:port")
	myServiceUser        *apb.User
)

const (
	CTXKEY = "go_easyops_ctx_key"
)

func initclient() {
	if easyopsClient != nil {
		return
	}

	// NOT loadbalanced on purpose. Each "cluster", "server", "workstation" has their own instance
	conn, err := grpc.Dial(
		*easyopsClientAddress,
		//	grpc.WithBlock(),
		grpc.WithTransportCredentials(client.GetClientCreds()),
	)
	utils.Bail("Misconfigured ip stack. Unable to connect to easyopsService", err)
	easyopsClient = goeasyops.NewEasyOpsClient(conn)
	ctx := context.Background()
	ubl := goeasyops.UserByTokenRequest{Token: tokens.GetServiceTokenParameter()}
	ar, err := easyopsClient.UserByToken(ctx, &ubl)
	if err != nil {
		utils.Bail("Failed to get serviceaccount: \n", err) // should not happen
	}
	if !ar.Valid {
		fmt.Printf("Invalid service token.\n")
	}
	myServiceUser = ar.User
}

func PrintContext(ctx context.Context) {
	fmt.Println(ContextToString(ctx))
}

func ContextToString(ctx context.Context) string {
	co := fromContext(ctx)
	return co.PrettyString()
}
func ContextWithLogin(username string, password string) (context.Context, error) {
	initclient()
	ctx := context.Background()
	ubl := goeasyops.UserByLoginRequest{Username: username, Password: password}
	ar, err := easyopsClient.UserByLogin(ctx, &ubl)
	if err != nil {
		return nil, err
	}
	if !ar.Valid {
		return nil, errors.AccessDenied(ctx, "no such user: %s ", username)
	}
	ctx = NewContextWithUserAndService(ctx, ar.User, myServiceUser)
	return ctx, nil
}

// adds this services' user & service accounts to context and returns it
func NewContextWithUserAndService(ctx context.Context, user, service *apb.User) context.Context {
	co := &contextObject{
		user:    user,
		service: service,
	}

	return co.NewContext(ctx)
}
