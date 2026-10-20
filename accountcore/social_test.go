package accountcore

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSocialRPCs(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		req := parseFields(raw[5:])
		if r.Header.Get("X-Internal-Key") != testKey || req.str(1) != "owner" {
			t.Error("missing principal or internal authentication")
		}
		var out pb
		switch r.URL.Path {
		case "/openpak.account.v1.Social/ListFriends", "/openpak.account.v1.Social/ListBlocks":
			for _, id := range []string{"first", "second"} {
				var friend pb
				friend.str(1, id)
				friend.str(2, "Name "+id)
				out.str(1, string(friend))
			}
		case "/openpak.account.v1.Social/CheckInteraction":
			if req.str(2) != "other" || req.num(3) != 1 {
				t.Error("wrong interaction")
			}
			out = append(out, 0x08, 0x01)
		default:
			t.Errorf("unexpected method %s", r.URL.Path)
		}
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
	for _, list := range []func(context.Context, string) ([]Friend, error){c.ListFriends, c.ListBlocks} {
		friends, err := list(context.Background(), "owner")
		if err != nil || len(friends) != 2 || friends[0].AccountID != "first" || friends[1].AccountID != "second" || friends[1].DisplayName != "Name second" {
			t.Fatalf("repeated projections lost: %+v %v", friends, err)
		}
	}
	allowed, err := c.CanDiscover(context.Background(), "owner", "other")
	if err != nil || !allowed {
		t.Fatalf("discovery: %v %v", allowed, err)
	}
}

func TestSocialRejectsTruncatedLists(t *testing.T) {
	for _, b := range [][]byte{{0x0a, 3, 'a'}, {0x0a, 0, 0x0a, 2, 'a'}, {0}, {0x08, 0x80}} {
		if _, err := socialMessages(b); err == nil {
			t.Errorf("accepted truncated protobuf %x", b)
		}
	}
}

func TestSocialPropagatesCoreFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/grpc")
		w.Header().Set(http.TrailerPrefix+"Grpc-Status", "14")
		w.Write([]byte{0, 0, 0, 0, 0})
	}))
	defer srv.Close()
	c, _ := New(srv.URL, testKey, srv.Client())
	_, err := c.ListFriends(context.Background(), "owner")
	var rpcErr *Error
	if !errors.As(err, &rpcErr) || rpcErr.Code != 14 {
		t.Fatalf("failure hidden: %v", err)
	}
}
