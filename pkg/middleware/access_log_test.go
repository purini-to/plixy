package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/purini-to/plixy/pkg/proxy"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAccessLog(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	r := proxy.New()
	r.Use(WithLogger(logger), AccessLog)

	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "test")
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())
}
