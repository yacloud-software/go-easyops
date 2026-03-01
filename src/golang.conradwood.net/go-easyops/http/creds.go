package http

import (
	"fmt"
)

type creds struct {
	username string
	password string
}

type cred_producer struct {
	host        string
	known_creds []*creds
	used        int
}

func (c *cred_producer) AddUsernamePassword(username, password string) {
	cr := &creds{username: username, password: password}
	c.known_creds = append([]*creds{cr}, c.known_creds...)
}
func (c *cred_producer) getNetRC() *creds {
	if *debug {
		fmt.Printf("Getting netrc for host \"%s\"\n", c.host)
	}
	cr := &creds{username: "foo", password: "bar"}
	return cr
}
func (c *cred_producer) GetCredentials() *creds {
	if c.used == len(c.known_creds) {
		c.used++
		//return c.getNetRC()
		return nil
	}
	if c.used > len(c.known_creds) {
		return nil
	}
	res := c.known_creds[c.used]
	c.used++
	return res
}

// idx 0 == no creds
func (c *cred_producer) GetCredentialsByIndex(idx int) *creds {
	if idx == 0 {
		return nil
	}
	if idx > len(c.known_creds) {
		return nil
	}
	return c.known_creds[idx-1]
}
func (c *cred_producer) Count() int {
	return len(c.known_creds)
}
