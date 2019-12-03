package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"

	"github.com/purini-to/plixy/pkg/trace"
)

func TestRequestID(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("should be set anew request id If there is no id in the request header", func(t *testing.T) {
		reqID := ""

		h := func(w http.ResponseWriter, r *http.Request) {
			reqID = trace.RequestIDFromContext(r.Context())
			assert.NotEmpty(t, reqID)
			_, _ = fmt.Fprint(w, "test")
		}
		r := WithLogger(logger)(RequestID(http.HandlerFunc(h)))

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
		assert.Equal(t, reqID, rec.Header().Get(trace.RequestIDHeader))
	})

	t.Run("should be set the request header id If there is an ID in the request header", func(t *testing.T) {
		reqID := ""

		h := func(w http.ResponseWriter, r *http.Request) {
			reqID = trace.RequestIDFromContext(r.Context())
			assert.Equal(t, "123456789", reqID)
			_, _ = fmt.Fprint(w, "test")
		}
		r := WithLogger(logger)(RequestID(http.HandlerFunc(h)))

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set(trace.RequestIDHeader, "123456789")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
		assert.Equal(t, "123456789", rec.Header().Get(trace.RequestIDHeader))
	})
}
