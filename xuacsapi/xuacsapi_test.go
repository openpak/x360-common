package xuacsapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusAndKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(KeyHeader) != "k" {
			w.WriteHeader(403)
			w.Write([]byte(`{"error":"forbidden"}`))
			return
		}
		w.Write([]byte(`{"online_xuid":"0009000000000007"}`))
	}))
	defer srv.Close()
	if x, err := (&Client{URL: srv.URL, Key: "k"}).OnlineXUID(context.Background(), "a"); err != nil || x != 0x0009000000000007 {
		t.Fatalf("%x %v", x, err)
	}
	_, err := (&Client{URL: srv.URL, Key: "x"}).OnlineXUID(context.Background(), "a")
	if se, ok := err.(*StatusError); !ok || se.Status != 403 || se.Msg != "forbidden" {
		t.Fatalf("%v", err)
	}
}
