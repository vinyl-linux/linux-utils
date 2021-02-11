// +build linux

package netctl

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/pelletier/go-toml"
	"github.com/vishvananda/netlink"
)

// netlink client
type netlinkHandle interface {
	LinkSetDown(netlink.Link) error
	LinkSetUp(netlink.Link) error
	AddrAdd(netlink.Link, *netlink.Addr) error
	RouteAdd(*netlink.Route) error
	LinkByName(string) (netlink.Link, error)
}

// dhcp4 dhcp client
type dhcpOfferer4 interface {
	DiscoverOffer(context.Context, ...dhcpv4.Modifier) (*dhcpv4.DHCPv4, error)
}

var (
	handle netlinkHandle
)

const DefaultPath = "/etc/vinyl/network.d"

var (
	// Location of dns resolver config
	ResolvFile = "/etc/resolv.conf"

	// Log calls, operations
	Verbose = false
)

// Netctl provides access to the files at /etc/vinyl/network
// which govern network connections on vinyl systems
type Netctl struct {
	Profiles []Profile
}

// NewDefaults calls New with all of the default values
func NewDefaults() (Netctl, error) {
	return New(DefaultPath)
}

// New returns a new Netctl
func New(p string) (n Netctl, err error) {
	handle, err = netlink.NewHandle()
	if err != nil {
		return
	}

	n = Netctl{}

	err = n.parse(p)

	return
}

// Profile returns the configured profile for this interface
func (n Netctl) Profile(iface string) (p Profile, err error) {
	for _, p = range n.Profiles {
		if p.Interface == iface {
			return
		}
	}

	return Profile{}, fmt.Errorf("interface %s is not configured", iface)
}

// parse reads all config files in n.path and updates n.profiles
// containing them all
func (n *Netctl) parse(p string) error {
	n.Profiles = make([]Profile, 0)

	return filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".toml") {
			p, err := readProfile(path)
			if err != nil {
				return err
			}

			n.Profiles = append(n.Profiles, p)
		}

		return nil
	})
}

// readProfile will... read a profile
func readProfile(filename string) (p Profile, err error) {
	p = NewProfile()

	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = toml.Unmarshal(d, &p)

	return
}

func wrap(s string, err error) error {
	if err == nil {
		return err
	}

	return fmt.Errorf("%s: %w", s, err)
}

func writeResolv(ns net.IP) (err error) {
	r, err := os.Create(ResolvFile)
	if err != nil {
		return
	}

	_, err = r.WriteString(fmt.Sprintf("nameserver %s\n", ns.String()))

	return
}
