package accountcore

import (
	"context"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBirthdayRequiresMatchingActiveAccount(t *testing.T) {
	for _, tc := range []struct {
		id      string
		status  byte
		wantErr bool
	}{{"owner", byte(StatusActive), false}, {"other", byte(StatusActive), true}, {"owner", 0, true}} {
		t.Run(tc.id+string(rune(tc.status)), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var out pb
				out.str(1, tc.id)
				out = append(out, 0x10, tc.status)
				out.str(6, "2000-10-20")
				frame := make([]byte, 5)
				binary.BigEndian.PutUint32(frame[1:], uint32(len(out)))
				w.Header().Set("Content-Type", "application/grpc")
				w.Header().Set(http.TrailerPrefix+"Grpc-Status", "0")
				w.Write(append(frame, out...))
			}))
			defer srv.Close()
			c, _ := New(srv.URL, testKey, srv.Client())
			got, err := c.Birthday(context.Background(), "owner")
			if (err != nil) != tc.wantErr || (!tc.wantErr && got != "2000-10-20") {
				t.Fatalf("birthday=%q err=%v", got, err)
			}
		})
	}
}
