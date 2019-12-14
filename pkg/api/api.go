package api

import (
	"context"
	"regexp"

	"github.com/asaskevich/govalidator"
)

type configApiKeyType int

const (
	apiContextKey configApiKeyType = iota
	varsContextKey
)

var varsReg = regexp.MustCompile(`\{(.+?)\}`)

type Definition struct {
	Apis    []*Api `yaml:"apis" valid:"required"`
	Version int64  `yaml:"-"`
}

func (d *Definition) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

type Api struct {
	Name  string `yaml:"name" valid:"required"`
	Proxy *Proxy `yaml:"proxy" valid:"required"`
}

type Proxy struct {
	Path     string    `yaml:"path" valid:"required,matches(^/)~path must be start with '/'"`
	Methods  []string  `yaml:"methods" valid:"matches(^(GET|HEAD|POST|PUT|PATCH|DELETE|CONNECT|OPTIONS|TRACE)$)~methods must be http methods. [GET|HEAD|POST|PUT|PATCH|DELETE|CONNECT|OPTIONS|TRACE]"`
	Upstream *Upstream `yaml:"upstream" valid:"required"`
}

type Upstream struct {
	Target    string   `yaml:"target" valid:"required,requrl~target must be url"`
	FixedPath bool     `yaml:"fixedPath"`
	Vars      []string `yaml:"-"`
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
