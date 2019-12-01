package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/purini-to/plixy/pkg/proxy"
	"github.com/purini-to/plixy/pkg/trace"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRecover(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Internal Server Error should be returned if the panic argument is not an error type", func(t *testing.T) {
		r := proxy.New()
		r.Use(WithLogger(logger), Recover, func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				panic("panic")
			}

			return http.HandlerFunc(fn)
		})

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError)+"\n", rec.Body.String())
	})

	t.Run("error.Error() should be returned if the panic argument is an error type", func(t *testing.T) {
		r := proxy.New()
		r.Use(Recover, func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Errorf("error panic"))
			}

			return http.HandlerFunc(fn)
		})

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set(trace.RequestIDHeader, "123456789")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "error panic\n", rec.Body.String())
	})
}
