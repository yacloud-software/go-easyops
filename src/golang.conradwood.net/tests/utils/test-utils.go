package main

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/echoservice"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/objectstore"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"time"
)

func main() {
	check_store()
	t := time.Now()
	fmt.Printf("Time now: %s\n", utils.TimeString(t))
	t, err := utils.LocalTime(context.Background())
	utils.Bail("failed to get timezone", err)
	fmt.Printf("Time local: %s\n", utils.TimeString(t))
	print(14)
	print(65)
	print(125)
	print(120)
	print(0)
	print(60*60 + 3 + 24)
	print(60*60*2 + 60*5 + 40)
	secs := uint32(0)
	fmt.Printf("'not set' as Age: %s\n", utils.TimestampAgeString(secs))
	s := fmt.Sprintf("I am a test hexdump buffer with some text and stuff\n")
	fmt.Println(utils.HexdumpWithLen(16, "hex ", []byte(s)))
}

func print(age int) {
	secs := uint32(time.Now().Unix()) - uint32(age)
	fmt.Printf("%d seconds as Age: %s\n", age, utils.TimestampAgeString(secs))
}
func check_store() {
	ctx := authremote.Context()
	pr := &echoservice.PingRequest{
		SequenceNumber: uint32(utils.RandomInt(500)),
		Payload:        utils.RandomString(100),
	}
	utils.Bail("failed to store proto", objectstore.StoreProto(ctx, "goeasyops-test", pr))
	npr := &echoservice.PingRequest{}
	utils.Bail("failed to retrieve proto", objectstore.RetrieveProto(ctx, "goeasyops-test", npr))
	if pr.SequenceNumber != npr.SequenceNumber {
		fmt.Printf("ERROR: Store/Retrieve proto: Sequence number mismatch\n (%d != %d)", pr.SequenceNumber, npr.SequenceNumber)
		os.Exit(10)
	}
	if pr.Payload != npr.Payload {
		fmt.Printf("ERROR: Store/Retrieve proto: payload mismatch\n (%s != %s)", pr.Payload, npr.Payload)
		os.Exit(10)
	}

}
