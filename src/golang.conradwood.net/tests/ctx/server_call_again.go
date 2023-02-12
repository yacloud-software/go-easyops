package main

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/common"
	ge "golang.conradwood.net/apis/getestservice"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/cmdline"
	"io"
)

// return current and alternative
func cur_versions() (int, int) {
	cur := cmdline.GetContextBuilderVersion()
	alt := 0
	if cur == 0 {
		alt = CONTEXT_VERSION
	}
	return cur, alt
}

func (g *geServer) CallUnaryFromStream(req *ge.RequiredContext, srv ge.CtxTest_CallUnaryFromStreamServer) error {
	cur, alt := cur_versions()
	cmdline.SetContextBuilderVersion(alt)
	_, err := ge.GetCtxTestClient().TestDeSer(srv.Context(), req)
	if err != nil {
		return err
	}
	cmdline.SetContextBuilderVersion(cur)
	_, err = ge.GetCtxTestClient().TestDeSer(srv.Context(), req)
	if err != nil {
		return err
	}
	return nil
}

func (g *geServer) CallUnaryFromUnary(ctx context.Context, req *ge.RequiredContext) (*common.Void, error) {
	cur, alt := cur_versions()
	cmdline.SetContextBuilderVersion(alt)
	_, err := ge.GetCtxTestClient().TestDeSer(ctx, req)
	if err != nil {
		fmt.Printf("ufu: %s\n", err)
		return nil, err
	}
	cmdline.SetContextBuilderVersion(cur)
	_, err = ge.GetCtxTestClient().TestDeSer(ctx, req)
	if err != nil {
		fmt.Printf("ufu: %s\n", err)
		return nil, err
	}
	return &common.Void{}, nil
}
func (g *geServer) CallStreamFromStream(req *ge.RequiredContext, srv ge.CtxTest_CallStreamFromStreamServer) error {
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
func (g *geServer) CallStreamFromUnary(ctx context.Context, req *ge.RequiredContext) (*common.Void, error) {
	m := map[string]string{"provides": "default"}

	cur, alt := cur_versions()
	cmdline.SetContextBuilderVersion(alt)

	srv2, err := ge.GetCtxTestClient().TestStream(ctx, req)
	if err != nil {
		return nil, err
	}
	err = checkSrv(srv2)
	if err != nil {
		return nil, err
	}

	cmdline.SetContextBuilderVersion(cur)
	srv2, err = ge.GetCtxTestClient().TestStream(ctx, req)
	if err != nil {
		return nil, err
	}
	err = checkSrv(srv2)
	if err != nil {
		return nil, err
	}

	cmdline.SetContextBuilderVersion(alt)
	nctx := authremote.DerivedContextWithRouting(ctx, m, true)
	srv2, err = ge.GetCtxTestClient().TestStream(nctx, req)
	if err != nil {
		return nil, err
	}
	err = checkSrv(srv2)
	if err != nil {
		return nil, err
	}

	cmdline.SetContextBuilderVersion(cur)
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
