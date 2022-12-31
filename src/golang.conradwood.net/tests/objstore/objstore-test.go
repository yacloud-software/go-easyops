package main

import (
	"fmt"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/client"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

func main() {
	fmt.Printf("objstore test\n")
	key := "FOOKEYOBJSTORETEST"
	blobs := [][]byte{
		make([]byte, 10000),
		make([]byte, 20000),
		make([]byte, 100000),
		make([]byte, 1000000),
		make([]byte, 10000000),
		make([]byte, 100000000),
	}
	for _, blob := range blobs {
		ctx := authremote.Context()
		exp := time.Now().Add(time.Duration(10) * time.Minute)
		fmt.Printf("Using blob with %d bytes\n", len(blob))
		err := client.PutWithIDAndExpiry(ctx, key, blob, exp)
		utils.Bail("failed to put", err)
		err = client.EvictNoResult(ctx, key)
		utils.Bail("failed to evict", err)
	}
}
