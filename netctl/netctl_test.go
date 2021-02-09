// +build linux

package netctl

import (
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
