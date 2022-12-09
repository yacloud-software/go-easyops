package http

import (
	"context"
	"golang.yacloud.eu/apis/urlcacher"
	"io"
	"net/http"
	"time"
)

// caching http
type cHTTP struct {
	ctx context.Context
}

func (h cHTTP) Cookie(name string) *http.Cookie {
	panic("cookie not supported")
}
func (h cHTTP) Cookies() []*http.Cookie {
	panic("cookies not supported")
}
func (h cHTTP) Delete(url string, body []byte) *HTTPResponse {
	panic("delete not supported")
}
func (h cHTTP) Get(url string) *HTTPResponse {
	ctx := h.ctx
	hr := &HTTPResponse{}
	empty := &urlcacher.GetRequest{URL: url}
	srv, err := urlcacher.GetURLCacherClient().Get(ctx, empty)
	if err != nil {
		hr.err = err
		return hr
	}
	l := uint64(0)
	for {
		data, err := srv.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			hr.err = err
			break
		}
		r := uint64(len(data.Data))
		l = l + r

	}
	return hr
}
func (h cHTTP) GetStream(url string) *HTTPResponse {
	return nil
}
func (h cHTTP) Head(url string) *HTTPResponse {
	panic("head not supported")
}
func (h cHTTP) Post(url string, body []byte) *HTTPResponse {
	panic("post not supported")
}
func (h cHTTP) Put(url string, body string) *HTTPResponse {
	panic("put not supported")
}
func (h cHTTP) SetHeader(key string, value string) {
	panic("setheader not supported")
}
func (h cHTTP) SetTimeout(dur time.Duration) {
	panic("settimeout not supported")
}
