package xstoken

import (
	"testing"
	"time"
)

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
