package transport

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Transport struct {
	*http.Transport
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.Transport.RoundTrip(req)
}

type transportOpt struct {
	dialTimeout           time.Duration
	keepAlive             time.Duration
	idleConnTimeout       time.Duration
	tlsHandshakeTimeout   time.Duration
	expectContinueTimeout time.Duration
	responseHeaderTimeout time.Duration
	maxIdleConns          int
	maxIdleConnsPerHost   int
	insecureSkipVerify    bool
}

func (t transportOpt) hash() string {
	return fmt.Sprintf("%v", t)
	//return strings.Join([]string{
	//	fmt.Sprintf("dialTimeout:%v", t.dialTimeout),
	//	fmt.Sprintf("keepAlive:%v", t.keepAlive),
	//	fmt.Sprintf("idleConnTimeout:%v", t.idleConnTimeout),
	//	fmt.Sprintf("tlsHandshakeTimeout:%v", t.tlsHandshakeTimeout),
	//	fmt.Sprintf("expectContinueTimeout:%v", t.expectContinueTimeout),
	//	fmt.Sprintf("responseHeaderTimeout:%v", t.responseHeaderTimeout),
	//	fmt.Sprintf("maxIdleConns:%v", t.maxIdleConns),
	//	fmt.Sprintf("maxIdleConnsPerHost:%v", t.maxIdleConnsPerHost),
	//	fmt.Sprintf("insecureSkipVerify:%v", t.insecureSkipVerify),
	//}, ";")
}

func New(opts ...Option) *http.Transport {
	// default
	t := transportOpt{
		dialTimeout:           3 * time.Second,
		keepAlive:             120 * time.Second,
		idleConnTimeout:       120 * time.Second,
		tlsHandshakeTimeout:   3 * time.Second,
		expectContinueTimeout: 1 * time.Second,
		responseHeaderTimeout: 5 * time.Second,
		maxIdleConns:          512,
		maxIdleConnsPerHost:   128,
		insecureSkipVerify:    false,
	}

	for _, opt := range opts {
		opt(&t)
	}

	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   t.dialTimeout,
			KeepAlive: t.keepAlive,
		}).DialContext,
		TLSHandshakeTimeout:   t.tlsHandshakeTimeout,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: t.insecureSkipVerify},
		MaxIdleConns:          t.maxIdleConns,
		MaxIdleConnsPerHost:   t.maxIdleConnsPerHost,
		IdleConnTimeout:       t.idleConnTimeout,
		ResponseHeaderTimeout: t.responseHeaderTimeout,
		ExpectContinueTimeout: t.expectContinueTimeout,
		ForceAttemptHTTP2:     true,
	}

	return tr
}
