package passwd

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
		expect      Passwd
		expectError bool
	}{
		{"missing file", "testdata/missing", Passwd{}, true},
		{"single line", "testdata/passwd", Passwd{Entry{Username: "root", Password: "x", UID: 0, GID: 0, Comment: "root", Home: "/root", LoginShell: "/bin/ash"}}, false},
		{"mangled entry", "testdata/mangled", Passwd{}, true},
		{"bad uid", "testdata/baduid", Passwd{}, true},
		{"bad gid", "testdata/badgid", Passwd{}, true},
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

func TestPasswd_UserExists(t *testing.T) {
	pwd := Passwd{Entry{Username: "root", Password: "x", UID: 0, GID: 0, Comment: "root", Home: "/root", LoginShell: "/bin/ash"}}

	for _, test := range []struct {
		user   string
		expect bool
	}{
		{"username", false},
		{"root", true},
	} {
		t.Run(test.user, func(t *testing.T) {
			got := pwd.UserExists(test.user)
			if test.expect != got {
				t.Errorf("expected %v, received %v", test.expect, got)
			}
		})
	}
}

func TestPasswd_UIDExists(t *testing.T) {
	pwd := Passwd{Entry{Username: "root", Password: "x", UID: 0, GID: 0, Comment: "root", Home: "/root", LoginShell: "/bin/ash"}}

	for _, test := range []struct {
		uid    int
		expect bool
	}{
		{0, true},
		{1, false},
	} {
		t.Run("", func(t *testing.T) {
			got := pwd.UIDExists(test.uid)
			if test.expect != got {
				t.Errorf("expected %v, received %v", test.expect, got)
			}
		})
	}
}

func TestPasswd_NextUID(t *testing.T) {
	pwd := Passwd{Entry{Username: "root", Password: "x", UID: 0, GID: 0, Comment: "root", Home: "/root", LoginShell: "/bin/ash"}}

	for _, test := range []struct {
		system bool
		expect int
	}{
		{true, 2},
		{false, 1000},
	} {
		t.Run("", func(t *testing.T) {
			got, _ := pwd.NextUID(test.system)

			if test.expect != got {
				t.Errorf("expected %v, received %v", test.expect, got)
			}
		})
	}

}
