// +build linux

package netctl

import (
	"context"
	"fmt"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/nclient4"
	"github.com/vishvananda/netlink"
)

// Address holds the configuration necessary for setting
// an interface's IP address, including both manual and
// dhcp settings
type Address struct {
	Address       string `toml:"address,omitempty"`
	AddressParsed net.IP `toml:"-"`

	Netmask       string     `toml:"netmask,omitempty"`
	NetmaskParsed net.IPMask `toml:"-"`

	Gateway       string `toml:"gateway,omitempty"`
	GatewayParsed net.IP `toml:"-"`

	// Enable needs to be true in order to configure
	Enable bool `toml:"enable,omitempty"`

	// DHCP true means Address/ Netmask/ Gateway can be ignored
	// (if set, they're ignored anyway)
	DHCP bool `toml:"dhcp,omitempty"`
}

// Parse will turn various user specified components of an Address
// into net.IP/net.IPMask types
func (a *Address) Parse() {
	if a.Address != "" {
		a.AddressParsed = net.ParseIP(a.Address)
	}

	if a.Netmask != "" {
		a.NetmaskParsed = net.IPMask(net.ParseIP(a.Netmask))
	}

	if a.Gateway != "" {
		a.GatewayParsed = net.ParseIP(a.Gateway)
	}

	return
}

// Profile contains configuration and data for an interface
type Profile struct {
	Interface string  `toml:"interface"`
	IPv4      Address `toml:",omitempty"`
	IPv6      Address `toml:",omitempty"`

	// link points to the underlying netlink object
	link netlink.Link

	// dhcp clients
	dclient4 dhcpOfferer4
}

// NewProfile returns a Profile with initial defaults set
func NewProfile() Profile {
	return Profile{}
}

// Up will bring an interface up
func (p Profile) Up() (err error) {
	p.link, err = handle.LinkByName(p.Interface)
	if err != nil {
		return
	}

	err = p.BringUp()
	if err != nil {
		return
	}

	for idx, addr := range []Address{
		p.IPv4,
		p.IPv6,
	} {
		if !addr.Enable {
			continue
		}

		addr.Parse()

		if addr.DHCP {
			err = p.PopulateFromDHCP(idx, &addr)
			if err != nil {
				return
			}
		}

		for _, f := range []func(Address) error{
			p.SetAddress,
			p.SetGateway,
		} {
			err = f(addr)
			if err != nil {
				return
			}
		}
	}

	return
}

// Down will bring an interface down
func (p Profile) Down() (err error) {
	return handle.LinkSetDown(p.link)
}

// PopulateFromDHCP accepts an index (where 0 is IPv4, and 1 is IPv6)
// and a pointer to an Address. It will:
//
// 1. Bring up the interface
// 2. Request either an IPv4 address or an IPv6 address (from index)
// 3. Update address with this data
// 4. Take the interface back down (to allow netlink to control the iface)
// 5. Return
func (p *Profile) PopulateFromDHCP(idx int, a *Address) (err error) {
	switch idx {
	case 0:
		a.AddressParsed, a.GatewayParsed, a.NetmaskParsed, err = p.negotiateIPV4()

	default:
		err = fmt.Errorf("unknown index #%d", idx)
	}

	return
}

func (p Profile) negotiateIPV4() (address, gateway net.IP, netmask net.IPMask, err error) {
	if p.dclient4 == nil {
		p.dclient4, err = nclient4.New(p.Interface, nclient4.WithDebugLogger())
		if err != nil {
			return
		}
	}

	offer, err := p.dclient4.DiscoverOffer(context.Background())
	if err != nil {
		return
	}

	address = offer.YourIPAddr
	netmask = net.IPMask(net.IP(offer.Options.Get(dhcpv4.OptionSubnetMask)))
	gateway = net.IP(offer.Options.Get(dhcpv4.OptionRouter))

	return
}

// SetAddress uses netlink to set the address of an interface
func (p Profile) SetAddress(a Address) (err error) {
	ipConfig := &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   a.AddressParsed,
			Mask: a.NetmaskParsed,
		},
	}

	return handle.AddrAdd(p.link, ipConfig)
}

// SetGateway uses netlink to set the gateway of an interface
func (p Profile) SetGateway(a Address) (err error) {
	route := &netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: p.link.Attrs().Index,
		Dst:       &net.IPNet{IP: a.GatewayParsed, Mask: net.CIDRMask(32, 32)},
	}

	return handle.RouteAdd(route)
}

// BringUp uses netlink to bring an interface up
func (p Profile) BringUp() (err error) {
	return handle.LinkSetUp(p.link)
}
