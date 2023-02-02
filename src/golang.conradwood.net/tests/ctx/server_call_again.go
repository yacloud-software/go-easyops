package main

import (
	"context"
	"golang.conradwood.net/apis/common"
	ge "golang.conradwood.net/apis/getestservice"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/cmdline"
	"io"
)

func (g *geServer) CallUnaryFromStream(req *common.Void, srv ge.CtxTest_CallUnaryFromStreamServer) error {
	b := cmdline.ContextWithBuilder()
	cmdline.SetContextWithBuilder(!b)
	_, err := ge.GetCtxTestClient().TestDeSer(srv.Context(), req)
	if err != nil {
		return err
	}
	cmdline.SetContextWithBuilder(b)
	_, err = ge.GetCtxTestClient().TestDeSer(srv.Context(), req)
	if err != nil {
		return err
	}
	return nil
}

func (g *geServer) CallUnaryFromUnary(ctx context.Context, req *common.Void) (*common.Void, error) {
	b := cmdline.ContextWithBuilder()
	cmdline.SetContextWithBuilder(!b)
	_, err := ge.GetCtxTestClient().TestDeSer(ctx, req)
	if err != nil {
		return nil, err
	}
	cmdline.SetContextWithBuilder(b)
	_, err = ge.GetCtxTestClient().TestDeSer(ctx, req)
	if err != nil {
		return nil, err
	}
	return &common.Void{}, nil
}
func (g *geServer) CallStreamFromStream(req *common.Void, srv ge.CtxTest_CallStreamFromStreamServer) error {
	srv2, err := ge.GetCtxTestClient().TestStream(srv.Context(), req)
	if err != nil {
		return err
	}
	for {
		c, err := srv2.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = srv.Send(c)
		if err != nil {
			return err
		}
	}

	return nil
}
func (g *geServer) CallStreamFromUnary(ctx context.Context, req *common.Void) (*common.Void, error) {
	m := map[string]string{"provides": "default"}

	b := cmdline.ContextWithBuilder()
	cmdline.SetContextWithBuilder(!b)

	srv2, err := ge.GetCtxTestClient().TestStream(ctx, req)
	if err != nil {
		return nil, err
	}
	err = checkSrv(srv2)
	if err != nil {
		return nil, err
	}

	cmdline.SetContextWithBuilder(b)
	srv2, err = ge.GetCtxTestClient().TestStream(ctx, req)
	if err != nil {
		return nil, err
	}
	err = checkSrv(srv2)
	if err != nil {
		return nil, err
	}

	cmdline.SetContextWithBuilder(!b)
	nctx := authremote.DerivedContextWithRouting(ctx, m, true)
	srv2, err = ge.GetCtxTestClient().TestStream(nctx, req)
	if err != nil {
		return nil, err
	}
	err = checkSrv(srv2)
	if err != nil {
		return nil, err
	}

	cmdline.SetContextWithBuilder(b)
	nctx = authremote.DerivedContextWithRouting(ctx, m, true)
	srv2, err = ge.GetCtxTestClient().TestStream(nctx, req)
	if err != nil {
		return nil, err
	}
	err = checkSrv(srv2)
	if err != nil {
		return nil, err
	}

	return &common.Void{}, nil
}
