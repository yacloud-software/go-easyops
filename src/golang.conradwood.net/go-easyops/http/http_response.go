package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type header struct {
	Name  string
	Value string
}
type HTTPResponse struct {
	httpCode         int
	ht               *HTTP
	xbody            []byte
	body_retrieved   bool
	err              error
	finalurl         string
	header           map[string]string
	allheaders       []*header // responses need multiple headers of same name, e.g. "Cookie"
	received_cookies []*http.Cookie
	resp             *http.Response
}

// if no error and http code indicates success, return true
func (h *HTTPResponse) IsSuccess() bool {
	if h.Error() != nil {
		return false
	}
	if h.HTTPCode() < 200 {
		return false
	}
	if h.HTTPCode() >= 299 {
		return false
	}
	return true
}

func (h *HTTPResponse) HTTPCode() int {
	return h.httpCode
}
func (h *HTTPResponse) Cookies() []*http.Cookie {
	return h.received_cookies
}
func (h *HTTPResponse) AllHeaders() []*header {
	return h.allheaders
}
func (h *HTTPResponse) Body() []byte {
	if h.body_retrieved {
		return h.xbody
	}
	if h.resp == nil {
		return nil
	}
	b := &bytes.Buffer{}
	_, err := io.Copy(b, h.BodyReader())
	if err != nil {
		fmt.Printf("[go-easyops] http - bodyreader for \"%s\" failed silently (%s)\n", h.finalurl, err)
		return nil
	}
	h.setBody(b.Bytes())
	return h.xbody
}
func (h *HTTPResponse) setBody(b []byte) {
	h.xbody = b
	h.body_retrieved = true
}
func (h *HTTPResponse) Error() error {
	return h.err
}
func (h *HTTPResponse) Header(name string) string {
	name = strings.ToLower(name)
	return h.header[name]
}

// the final url (if we followed redirects the last one)
func (h *HTTPResponse) FinalURL() string {
	return h.finalurl
}
func (h *HTTPResponse) BodyReader() io.Reader {
	return h.resp.Body
}
