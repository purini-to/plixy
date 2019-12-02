package config

import "context"

type configApiKeyType int

const (
	apiContextKey configApiKeyType = iota
)

type ApiDefinition struct {
	Apis []*Api `json:"apis"`
}

type Api struct {
	Name  string `json:"name"`
	Proxy *Proxy `json:"proxy"`
}

type Proxy struct {
	Path     string    `json:"path"`
	Methods  []string  `json:"methods"`
	Upstream *Upstream `json:"upstream"`
}

type Upstream struct {
	Target string `json:"target"`
}

func ApiFromContext(ctx context.Context) *Api {
	return ctx.Value(apiContextKey).(*Api)
}

func ApiToContext(ctx context.Context, api *Api) context.Context {
	return context.WithValue(ctx, apiContextKey, api)
}
