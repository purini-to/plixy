package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRealIP(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "test")
	}
	r := RealIP(http.HandlerFunc(h))

	t.Run("should set the header first IP to RemoteAddr, if there is X-Forwarded-For in the header", func(tt *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Forwarded-For", "203.0.113.195, 70.41.3.18, 150.172.238.178")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
		assert.Equal(t, "203.0.113.195", req.RemoteAddr)
	})

	t.Run("should set the header IP to RemoteAddr, if there is X-Real-IP in the header", func(tt *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Real-IP", "203.0.113.195")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
		assert.Equal(t, "203.0.113.195", req.RemoteAddr)
	})
}
