package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/purini-to/plixy/pkg/router"
	"github.com/purini-to/plixy/pkg/trace"
)

func TestRequestID(t *testing.T) {
	t.Run("should be set anew request id If there is no id in the request header", func(t *testing.T) {
		r := router.New()
		r.Use(RequestID)

		reqID := ""
		r.GET("/", func(w http.ResponseWriter, r *http.Request) {
			reqID = trace.RequestIDFromContext(r.Context())
			assert.NotEmpty(t, reqID)
			fmt.Fprint(w, "test")
		})

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
		assert.Equal(t, reqID, rec.Header().Get(trace.RequestIDHeader))
	})

	t.Run("should be set the request header id If there is an ID in the request header", func(t *testing.T) {
		r := router.New()
		r.Use(RequestID)

		reqID := ""
		r.GET("/", func(w http.ResponseWriter, r *http.Request) {
			reqID = trace.RequestIDFromContext(r.Context())
			assert.Equal(t, "123456789", reqID)
			fmt.Fprint(w, "test")
		})

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set(trace.RequestIDHeader, "123456789")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
		assert.Equal(t, "123456789", rec.Header().Get(trace.RequestIDHeader))
	})
}
