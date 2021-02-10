// +build linux

package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestNetctl_Up(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	for _, test := range []struct {
		name        string
		args        []string
		expectError bool
	}{
		{"all profiles, missing (default) dir", []string{"netctl", "up", "all"}, true},
		{"all profiles, empty list", []string{"netctl", "up", "all", "-b", dir}, false},
		{"all profiles, with errors", []string{"netctl", "up", "all", "-b", "testdata/netctl"}, true},
		{"single iface, with errors", []string{"netctl", "up", "test0", "-b", "testdata/netctl"}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			reset()

			rootCmd.SetArgs(test.args)

			b := &bytes.Buffer{}
			rootCmd.SetOut(b)

			err := rootCmd.Execute()
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error %#v", err)
			}
		})
	}
}
