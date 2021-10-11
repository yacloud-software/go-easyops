package common

import (
	ed "crypto/ed25519"
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "golang.conradwood.net/apis/auth"
	"io/ioutil"
	"os"
)

var (

	// Variable "authpubkey": byte contents of file '/tmp/authkey.pub' converted on 2020-04-17 06:43:52
	authpubkey = []byte{
		0x88, 0x8A, 0x3F, 0x33, 0x1F, 0x41, 0xC7, 0x71, 0xDB, 0xC8, 0xF4,
		0xA0, 0xE8, 0x5C, 0x48, 0xF2, 0xE4, 0xE8, 0xC1, 0xF6, 0x69, 0xE3,
		0x5A, 0x23, 0x2B, 0x90, 0x38, 0xF6, 0x36, 0xF8, 0xFE, 0x38,
	}
	got_pub_key = false
)

// get the bytes from a proto that ought to be signed
func signbytes(in *pb.User) []byte {
	b := []byte(in.ID)
	ts := in.SignedAt
	x := ts
	for i := 0; i < 4; i++ {
		b = append(b, byte(x&0xFF))
		x = x << 8
	}
	return b
}

// get all the bytes from a proto that ought to be signed
func SignAllbytes(in *pb.User) []byte {
	b := []byte(in.ID)
	ts := in.SignedAt
	x := ts
	for i := 0; i < 4; i++ {
		b = append(b, byte(x&0xFF))
		x = x << 8
	}
	b = append(b, []byte(in.Email)...)
	b = append(b, []byte(in.FirstName)...)
	b = append(b, []byte(in.LastName)...)
	b = append(b, []byte(in.Abbrev)...)
	z := byte(0)
	if in.Active {
		z = z | 1<<3
	}
	if in.ServiceAccount {
		z = z | 1<<2
	}
	if in.EmailVerified {
		z = z | 1<<1
	}
	b = append(b, z)
	for _, g := range in.Groups {
		b = append(b, []byte(g.ID)...)
	}
	return b
}

/*
user has 2 signatures, one for the ID only and one "full" over all fields.
this one verifies the "full" signature (e.g. true indicates that the all fields in the user object
have been created by a 'real' auth-service and can be trusted)
*/
func VerifySignature(u *pb.User) bool {
	if u == nil {
		return false
	}
	v := ed.Verify(signPublicKey(), SignAllbytes(u), u.SignatureFull)
	return v
}

// check signature from signed user, and if valid return user. otherwise nil
func VerifySignedUser(u *pb.SignedUser) *pb.User {
	if u == nil {
		return nil
	}

	v := ed.Verify(signPublicKey(), u.User, u.Signature)
	if !v {
		return nil
	}
	res := &pb.User{}
	err := proto.Unmarshal(u.User, res)
	if err != nil {
		fmt.Printf("[go-easyops] invalid signed user (%s)\n", err)
		return nil
	}
	return res
}

func signPublicKey() ed.PublicKey {
	if got_pub_key {
		return authpubkey
	}

	fname := "/tmp/authkey.pub"
	_, err := os.Stat(fname)
	if err == nil {
		b, err := ioutil.ReadFile(fname)
		if err != nil {
			fmt.Printf("[go-easyops] failed to read pubkey file: %s\n", err)
		} else {
			fmt.Printf("[go-easyops] read temporary pubkey file: %s\n", fname)
			authpubkey = b
		}
		got_pub_key = true
		return authpubkey
	}

	return authpubkey
}
func SetPublicSigningKey(k *pb.KeyResponse) {
	got_pub_key = true
	authpubkey = k.Key
}
func GetPublicSigningKey() *pb.KeyResponse {
	return &pb.KeyResponse{Key: signPublicKey()}
}
