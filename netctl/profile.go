// +build linux

package netctl

import (
	"context"
	"fmt"
	"net"
	"regexp"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/nclient4"
	"github.com/vishvananda/netlink"
)

var (
	// Test whether an interface is loopback by testing whether the name
	// is suggests loopback by convention.
	//
	// This is a potentially weak way of doing this; realistically a loopback
	// device can be named anything, and can only reliably be tested by checking
	// for the IFF_LOOPBACK flag.
	//
	// However, it's probably fine /shrug
	loopback = regexp.MustCompile(`^lo[0-9]*$`)
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

	// Loopback devices are special; we can go ahead and set them
	// up the same way each time. In fact, the loopback file only needs
	// the value of `Interface` to be set
	if loopback.Match([]byte(p.Interface)) {
		return p.UpLoopback()
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

// UpLoopback brings a loopback device up
func (p Profile) UpLoopback() (err error) {
	return p.BringUp()
}

// Down will bring an interface down
func (p Profile) Down() (err error) {
	return wrap("iface Down", handle.LinkSetDown(p.link))
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
	var nameserver net.IP

	switch idx {
	case 0:
		a.AddressParsed, a.GatewayParsed, nameserver, a.NetmaskParsed, err = p.negotiateIPV4()

	default:
		err = fmt.Errorf("unknown index #%d", idx)
	}

	if err != nil {
		return
	}

	return writeResolv(nameserver)
}

func (p Profile) negotiateIPV4() (address, gateway, nameserver net.IP, netmask net.IPMask, err error) {
	if p.dclient4 == nil {
		if Verbose {
			p.dclient4, err = nclient4.New(p.Interface, nclient4.WithDebugLogger())
		} else {
			p.dclient4, err = nclient4.New(p.Interface)
		}
		if err != nil {
			err = wrap("nclient4.New", err)

			return
		}
	}

	offer, err := p.dclient4.DiscoverOffer(context.Background())
	if err != nil {
		err = wrap("dhcpv4 negotiation", err)
		return
	}

	address = offer.YourIPAddr
	netmask = net.IPMask(net.IP(offer.Options.Get(dhcpv4.OptionSubnetMask)))
	gateway = net.IP(offer.Options.Get(dhcpv4.OptionRouter))
	nameserver = net.IP(offer.Options.Get(dhcpv4.OptionDomainNameServer))

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

	return wrap("AddAddress", handle.AddrAdd(p.link, ipConfig))
}

// SetGateway uses netlink to set the gateway of an interface
func (p Profile) SetGateway(a Address) (err error) {
	route := &netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: p.link.Attrs().Index,
		Gw:        a.GatewayParsed,
	}

	return wrap("SetGateway", handle.RouteAdd(route))
}

// BringUp uses netlink to bring an interface up
func (p Profile) BringUp() (err error) {
	return wrap("iface Up", handle.LinkSetUp(p.link))
}
