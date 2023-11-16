package topology

import (
	"net/http"
	"testing"
)

func TestHttpForm(t *testing.T) {
	request, err := http.NewRequest("GET", "http://example.org?name=wenhaozhou", nil)
	if err != nil {
		return
	}
	t.Log(request.FormValue("name"))
}
