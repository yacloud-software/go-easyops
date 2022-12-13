package http

import (
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
	body             []byte
	err              error
	finalurl         string
	header           map[string]string
	allheaders       []*header // responses need multiple headers of same name, e.g. "Cookie"
	received_cookies []*http.Cookie
	resp             *http.Response
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
	return h.body
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
