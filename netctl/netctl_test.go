// +build linux

package netctl

import (
	"io/ioutil"
	"net"
	"testing"
)

func TestNewDefaults(t *testing.T) {
	_, err := NewDefaults()
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestNew(t *testing.T) {
	for _, test := range []struct {
		path        string
		expectError bool
	}{
		{"testdata/happy-path", false},
		{"testdata/mangled", true},
	} {
		t.Run(test.path, func(t *testing.T) {
			_, err := New(test.path)
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %#v", err)
			}
		})
	}
}

func TestNetctl_Profile(t *testing.T) {
	n := Netctl{
		Profiles: []Profile{
			{Interface: "test0"},
			{Interface: "test1"},
		},
	}

	_, err := n.Profile("test0")
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}

	_, err = n.Profile("eth0")
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestWriteResolve(t *testing.T) {
	origResolv := ResolvFile
	defer func() {
		ResolvFile = origResolv
	}()

	for _, test := range []struct {
		name        string
		fn          string
		expect      string
		expectError bool
	}{
		{"happy path", "/tmp/resolv.conf", "nameserver 1.1.1.1\n", false},
		{"no such dir", "/1/2/3/4/5/resovl.conf", "", true},
		{"no write permission", "/no-access", "", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			ResolvFile = test.fn

			ip := net.ParseIP("1.1.1.1")
			err := writeResolv(ip)
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %#v", err)
			}

			got, _ := ioutil.ReadFile(test.fn)
			gots := string(got)

			if test.expect != gots {
				t.Errorf("expected %q, received %q", test.expect, gots)
			}
		})
	}
}
