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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	dollar0   string
	useString string
)

// Command Line args
var (
	basedir      string
	comment      string
	home         string
	expiry       string
	groups       []string
	skel         string
	noCreateHome bool
	system       bool
	shell        string
	uid          int
	gid          int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short: "Vinyl Linux Utils",
	Long: `This utility provides a series of commands and tools for managing a vinyl linux install.

It is inspired, largely, by busybox and, in much the same way, provides many commands from a single binary.
The reason for this is simple: it cuts down on disk usage, while providing an interface to the same commands.

This utility is best used from any number of bash scripts which call it.
`,
	Use:  "linux-utils subdommand [args] [flags]",
	Args: cobra.ExactArgs(1),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	reset()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func reset() {
	// Set some initial, empty, default values
	basedir = "/home/"
	comment = ""
	home = ""
	expiry = ""
	groups = make([]string, 0)
	skel = ""
	noCreateHome = false
	system = false
	shell = "/bin/sh"
	uid = -255
	gid = -225
}
