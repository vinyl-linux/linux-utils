package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/vinyl-linux/linux-utils/group"
	"github.com/vinyl-linux/linux-utils/passwd"
)

func tmpFileAndCopy(in string) (fn string, err error) {
	// copy test.pwdFile to some tmpfile
	// re-open pwd pointing to it
	tmp, err := os.CreateTemp("", "")
	if err != nil {
		return
	}

	fn = tmp.Name()

	defer tmp.Close()

	orig, err := os.Open(in)
	if err != nil {
		return
	}

	defer orig.Close()

	_, err = io.Copy(tmp, orig)
	return
}

func TestUseradd(t *testing.T) {
	for _, test := range []struct {
		name        string
		args        []string
		pwdFile     string
		grpFile     string
		expectError bool
	}{
		{"happy path", []string{"useradd", "gopher"}, "testdata/simple.passwd", "testdata/simple.group", false},
		{"new user, plus extra groups", []string{"useradd", "-G", "root", "gopher"}, "testdata/simple.passwd", "testdata/simple.group", false},
		{"new user, add to existing group", []string{"useradd", "gopher", "--group", "root"}, "testdata/simple.passwd", "testdata/simple.group", false},
		{"user exists", []string{"useradd", "root"}, "testdata/simple.passwd", "testdata/simple.group", true},
		{"uid exists", []string{"useradd", "toor", "-u", "0"}, "testdata/simple.passwd", "testdata/simple.group", true},
		{"run out of system uids", []string{"useradd", "toor", "--system"}, "testdata/full-system.passwd", "testdata/simple.group", true},
		{"primary gid doesn't exist", []string{"useradd", "toor", "-g", "999"}, "testdata/simple.passwd", "testdata/simple.group", true},
		{"primary group doesn't exist", []string{"useradd", "toor", "--group", "foo"}, "testdata/simple.passwd", "testdata/simple.group", true},
		{"skel copies stuff", []string{"useradd", "gopher", "-k", "testdata/"}, "testdata/simple.passwd", "testdata/simple.group", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			var err error

			passwd.File, err = tmpFileAndCopy(test.pwdFile)
			if err != nil {
				t.Fatalf("unexpected err: %#v", err)
			}

			t.Logf("passwd: %s", passwd.File)

			group.File, err = tmpFileAndCopy(test.grpFile)
			if err != nil {
				t.Fatalf("unexpected err: %#v", err)
			}

			t.Logf("group: %s", group.File)

			pwd, err = passwd.Read()
			if err != nil {
				t.Fatalf("unexpected err: %#v", err)
			}

			grp, err = group.Read()
			if err != nil {
				t.Fatalf("unexpected err: %#v", err)
			}

			rootCmd.SetArgs(test.args)

			b := &bytes.Buffer{}
			rootCmd.SetOut(b)

			reset()
			basedir, _ = os.MkdirTemp("", "")

			err = rootCmd.Execute()
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error %#v", err)
			}
		})
	}
}
