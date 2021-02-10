// +build linux

/*
Copyright © 2021 James Condron <james@zero-internet.org.uk>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vinyl-linux/linux-utils/netctl"
)

// upCmd represents the up command
var netctl_upCmd = &cobra.Command{
	Use:   "up [iface | all]",
	Short: "Start stored network profiles",
	Long: `Start stored network profiles.

This command takes either an interface (eth0, wlan1, etc.) or the word "all" which brings all interfaces up
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		n, err := netctl.New(netctlDir)
		if err != nil {
			return
		}

		return netctlUpDo(args[0], n)
	},
}

func netctlUpDo(iface string, n netctl.Netctl) error {
	netctl.Verbose = verbose

	if iface == "all" {
		return netctlUpDoAll(n)
	}

	return netctlUpIface(iface, n)
}

func netctlUpDoAll(n netctl.Netctl) (err error) {
	for _, i := range n.Profiles {
		err = i.Up()
		if err != nil {
			return
		}
	}

	return
}

func netctlUpIface(iface string, n netctl.Netctl) (err error) {
	p, err := n.Profile(iface)
	if err != nil {
		return
	}

	return p.Up()
}

func init() {
	netctlCmd.AddCommand(netctl_upCmd)

	netctl_upCmd.Flags().StringVarP(&netctlDir, "basedir", "b", netctl.DefaultPath, "location of network config files")
	netctl_upCmd.Flags().BoolVarP(&verbose, "", "v", false, "log dhcp, network calls")
}
