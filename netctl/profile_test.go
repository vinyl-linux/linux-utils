// +build linux

package netctl

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/vishvananda/netlink"
)

type testNetLinkHandle struct {
	err bool
}

func (h testNetLinkHandle) LinkSetDown(_ netlink.Link) error {
	if h.err {
		return fmt.Errorf("an error")
	}

	return nil
}

func (h testNetLinkHandle) LinkSetUp(_ netlink.Link) error {
	if h.err {
		return fmt.Errorf("an error")
	}

	return nil
}

func (h testNetLinkHandle) AddrAdd(_ netlink.Link, _ *netlink.Addr) error {
	if h.err {
		return fmt.Errorf("an error")
	}

	return nil
}

func (h testNetLinkHandle) RouteAdd(_ *netlink.Route) error {
	if h.err {
		return fmt.Errorf("an error")
	}

	return nil
}

func (h testNetLinkHandle) LinkByName(_ string) (l netlink.Link, err error) {
	l = &netlink.Dummy{netlink.NewLinkAttrs()}

	if h.err {
		err = fmt.Errorf("an error")
	}

	return
}

type testDHCPClient struct {
	err bool
}

func (c testDHCPClient) DiscoverOffer(_ context.Context, _ ...dhcpv4.Modifier) (o *dhcpv4.DHCPv4, err error) {
	if c.err {
		err = fmt.Errorf("an error")
	}

	o = &dhcpv4.DHCPv4{
		YourIPAddr: net.ParseIP("192.168.1.3"),
		Options: dhcpv4.Options{
			uint8(dhcpv4.OptionSubnetMask): []byte{0xff, 0xff, 0xff, 0x0},
			uint8(dhcpv4.OptionRouter):     []byte{0xc0, 0xa8, 0x1, 0x1},
		},
	}

	return
}

func TestAddress_Parse(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Fatalf("unexpected panic %#v", err)
		}
	}()

	a := Address{
		Address: "192.168.1.2",
		Netmask: "255.255.255.0",
		Gateway: "192.168.1.1",
	}

	a.Parse()
}

func TestUp(t *testing.T) {
	nonDHCP := Profile{
		Interface: "test0",
		IPv4: Address{
			Address: "192.168.1.2",
			Netmask: "255.255.255.0",
			Gateway: "192.168.1.1",
			Enable:  true,
			DHCP:    false,
		},
		IPv6: Address{
			Address: "fd12:3456:789a:1::8",
			Netmask: "ffffffffffffffff0000000000000000",
			Gateway: "fd12:3456:789a:1::1",
			Enable:  true,
			DHCP:    false,
		},
	}

	ip4DHCP := Profile{
		Interface: "test0",
		IPv4: Address{
			Enable: true,
			DHCP:   true,
		},
		dclient4: testDHCPClient{},
	}

	ip4DHCPErr := Profile{
		Interface: "test0",
		IPv4: Address{
			Enable: true,
			DHCP:   true,
		},
		dclient4: testDHCPClient{err: true},
	}

	ip4DHCPClientErr := Profile{
		Interface: "test0",
		IPv4: Address{
			Enable: true,
			DHCP:   true,
		},
	}

	for _, test := range []struct {
		name        string
		profile     Profile
		handle      netlinkHandle
		expectError bool
	}{
		{"happy path", nonDHCP, testNetLinkHandle{}, false},
		{"netlink errors", nonDHCP, testNetLinkHandle{err: true}, true},
		{"with dhcp", ip4DHCP, testNetLinkHandle{}, false},
		{"with dhcp errors", ip4DHCPErr, testNetLinkHandle{}, true},
		{"dhcp client errors", ip4DHCPClientErr, testNetLinkHandle{}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			handle = test.handle

			err := test.profile.Up()
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %#v", err)
			}
		})
	}
}
