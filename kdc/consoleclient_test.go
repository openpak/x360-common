package kdc

import "testing"

func TestConsoleClient(t *testing.T) {
	c := make([]byte, CertLen)
	if ConsoleClient(c) != "" {
		t.Fatal("real console labelled")
	}
	copy(c[7:], "OPENPAKDEV")
	if ConsoleClient(c) != ClientWharf {
		t.Fatal("wharf not labelled")
	}
}
