package kdc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/openpak/x360-common/pgstate"
)

// pgSchema returns a private schema for the test and drops it afterwards.
// Tests that need Postgres skip without X360_TEST_DATABASE_URL.
func pgSchema(t *testing.T) (dsn, schema string) {
	t.Helper()
	dsn = os.Getenv("X360_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("X360_TEST_DATABASE_URL not set")
	}
	h := sha256.Sum256([]byte(t.Name() + t.TempDir()))
	schema = "t_kdc_" + hex.EncodeToString(h[:6])
	t.Cleanup(func() {
		if p, err := pgstate.Open(context.Background(), dsn, schema); err == nil {
			p.Pool().Exec(context.Background(), "DROP SCHEMA "+schema+" CASCADE")
			p.Close()
		}
	})
	return dsn, schema
}

func TestOpenPGStoreRefusesWithoutDSN(t *testing.T) {
	if _, err := OpenPGStore(context.Background(), "", Schema, ""); err == nil {
		t.Fatal("opened without a database URL")
	}
}

// A restart (close and reopen) gives the same account; a second store on the
// same schema (the gateway) sees what the first one issued after Reload.
func TestPGStoreRestartAndReload(t *testing.T) {
	dsn, schema := pgSchema(t)
	ctx := context.Background()
	s, err := OpenPGStore(ctx, dsn, schema, "")
	if err != nil {
		t.Fatal(err)
	}
	reader, err := OpenPGStore(ctx, dsn, schema, "")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	cert := make([]byte, CertLen)
	copy(cert[2:7], []byte{1, 2, 3, 4, 5})
	a, issued, err := s.Provision(cert, t0)
	if err != nil || !issued {
		t.Fatalf("%v %v", err, issued)
	}
	name := ConsoleDigits(cert) + "@" + AccountDomain
	if reader.Get(name) != nil {
		t.Fatal("reader saw the account before Reload")
	}
	if err := reader.Reload(ctx); err != nil || reader.Get(name) == nil || reader.Get(name).Record != a.Record {
		t.Fatalf("reload: %v", err)
	}
	s.Close()
	s2, err := OpenPGStore(ctx, dsn, schema, "")
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Close()
	b, issued, err := s2.Provision(cert, t0.Add(time.Minute))
	if err != nil || issued || b.Record != a.Record || !b.ConsoleSignatureVerified {
		t.Fatalf("restart lost the account: %v issued=%v", err, issued)
	}
}

// The legacy accounts.json is imported once, left in place, and never
// overrides what the database holds afterwards.
func TestPGStoreImportsLegacyFileOnce(t *testing.T) {
	dsn, schema := pgSchema(t)
	ctx := context.Background()
	file := filepath.Join(t.TempDir(), "accounts.json")
	fs, err := OpenStore(file)
	if err != nil {
		t.Fatal(err)
	}
	cert := make([]byte, CertLen)
	copy(cert[2:7], []byte{9, 8, 7, 6, 5})
	a, _, err := fs.Provision(cert, t0)
	if err != nil {
		t.Fatal(err)
	}
	s, err := OpenPGStore(ctx, dsn, schema, file)
	if err != nil {
		t.Fatal(err)
	}
	name := ConsoleDigits(cert) + "@" + AccountDomain
	if g := s.Get(name); g == nil || g.Record != a.Record || g.Certificate != a.Certificate {
		t.Fatal("record not imported")
	}
	s.Close()
	if _, err := os.Stat(file); err != nil {
		t.Fatal("legacy file removed")
	}
	// A different file now must not replace the database content.
	other := filepath.Join(t.TempDir(), "other.json")
	os.WriteFile(other, []byte("{}"), 0o600)
	s2, err := OpenPGStore(ctx, dsn, schema, other)
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Close()
	if s2.Len() != 1 || s2.Get(name) == nil {
		t.Fatal("database content replaced by the file")
	}
}
