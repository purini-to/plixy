package transport

import "time"

type Option func(*transportOpt)

func WithDialTimeout(v time.Duration) Option {
	return func(t *transportOpt) {
		t.dialTimeout = v
	}
}

func WithKeepAlive(v time.Duration) Option {
	return func(t *transportOpt) {
		t.keepAlive = v
	}
}

func WithIdleConnTimeout(v time.Duration) Option {
	return func(t *transportOpt) {
		t.idleConnTimeout = v
	}
}

func WithTLSHandshakeTimeout(v time.Duration) Option {
	return func(t *transportOpt) {
		t.tlsHandshakeTimeout = v
	}
}

func WithExpectContinueTimeout(v time.Duration) Option {
	return func(t *transportOpt) {
		t.expectContinueTimeout = v
	}
}

func WithResponseHeaderTimeout(v time.Duration) Option {
	return func(t *transportOpt) {
		t.responseHeaderTimeout = v
	}
}

func WithMaxIdleConns(v int) Option {
	return func(t *transportOpt) {
		t.maxIdleConns = v
	}
}

func WithMaxIdleConnsPerHost(v int) Option {
	return func(t *transportOpt) {
		t.maxIdleConnsPerHost = v
	}
}

func WithInsecureSkipVerify(v bool) Option {
	return func(t *transportOpt) {
		t.insecureSkipVerify = v
	}
}
