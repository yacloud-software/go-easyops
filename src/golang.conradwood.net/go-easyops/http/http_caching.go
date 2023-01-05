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

var (
	foo_ctr int
)

// caching http
type cHTTP struct {
	timeout               time.Duration
	ctx                   context.Context
	ctx_cancel            context.CancelFunc
	tio                   *timeouter
	last_message_received time.Time
	cancelled             bool
	debug                 bool
}
type timeouter struct {
	ch  chan bool
	idx int
}

func (h *cHTTP) SetDebug(b bool) {
	h.debug = b
}
func (h *cHTTP) Cookie(name string) *http.Cookie {
	panic("cookie not supported")
}
func (h *cHTTP) Cookies() []*http.Cookie {
	panic("cookies not supported")
}
func (h *cHTTP) Delete(url string, body []byte) *HTTPResponse {
	panic("delete not supported")
}
func (h *cHTTP) Get(url string) *HTTPResponse {
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
		if h.ctx.Err() != nil {
			hr.err = fmt.Errorf("ctx error after %0.2fs: %s", time.Since(started).Seconds(), h.ctx.Err())
			break
		}
		if h.cancelled {
			hr.err = fmt.Errorf("Cancelled (timeout)")
			break
		}
		h.last_message_received = time.Now()
		if (data != nil) && (len(data.Data)) > 0 {
			buf = append(buf, data.Data...)
		}
		if (data != nil) && (data.Result != nil) {
			r := data.Result
			hr.httpCode = int(r.HTTPCode)
			if !r.Success {
				if hr.httpCode == 404 {
					hr.err = errors.NotFound(ctx, "url \"%s\" not found (code %d): %s", url, hr.httpCode, r.Message)
				} else {
					hr.err = fmt.Errorf("failed to retrieve url \"%s\" (code %d): %s", url, hr.httpCode, r.Message)
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
func (h *cHTTP) GetStream(url string) *HTTPResponse {
	return nil
}
func (h *cHTTP) Head(url string) *HTTPResponse {
	panic("head not supported")
}
func (h *cHTTP) Post(url string, body []byte) *HTTPResponse {
	panic("post not supported")
}
func (h *cHTTP) Put(url string, body string) *HTTPResponse {
	panic("put not supported")
}
func (h *cHTTP) SetHeader(key string, value string) {
	panic("setheader not supported")
}
func (h *cHTTP) SetTimeout(dur time.Duration) {
	h.debugf("setting timeout to %0.2fs\n", dur.Seconds())
	h.timeout = dur
	h.stop_timeouter()
	foo_ctr++
	h.tio = &timeouter{ch: make(chan bool), idx: foo_ctr}
	if h.tio == nil {
		panic("no timeouter")
	}
	//	h.debugf("h=%v\n", h.tio)
	go h.timeouter(h.tio)
}
func (h *cHTTP) stop_timeouter() {
	if h.tio == nil {
		h.debugf("no timeouter to stop\n")
		return
	}
	h.debugf("stopping timeouter %d\n", h.tio.idx)
	select {
	case h.tio.ch <- false:
	case <-time.After(time.Duration(10) * time.Millisecond):
	}
}
func (h *cHTTP) SetCreds(username, password string) {
	panic("cannot use credentials for caching http")
}

func (h *cHTTP) timeouter(t *timeouter) {
	h.debugf("timeouter %d started\n", t.idx)
	var b bool
	ch := t.ch
	select {
	case b = <-ch:
		h.debugf("timeouter %d received %v\n", t.idx, b)
	case <-time.After(h.timeout):
		h.debugf("timeouter %d timer-outed\n", t.idx)
		b = true
	}
	if b {
		h.ctx_cancel()
		h.cancelled = true
		h.debugf("timeouter %d cancelled context\n", t.idx)
	}
	h.debugf("timeouter %d done\n", t.idx)

}

func (h *cHTTP) debugf(format string, args ...interface{}) {
	if !*debug {
		return
	}
	sn := "[cHTTP] "
	sx := fmt.Sprintf(format, args...)
	fmt.Print(sn + sx)
}
