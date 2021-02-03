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
	"github.com/vinyl-linux/linux-utils/group"
	"github.com/vinyl-linux/linux-utils/passwd"
)

// useraddCmd represents the useradd command
var useraddCmd = &cobra.Command{
	Use:   "useradd [flags] username",
	Short: "Add a user to this system",
	Long:  `Add a user to this system`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		grp, err = group.Read()
		if err != nil {
			return
		}

		pwd, err = passwd.Read()
		if err != nil {
			return
		}

		return Useradd(args[0])
	},
}

func init() {
	rootCmd.AddCommand(useraddCmd)

	useraddCmd.Flags().StringVarP(&basedir, "basedir", "b", "/home/", "directory to prepend to the new username in order to create the new user homedir. Ignored when -M")
	useraddCmd.Flags().StringVarP(&comment, "comment", "c", "", "comment to add to new user. Usually a user name or long indentifier")
	useraddCmd.Flags().StringVarP(&home, "home", "d", "", "directory to be used as homedir. If set, -d is ignored")
	useraddCmd.Flags().StringVarP(&expiry, "expiry", "e", "", "if set, the day on which this account is to be disabled")
	useraddCmd.Flags().IntVarP(&gid, "gid", "g", -1, "groupid, to set as this user's primary group. Must exist. If empty, a new group is created with the same name as the requested user")
	useraddCmd.Flags().StringSliceVarP(&groups, "extra-groups", "G", []string{}, "additional groups to add user to")
	useraddCmd.Flags().StringVarP(&skel, "skel", "k", "", "a skeleton directory contains files and directories to be copied into the new homedir")
	useraddCmd.Flags().BoolVarP(&system, "system", "r", false, "create a system account (lower UID, no expiry, no homedir unless specified)")
	useraddCmd.Flags().StringVarP(&shell, "shell", "s", "/bin/sh", "login shell (can be changed later)")
	useraddCmd.Flags().BoolVarP(&noCreateHome, "no-home", "M", false, "do not create a homedir")
	useraddCmd.Flags().IntVarP(&uid, "uid", "u", -1, "uid to assign to user. The default is to use the next available")
}
