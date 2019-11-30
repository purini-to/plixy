package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestRouter_Router(t *testing.T) {
	r := New()
	r.Use(
		func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("1-Middleware", "1")
				next.ServeHTTP(w, r)
			}
			return http.HandlerFunc(fn)
		},
		func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("2-Middleware", "2")
				next.ServeHTTP(w, r)
			}
			return http.HandlerFunc(fn)
		},
	)
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "test")
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())
	assert.Equal(t, "1", rec.Header().Get("1-Middleware"))
	assert.Equal(t, "2", rec.Header().Get("2-Middleware"))
}
