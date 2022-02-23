//go:build linux
// +build linux

package wifi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vinyl-linux/linux-utils/netctl"
	"pifke.org/wpasupplicant"
)

var (
	networkPath = netctl.DefaultPath
)

type Networks []Network

func (ns Networks) String() string {
	lines := make([]string, 0)
	lines = append(lines, fmt.Sprintf("%s\t\t%s\t%s", "SSID", "Active", "Has Configuration"))

	for _, n := range ns {
		lines = append(lines, n.String())
	}

	return strings.Join(lines, "\n")
}

type Network struct {
	SSID       string
	Type       string
	Active     bool
	Configured bool
}

func (n Network) String() string {
	var active, configured string
	if n.Active {
		active = "*"
	}

	if n.Configured {
		configured = "*"
	}

	return fmt.Sprintf("%q\t\t%s\t\t%s", n.SSID, active, configured)
}

type Wifi struct {
	conn wpasupplicant.Conn
}

func New() (w Wifi, err error) {
	n, err := netctl.New(networkPath)
	if err != nil {
		return
	}

	foundCount := 0
	iface := ""
	for _, p := range n.Profiles {
		if p.Wifi {
			iface = p.Interface
			foundCount++
		}
	}

	switch foundCount {
	case 0:
		err = fmt.Errorf("no wireless interfaces configured")

	case 1:
		// nop - expected case

	default:
		err = fmt.Errorf("%d wireless interfaces configured, 1 expected", foundCount)
	}

	if err != nil {
		return
	}

	w.conn, err = wpasupplicant.Unixgram(iface)

	return
}

// List returns a list of the networks visible to the nic
func (w Wifi) List() (nets Networks, err error) {
	err = w.conn.Scan()
	if err != nil {
		return
	}

	conns, err := w.conn.ListNetworks()
	if err != nil {
		return
	}

	status, err := w.conn.Status()
	if err != nil {
		return
	}

	curr := status.SSID()

	scanResults, errs := w.conn.ScanResults()
	if len(errs) != 0 {
		err = flattenErrs(errs)

		return
	}

	nets = make(Networks, len(scanResults))

	for i, n := range scanResults {
		ssid := n.SSID()

		nets[i] = Network{
			SSID:       ssid,
			Active:     ssid == curr,
			Configured: contains(ssid, conns),
		}
	}

	return
}

// Create takes an SSID, pre-shared key, and creates a connection
func (w Wifi) Create(ssid, psk string) (err error) {
	id, err := w.conn.AddNetwork()
	if err != nil {
		return
	}

	err = w.conn.SetNetwork(id, "ssid", ssid)
	if err != nil {
		return
	}

	return w.conn.SetNetwork(id, "psk", psk)
}

// Connect will connect to an SSID
func (w Wifi) Connect(ssid string) (err error) {
	nets, err := w.conn.ListNetworks()
	if err != nil {
		return
	}

	var id int
	for _, n := range nets {
		if ssid == n.SSID() {
			id, err = strconv.Atoi(n.NetworkID())
			if err != nil {
				return
			}

			return w.conn.SelectNetwork(id)
		}
	}

	return fmt.Errorf("ssid %s not found", ssid)
}

// Disconnect will disconnect from all wifi networks
func (w Wifi) Disconnect() (err error) {
	status, err := w.conn.Status()
	if err != nil {
		return
	}

	if status.WPAState() != "COMPLETED" {
		return
	}

	nets, err := w.conn.ListNetworks()
	if err != nil {
		return
	}

	var id int
	for _, n := range nets {
		if n.SSID() != status.SSID() {
			continue
		}

		id, err = strconv.Atoi(n.NetworkID())
		if err != nil {
			return
		}

		return w.conn.DisableNetwork(id)
	}

	return fmt.Errorf("not connected to anything, nothing to disconnect")
}

// Save will save wifi config.
//
// It exists outside of, say, the create function because in some contexts
// we may not want to persist config. For instance: we may store keys and
// details externally (say with vault, or on an external disk).
func (w Wifi) Save() error {
	return w.conn.SaveConfig()
}

func contains(ssid string, cn []wpasupplicant.ConfiguredNetwork) bool {
	for _, n := range cn {
		if ssid == n.SSID() {
			return true
		}
	}

	return false
}

func flattenErrs(errs []error) (err error) {
	err = fmt.Errorf("error(s): ")

	for _, e := range errs {
		err = fmt.Errorf("%w %s", err, e.Error())
	}

	return
}
