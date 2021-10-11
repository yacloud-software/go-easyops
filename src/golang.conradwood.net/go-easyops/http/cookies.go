package http

import (
	"fmt"
	"net/http"
	"net/url"
)

type Cookies struct {
	cookies []*http.Cookie
}

func (c *Cookies) SetCookies(u *url.URL, cookies []*http.Cookie) {
	c.cookies = append(c.cookies, cookies...)
}
func (c *Cookies) Cookies(u *url.URL) []*http.Cookie {
	return c.cookies
}
func (c *Cookies) Print() {
	for _, ck := range c.cookies {
		fmt.Printf("Cookie: %s\n", ck.Name)
	}
}
