package utils

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"math/rand"
	"os"
	"time"
)

var (
	randsrc = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// return random integer between 0 and n
func RandomInt(max int64) int {
	t := randsrc.Int63n(max)
	return int(t)
}

func RandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits

	)

	b := make([]byte, n)
	// A randsrc.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randsrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randsrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// never returns - if error != nil, print it and exit
func Bail(txt string, err error) {
	if err == nil {
		return
	}
	fmt.Printf("%s: %s\n", txt, ErrorString(err))
	os.Exit(10)
}

// return true if string has letters only
func IsLettersOnly(txt string) bool {
	return IsOnlyChars(txt, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// true only if string "txt" is made up exclusively of characters in "valid"
func IsOnlyChars(txt string, valid string) bool {
	for _, x := range txt {
		b := false
		for _, y := range valid {
			if x == y {
				b = true
			}
		}
		if !b {
			return false
		}
	}
	return true
}

// stall for a random amount of "upto" minutes
func RandomStall(minutes int) {
	if minutes == 0 {
		return
	}
	t := randsrc.Int63n(int64(minutes * 60))
	time.Sleep(time.Duration(t) * time.Second)
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func MinInt32(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func MaxInt32(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func MinInt64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func MaxInt64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func PrettyNumber(number uint64) string {
	return humanize.Bytes(number)
}
