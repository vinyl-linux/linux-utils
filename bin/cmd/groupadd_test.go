package cmd

import (
	"bytes"
	"testing"

	"github.com/vinyl-linux/linux-utils/group"
)

func TestGroupadd(t *testing.T) {
	for _, test := range []struct {
		name        string
		args        []string
		grpFile     string
		expectError bool
	}{
		{"happy path", []string{"groupadd", "gopher"}, "testdata/simple.group", false},
		{"group exists", []string{"groupadd", "root"}, "testdata/simple.group", true},
		{"gid exists", []string{"groupadd", "toor", "-g", "0"}, "testdata/simple.group", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var err error

			group.File, err = tmpFileAndCopy(test.grpFile)
			if err != nil {
				t.Fatalf("unexpected err: %#v", err)
			}

			t.Logf("group: %s", group.File)

			grp, err = group.Read()
			if err != nil {
				t.Fatalf("unexpected err: %#v", err)
			}

			rootCmd.SetArgs(test.args)

			b := &bytes.Buffer{}
			rootCmd.SetOut(b)

			reset()

			err = rootCmd.Execute()
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error %#v", err)
			}
		})
	}

}
