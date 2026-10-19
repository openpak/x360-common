package accountcore

import (
	"context"
	"encoding/binary"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAvatar(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/openpak.account.v1.Media/GetAvatar" || r.Header.Get("X-Internal-Key") != testKey {
			t.Error("wrong RPC or credentials")
		}
		b, _ := io.ReadAll(r.Body)
		f := parseFields(b[5:])
		if f.str(1) != "account" || f.num(2) != 256 || f.str(3) != "previous" {
			t.Errorf("wrong request: %+v", f)
		}
		var out pb
		out.str(1, string([]byte{0, 255, 0, 128}))
		out.str(2, "version")
		out = append(out, 0x18)
		out = binary.AppendUvarint(out, 1700000000)
		frame := make([]byte, 5)
		binary.BigEndian.PutUint32(frame[1:], uint32(len(out)))
		w.Header().Set("Content-Type", "application/grpc")
		w.Header().Set(http.TrailerPrefix+"Grpc-Status", "0")
		w.Write(append(frame, out...))
	}))
	defer srv.Close()
	c, err := New(srv.URL, testKey, srv.Client())
	if err != nil {
		t.Fatal(err)
	}
	a, err := c.GetAvatar(context.Background(), "account", 256, "previous")
	if err != nil {
		t.Fatal(err)
	}
	if string(a.PNG) != string([]byte{0, 255, 0, 128}) || a.Version != "version" || a.UpdatedAtUnix != 1700000000 || a.NotModified {
		t.Fatalf("bad response: %+v", a)
	}
}
