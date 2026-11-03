package pgstate

import (
	"context"
	"os"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	dsn := os.Getenv("X360_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("X360_TEST_DATABASE_URL not set")
	}
	ctx := context.Background()
	a, err := Open(ctx, dsn, "x360_pgstate_test")
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	type doc struct{ N int }
	if err := a.Save(ctx, "k", doc{7}); err != nil {
		t.Fatal(err)
	}
	b, err := Open(ctx, dsn, "x360_pgstate_test") // a restart
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	var d doc
	ok, err := b.Load(ctx, "k", &d)
	if err != nil || !ok || d.N != 7 {
		t.Fatalf("ok=%v err=%v d=%v", ok, err, d)
	}
	if ok, _ := b.Load(ctx, "missing", &d); ok {
		t.Fatal("missing key found")
	}
}

func TestRefusesEmptyDSN(t *testing.T) {
	if _, err := Open(context.Background(), "", "x"); err != ErrNoDSN {
		t.Fatalf("err=%v", err)
	}
}

// A schema made ahead of time by an admin is used without CREATE on the database.
func TestExistingSchemaNeedsNoCreateSchema(t *testing.T) {
	dsn := os.Getenv("X360_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("X360_TEST_DATABASE_URL not set")
	}
	ctx := context.Background()
	a, err := Open(ctx, dsn, "x360_pgstate_pre")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		a.Pool().Exec(ctx, "DROP SCHEMA x360_pgstate_pre CASCADE")
		a.Close()
	}()
	b, err := Open(ctx, dsn, "x360_pgstate_pre") // schema exists: no CREATE SCHEMA issued
	if err != nil {
		t.Fatal(err)
	}
	b.Close()
}
