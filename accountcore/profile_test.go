package accountcore

import "testing"

func TestPatchClearIsPresentField(t *testing.T) {
	empty, bio := "", "hi"
	f := parseFields(profileRequest("a", ProfilePatch{Bio: &bio, Motto: &empty}))
	if _, ok := f.s[5]; !ok || f.str(4) != "hi" || f.str(1) != "a" {
		t.Fatalf("fields %v", f.s)
	}
	if _, ok := f.s[2]; ok {
		t.Fatal("omitted name must stay absent")
	}
}
