package wire

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGuard(t *testing.T) {
	var seen string
	h := Guard("k", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { seen = r.Header.Get(UsersHeader) }))
	for key, want := range map[string]string{"": "", "bad": "", "k": "x"} {
		r := httptest.NewRequest("POST", "/", nil)
		r.Header.Set(UsersHeader, "x")
		if key != "" {
			r.Header.Set(GatewayKeyHeader, key)
		}
		h.ServeHTTP(httptest.NewRecorder(), r)
		if seen != want {
			t.Fatalf("key %q: users %q", key, seen)
		}
	}
}
