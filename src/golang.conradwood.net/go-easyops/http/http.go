package http

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/prometheus"
	"google.golang.org/grpc/codes"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	durationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "goeasyops_httpclient_duration",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			Help:       "V=1 unit=s DESC=execution time of successful http calls",
		}, []string{"name"},
	)
	callcounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goeasyops_httpclient_total_calls",
			Help: "V=1 unit=ops DESC=total number of outbound http calls",
		}, []string{"name"},
	)
	failcounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goeasyops_httpclient_failures",
			Help: "V=1 unit=ops DESC=number of failed outbound http calls",
		}, []string{"name"},
	)

	tr = &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:          50,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       3 * time.Second,
		ResponseHeaderTimeout: 3 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
	}
	mytr  = &transport{t: tr}
	debug = flag.Bool("ge_debug_http", false, "debug http code")
)

func init() {
	prometheus.MustRegister(durationSummary, callcounter, failcounter)
}

type HTTP struct {
	MetricName string // if not "", will export metrics for this call
	username   string
	password   string
	err        error
	headers    map[string]string
	jar        *Cookies
	transport  *transport // nil for default
}

func (h *HTTP) SetTimeout(dur time.Duration) {
	if h.transport == nil {
		h.transport = &transport{t: &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:          50,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       dur,
			ResponseHeaderTimeout: dur,
			ExpectContinueTimeout: dur,
			DialContext: (&net.Dialer{
				Timeout:   dur,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
		},
		}
	} else {
		h.transport.t.IdleConnTimeout = dur
		h.transport.t.ResponseHeaderTimeout = dur
		h.transport.t.ExpectContinueTimeout = dur
		h.transport.t.DialContext = (&net.Dialer{
			Timeout:   dur,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext
	}
}
func (h *HTTP) promLabels() prometheus.Labels {
	return prometheus.Labels{"name": h.MetricName}
}
func (h *HTTP) doMetric() bool {
	return h.MetricName != ""
}
func WithAuth(username string, password string) *HTTP {
	res := &HTTP{username: username, password: password}
	if username == "" {
		res.err = fmt.Errorf("Missing username")
	}
	return res
}

type header struct {
	Name  string
	Value string
}
type HTTPResponse struct {
	httpCode   int
	ht         *HTTP
	body       []byte
	err        error
	finalurl   string
	header     map[string]string
	allheaders []*header // responses need multiple headers of same name, e.g. "Cookie"
}

func (h *HTTPResponse) HTTPCode() int {
	return h.httpCode
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
func (h *HTTP) SetHeader(key string, value string) {
	if h.headers == nil {
		h.headers = make(map[string]string)
	}
	h.headers[key] = value
}

func (h *HTTP) Head(url string) *HTTPResponse {
	hr := &HTTPResponse{ht: h}
	if h.err != nil {
		hr.err = h.err
		return hr
	}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		hr.err = err
		return hr
	}
	hr.do(req)
	return hr
}
func (h *HTTP) Get(url string) *HTTPResponse {
	hr := &HTTPResponse{ht: h}
	if h.err != nil {
		hr.err = h.err
		return hr
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		hr.err = err
		return hr
	}
	hr.do(req)
	return hr
}
func (h *HTTP) Delete(url string, body []byte) *HTTPResponse {
	hr := &HTTPResponse{ht: h}
	if h.err != nil {
		hr.err = h.err
		return hr
	}
	b := strings.NewReader(string(body))
	req, err := http.NewRequest("DELETE", url, b)
	if err != nil {
		hr.err = err
		return hr
	}
	hr.do(req)
	return hr
}
func (h *HTTP) Post(url string, body []byte) *HTTPResponse {
	hr := &HTTPResponse{ht: h}
	if h.err != nil {
		hr.err = h.err
		return hr
	}
	if *debug {
		fmt.Printf("Body: \"%s\"\n", string(body))
	}
	b := strings.NewReader(string(body))
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		hr.err = err
		return hr
	}
	hr.do(req)
	return hr
}
func (h *HTTP) Put(url string, body string) *HTTPResponse {
	hr := &HTTPResponse{ht: h}
	if h.err != nil {
		hr.err = h.err
		return hr
	}
	b := strings.NewReader(body)
	req, err := http.NewRequest("PUT", url, b)
	if err != nil {
		hr.err = err
		return hr
	}
	hr.do(req)
	return hr
}

/************************** direct calls ****************************/
func Get(url string) ([]byte, error) {
	h := &HTTP{}
	res := h.Get(url)
	return res.Body(), res.Error()
}
func Post(url string, body []byte) ([]byte, error) {
	h := &HTTP{}
	res := h.Post(url, body)
	return res.Body(), res.Error()
}
func Put(url string, body string) ([]byte, error) {
	h := &HTTP{}
	res := h.Put(url, body)
	return res.Body(), res.Error()
}
func (h *HTTPResponse) Header(name string) string {
	name = strings.ToLower(name)
	return h.header[name]
}

// the final url (if we followed redirects the last one)
func (h *HTTPResponse) FinalURL() string {
	return h.finalurl
}

func (h *HTTP) Cookies() []*http.Cookie {
	if h.jar == nil {
		return nil
	}
	return h.jar.cookies
}
func (h *HTTP) Cookie(name string) *http.Cookie {
	if h.jar == nil {
		return nil
	}
	for _, c := range h.jar.cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (h *HTTPResponse) do(req *http.Request) *HTTPResponse {
	if h.ht.jar == nil {
		h.ht.jar = &Cookies{}
	}

	ctx := context.Background()
	if h.ht.username != "" {
		req.SetBasicAuth(h.ht.username, h.ht.password)
	}
	tr := mytr
	if h.ht.transport != nil {
		tr = h.ht.transport
	}
	hclient := &http.Client{Transport: tr, Jar: h.ht.jar}
	h.ht.jar.Print()
	if h.ht.headers != nil {
		for k, v := range h.ht.headers {
			if *debug {
				fmt.Printf("Header \"%s\" = \"%s\"\n", k, v)
			}
			req.Header.Set(k, v)

			if strings.ToLower(k) == "host" {
				req.Host = v
			}

		}
	}

	for _, c := range h.ht.jar.cookies {
		fmt.Printf("Adding cookie %s\n", c.Name)
		req.Header.Add("Cookie", fmt.Sprintf("%s=%s", c.Name, c.Value))
	}
	started := time.Now()
	if h.ht.doMetric() {
		callcounter.With(h.ht.promLabels()).Inc()
	}
	resp, err := hclient.Do(req)
	if resp != nil {
		h.httpCode = resp.StatusCode
		h.finalurl = resp.Request.URL.String()

	}
	if err != nil {
		if h.ht.doMetric() {
			failcounter.With(h.ht.promLabels()).Inc()
		}
		h.err = err
		return h
	}

	if *debug {
		fmt.Printf("Request to %s complete (code=%d)\n", h.FinalURL(), h.HTTPCode())
	}
	h.header = make(map[string]string)
	for k, va := range resp.Header {
		if len(va) == 0 {
			continue
		}
		k = strings.ToLower(k)
		for _, v := range va {
			h.allheaders = append(h.allheaders, &header{Name: k, Value: v})
		}
		h.header[k] = va[0]
	}

	defer resp.Body.Close()
	pbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.err = err
		return h
	}
	if resp.StatusCode == 404 {
		if h.ht.doMetric() {
			failcounter.With(h.ht.promLabels()).Inc()
		}
		h.err = errors.Error(ctx, codes.NotFound, "not found", "%s not found", req.URL)
	} else if resp.StatusCode > 299 || resp.StatusCode < 200 {
		if h.ht.doMetric() {
			failcounter.With(h.ht.promLabels()).Inc()
		}
		h.err = fmt.Errorf("Http to \"%s\" failed with code %d", req.URL, resp.StatusCode)
	}
	if h.err == nil {
		durationSummary.With(h.ht.promLabels()).Observe(time.Since(started).Seconds())
	}
	h.body = pbody
	return h
}

type transport struct {
	t *http.Transport
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to keep track
// of the current request.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if *debug {
		//		fmt.Printf("Request: %#v\n", req)
		//		fmt.Printf("Body: \"%v\"\n", req.Body)
		fmt.Printf("URL: %s\n", req.URL)
	}
	return t.t.RoundTrip(req)
}
