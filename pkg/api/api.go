package api

import (
	"context"
	"regexp"
)

type configApiKeyType int

const (
	apiContextKey configApiKeyType = iota
	varsContextKey
)

var varsReg = regexp.MustCompile(`\{(.+?)\}`)

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
	Target string   `yaml:"target"`
	Vars   []string `yaml:"-"`
}

func (u *Upstream) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var m map[string]interface{}
	if err := unmarshal(&m); err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "target":
			v := v.(string)
			u.Target = v
			group := varsReg.FindAllSubmatch([]byte(v), -1)
			if group == nil {
				continue
			}
			for _, g := range group {
				u.Vars = append(u.Vars, string(g[1]))
			}
		default:
			continue
		}
	}
	return nil
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

func VarsFromContext(ctx context.Context) map[string]string {
	return ctx.Value(varsContextKey).(map[string]string)
}

func VarsToContext(ctx context.Context, api map[string]string) context.Context {
	return context.WithValue(ctx, varsContextKey, api)
}
