package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAccessLog(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	h := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "test")
	}
	r := WithLogger(logger)(AccessLog(http.HandlerFunc(h)))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())
}
