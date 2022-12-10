package http

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/errors"
	"golang.yacloud.eu/apis/urlcacher"
	"io"
	"net/http"
	"time"
)

// caching http
type cHTTP struct {
	timeout    time.Duration
	ctx        context.Context
	ctx_cancel context.CancelFunc
	tich       chan bool
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
	defer h.stop_timeouter()
	ctx := h.ctx
	hr := &HTTPResponse{}
	empty := &urlcacher.GetRequest{URL: url}
	srv, err := urlcacher.GetURLCacherClient().Get(ctx, empty)
	if err != nil {
		hr.err = err
		return hr
	}

	l := uint64(0)
	started := time.Now()
	var buf []byte
	for {
		if h.timeout != 0 {
			dur := time.Since(started)
			if dur > h.timeout {
				hr.err = fmt.Errorf("timeout after %0.2fs seconds", dur.Seconds())
			}
		}

		data, err := srv.Recv()
		if (data != nil) && (len(data.Data)) > 0 {
			buf = append(buf, data.Data...)
		}
		if (data != nil) && (data.Result != nil) {
			r := data.Result
			hr.httpCode = int(r.HTTPCode)
			if !r.Success {
				if hr.httpCode == 404 {
					hr.err = errors.NotFound(ctx, "failed to retrieve url (code %d): %s", hr.httpCode, r.Message)
				} else {
					hr.err = fmt.Errorf("failed to retrieve url (code %d): %s", hr.httpCode, r.Message)
				}
				break
			}
		}

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
	hr.body = buf
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
	h.timeout = dur
	h.stop_timeouter()
	h.tich = make(chan bool)
	go h.timeouter()
}
func (h cHTTP) stop_timeouter() {
	if h.tich == nil {
		return
	}
	h.tich <- false
}
func (h cHTTP) timeouter() {
	var b bool
	select {
	case b = <-h.tich:
	case <-time.After(h.timeout):
		b = false
	}
	if b {
		h.ctx_cancel()
	}

}
