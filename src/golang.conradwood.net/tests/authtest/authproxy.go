package main

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/getestservice"
	git "golang.conradwood.net/apis/gitserver"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/tokens"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

func StartAuthProxy() {
	//p := utils.RandomInt(100)
	//	go create.NewEasyOpsTest(&testserver{}, *port+p)
	time.Sleep(time.Duration(500) * time.Millisecond) // give server some time to register...
	fmt.Printf("Creating easyops client...\n")
	_, err := pb.GetEasyOpsTestClient().CheckSerialisation(authremote.Context(), &pb.Count{})
	utils.Bail("failed", err)

}
func (*testserver) CheckSerialisation(ctx context.Context, req *pb.Count) (*common.Void, error) {
	fmt.Printf("Running serialisation test...(Count=%d)\n", req.Count)
	fmt.Printf("Inbound User: %s\n", auth.CurrentUserString(ctx))
	if req.Count == 0 {
		tokens.SetServiceTokenParameter(servicetokens[0])
		return pb.GetEasyOpsTestClient().CheckSerialisation(ctx, &pb.Count{Count: req.Count + 1})
	}
	foo, err := auth.SerialiseContext(ctx)
	utils.Bail("Failed to serialise context", err)
	fmt.Printf("Serialised context: %s\n", string(foo))
	new_ctx, err := auth.RecreateContextWithTimeout(time.Duration(10)*time.Second, foo)
	utils.Bail("Failed to fix context", err)
	fmt.Printf("Restored User: %s\n", auth.CurrentUserString(new_ctx))
	// now try a gitserver call...
	r, err := git.GetGIT2Client().RepoByID(new_ctx, &git.ByIDRequest{ID: 2})
	utils.Bail("Failed to call git", err)
	fmt.Printf("Git Repo: %s\n", r.ArtefactName)
	return &common.Void{}, nil
}
