package xstoken

import (
	"testing"
	"time"
)

func TestMatchesXUID(t *testing.T) {
	for _, tc := range []struct {
		claim, path string
		want        bool
	}{
		{"0009005f5139dadb", "2533684175035099", true},
		{"0009005f5139dadb", "0009005F5139DADB", true},
		{"0009005f5139dadb", "2533684175035100", false},
		{"ffffffffffffffff", "18446744073709551615", true},
		{"ffffffffffffffff", "18446744073709551616", false},
		{"invalid", "0", false},
		{"0", "", false},
	} {
		if got := (Claims{XUID: tc.claim}).MatchesXUID(tc.path); got != tc.want {
			t.Errorf("claim %q path %q: got %v, want %v", tc.claim, tc.path, got, tc.want)
		}
	}
}

func TestSignVerify(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	s := &Signer{Key: []byte("k"), Now: func() time.Time { return now }}
	tok := s.Sign(Claims{Kind: "user", XUID: "0009005f5139dadb", Account: "a", Audience: "http://xboxlive.com", Expires: now.Add(time.Hour).Unix()})
	c, err := s.Verify(tok)
	if err != nil || c.XUID != "0009005f5139dadb" || c.Account != "a" {
		t.Fatalf("%+v %v", c, err)
	}
	if _, err := s.Verify(tok + "x"); err == nil {
		t.Fatal("tampered accepted")
	}
	s.Now = func() time.Time { return now.Add(2 * time.Hour) }
	if _, err := s.Verify(tok); err == nil {
		t.Fatal("expired accepted")
	}
	if got, ok := FromAuthorization("XBL2.0 x=uhs;" + tok); !ok || got != tok {
		t.Fatalf("uhs form: %q", got)
	}
	if got, ok := FromAuthorization("XBL2.0 x=" + tok); !ok || got != tok {
		t.Fatalf("plain form: %q", got)
	}
	if _, ok := FromAuthorization("Bearer x"); ok {
		t.Fatal("wrong scheme accepted")
	}
}
