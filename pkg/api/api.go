package api

import "context"

type configApiKeyType int

const (
	apiContextKey configApiKeyType = iota
)

type Definition struct {
	Apis    []*Api `yaml:"apis"`
	Version int64
}

type Api struct {
	Name  string `yaml:"name"`
	Proxy *Proxy `yaml:"proxy"`
}

type Proxy struct {
	Path      string    `yaml:"path"`
	Methods   []string  `yaml:"methods"`
	Upstream  *Upstream `yaml:"upstream"`
	FixedPath bool      `yaml:"fixedPath"`
}

type Upstream struct {
	Target string `yaml:"target"`
}

type DefinitionChanged struct {
	Definition *Definition
}

func FromContext(ctx context.Context) *Api {
	return ctx.Value(apiContextKey).(*Api)
}

func ToContext(ctx context.Context, api *Api) context.Context {
	return context.WithValue(ctx, apiContextKey, api)
}
