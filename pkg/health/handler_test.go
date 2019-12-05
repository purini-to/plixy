package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestHandler(t *testing.T) {
	r := http.HandlerFunc(Handler)

	req := httptest.NewRequest("GET", "/__health__", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
