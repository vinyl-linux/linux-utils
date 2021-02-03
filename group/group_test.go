package group

import (
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	origFile := File
	defer func() {
		File = origFile
	}()

	for _, test := range []struct {
		name        string
		input       string
		expect      Group
		expectError bool
	}{
		{"missing file", "testdata/missing", Group{}, true},
		{"single line", "testdata/group", Group{Entry{Name: "sys", Password: "x", GID: 3, Users: []string{"root", "bin", "adm"}}}, false},
		{"mangled entry", "testdata/mangled", Group{}, true},
		{"bad gid", "testdata/badgid", Group{}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			File = test.input

			p, err := Read()
			t.Logf("%#v", err)

			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %#v", err)
			}

			if !reflect.DeepEqual(test.expect, p) {
				t.Errorf("expected %#v, received %#v", test.expect, p)
			}
		})
	}
}

func TestGroup_GroupExists(t *testing.T) {
	grp := Group{Entry{Name: "sys", Password: "x", GID: 3, Users: []string{"root", "bin", "adm"}}}

	for _, test := range []struct {
		user   string
		expect bool
	}{
		{"games", false},
		{"sys", true},
	} {
		t.Run(test.user, func(t *testing.T) {
			got := grp.GroupExists(test.user)
			if test.expect != got {
				t.Errorf("expected %v, received %v", test.expect, got)
			}
		})
	}
}

func TestGroup_UIDExists(t *testing.T) {
	grp := Group{Entry{Name: "sys", Password: "x", GID: 3, Users: []string{"root", "bin", "adm"}}}

	for _, test := range []struct {
		uid    int
		expect bool
	}{
		{3, true},
		{1, false},
	} {
		t.Run("", func(t *testing.T) {
			got := grp.GIDExists(test.uid)
			if test.expect != got {
				t.Errorf("expected %v, received %v", test.expect, got)
			}
		})
	}
}

func TestGroup_NextUID(t *testing.T) {
	grp := Group{Entry{Name: "sys", Password: "x", GID: 3, Users: []string{"root", "bin", "adm"}}}

	for _, test := range []struct {
		system bool
		expect int
	}{
		{true, 4},
	} {
		t.Run("", func(t *testing.T) {
			got, _ := grp.NextGID()

			if test.expect != got {
				t.Errorf("expected %v, received %v", test.expect, got)
			}
		})
	}
}
