package rate

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/purini-to/plixy/pkg/api"

	"github.com/asaskevich/govalidator"

	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"

	"github.com/purini-to/plixy/pkg/plugin"

	"github.com/pkg/errors"
)

const (
	defaultPer          = "s"
	defaultMaxStoreSize = 65536
)

func init() {
	plugin.Register("rate", &plugin.Plugin{
		ValidateConfig: ValidateConfigFunc,
		BeforeProxy:    BeforeProxy,
	})
}

type Config struct {
	Limit        int    `json:"limit" valid:"required"`
	Burst        *int   `json:"burst"`
	Per          string `json:"per" valid:"required,in(s|m|h|d)~must be contains [s|m|h|d]"`
	MaxStoreSize int    `json:"maxStoreSize"`
}

func ValidateConfigFunc(config map[string]interface{}) error {
	return parseConfig(config, &Config{})
}

func BeforeProxy(config map[string]interface{}) (func(next http.Handler) http.Handler, error) {
	c := &Config{
		Per:          defaultPer,
		MaxStoreSize: defaultMaxStoreSize,
	}
	err := parseConfig(config, c)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("cannot parse config by rate plugin"))
	}
	if c.Burst == nil {
		b := c.Limit - 1
		c.Burst = &b
	}

	store, err := memstore.New(c.MaxStoreSize)
	if err != nil {
		return nil, errors.Wrap(err, "could not new memstore by rate plugin")
	}

	var rate throttled.Rate
	switch c.Per {
	case "s":
		rate = throttled.PerSec(c.Limit)
	case "m":
		rate = throttled.PerMin(c.Limit)
	case "h":
		rate = throttled.PerHour(c.Limit)
	case "d":
		rate = throttled.PerDay(c.Limit)
	}
	quota := throttled.RateQuota{MaxRate: rate, MaxBurst: *c.Burst}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		return nil, errors.Wrap(err, "could not new memstore by rate plugin")
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy: &throttled.VaryBy{
			Headers: []string{api.NameHeaderKey},
		},
	}

	return func(next http.Handler) http.Handler {
		return httpRateLimiter.RateLimit(next)
	}, nil
}

func parseConfig(c map[string]interface{}, v interface{}) error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "error marshal config by rate plugin")
	}

	if err = json.Unmarshal(bytes, v); err != nil {
		return errors.Wrap(err, "error unmarshal config by rate plugin")
	}

	_, err = govalidator.ValidateStruct(v)
	if err != nil {
		return err
	}

	return nil
}
