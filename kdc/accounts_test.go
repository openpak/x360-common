package kdc

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStoreProvisionPersistsWithMode0600(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "accounts.json")
	s, err := OpenStore(path)
	if err != nil {
		t.Fatal(err)
	}
	cert := make([]byte, CertLen)
	copy(cert[2:7], []byte{1, 2, 3, 4, 5})
	a, issued, err := s.Provision(cert, t0)
	if err != nil || !issued || !a.ConsoleSignatureVerified {
		t.Fatalf("%v %v %+v", err, issued, a)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("mode %o", fi.Mode().Perm())
	}
	if _, err := os.Stat(path + ".pending"); !os.IsNotExist(err) {
		t.Fatal("pending file left behind")
	}
	// Same certificate: same account, nothing reissued.
	b, issued, err := s.Provision(cert, t0.Add(time.Minute))
	if err != nil || issued || b.Record != a.Record {
		t.Fatal("reuse failed")
	}
	// Reload from disk gives the same record and provenance.
	s2, err := OpenStore(path)
	if err != nil {
		t.Fatal(err)
	}
	name := ConsoleDigits(cert) + "@" + AccountDomain
	got := s2.Get(name)
	if got == nil || got.Record != a.Record || got.Certificate != a.Certificate || !got.ConsoleSignatureVerified || got.IssuedAt != a.IssuedAt {
		t.Fatalf("reload: %+v", got)
	}
	k1, _ := a.Key()
	k2, _ := got.Key()
	if !bytes.Equal(k1, k2) || len(k1) != 16 {
		t.Fatal("key")
	}
	// A changed certificate for the same console replaces the account with a new key.
	cert2 := append([]byte{}, cert...)
	cert2[0x1c] = 'X'
	c, issued, err := s2.Provision(cert2, t0)
	if err != nil || !issued || c.Record == a.Record || c.Certificate == a.Certificate {
		t.Fatal("reissue on certificate change failed")
	}
	// An entry without provenance is refused, not upgraded.
	s2.Get(name).ConsoleSignatureVerified = false
	if _, _, err := s2.Provision(cert2, t0); err == nil {
		t.Fatal("unverified entry reused")
	}
}

func TestStoreUnknownAndBadFile(t *testing.T) {
	dir := t.TempDir()
	s, err := OpenStore(filepath.Join(dir, "missing.json"))
	if err != nil || s.Len() != 0 || s.Get("x") != nil {
		t.Fatal(err)
	}
	bad := filepath.Join(dir, "bad.json")
	os.WriteFile(bad, []byte("{"), 0o600)
	if _, err := OpenStore(bad); err == nil {
		t.Fatal("bad JSON accepted")
	}
}

func TestTicketKeyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ticket.bin")
	k, err := LoadTicketKey(path)
	if err != nil || len(k) != 16 {
		t.Fatal(err)
	}
	fi, _ := os.Stat(path)
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("mode %o", fi.Mode().Perm())
	}
	k2, err := LoadTicketKey(path)
	if err != nil || !bytes.Equal(k, k2) {
		t.Fatal("reload")
	}
	os.WriteFile(path, make([]byte, 15), 0o600)
	if _, err := LoadTicketKey(path); err == nil {
		t.Fatal("15-byte key accepted")
	}
}

func TestEraseAccount(t *testing.T) {
	const a1, a2 = "11111111-2222-3333-4444-555555555555", "66666666-7777-8888-9999-000000000000"
	s, err := OpenStore(filepath.Join(t.TempDir(), "accounts.json"))
	if err != nil {
		t.Fatal(err)
	}
	var names [3]string
	for i := range names {
		cert := make([]byte, CertLen)
		cert[2], cert[3] = byte(i+1), 9
		if _, _, err := s.Provision(cert, t0); err != nil {
			t.Fatal(err)
		}
		names[i] = ConsoleDigits(cert) + "@" + AccountDomain
	}
	if err := s.BindCore(names[0], "l0", a1); err != nil {
		t.Fatal(err)
	}
	if err := s.BindCore(names[1], "l1", a1); err != nil {
		t.Fatal(err)
	}
	if err := s.BindCore(names[2], "l2", a2); err != nil {
		t.Fatal(err)
	}
	n, err := s.EraseAccount(context.Background(), a1)
	if err != nil || n != 2 {
		t.Fatalf("first run: %d %v", n, err)
	}
	if n, err := s.EraseAccount(context.Background(), a1); err != nil || n != 0 {
		t.Fatalf("repeat: %d %v", n, err)
	}
	if s.Get(names[0]) != nil || s.Get(names[1]) != nil || s.Get(names[2]) == nil {
		t.Fatal("wrong rows removed")
	}
	again, err := OpenStore(s.path)
	if err != nil || again.Len() != 1 {
		t.Fatalf("not persisted: %v %d", err, again.Len())
	}
	if _, err := s.EraseAccount(context.Background(), ""); err == nil {
		t.Fatal("empty id must be refused")
	}
}
