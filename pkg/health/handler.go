package health

import (
	"fmt"
	"net/http"

	"github.com/purini-to/plixy/pkg/config"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "health: OK, version: %s", config.Version)
}
