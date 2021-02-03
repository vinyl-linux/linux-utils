# Vinyl Linux Utils

This repo holds linux binaries, designed for Vinyl Linux.

These are designed to be static, well tested, and with strict memory and type safety.

This project provides a single binary which is accessed with some wrapper scripts. This is similar to `busybox`, though without a single binary which reacts due to the name the binary is called with. This is to make it easier to call the binary directly/ script with.

## Contents

This binary provides:

### `useradd`

**Path:** `/bin/useradd`
**Direct Call:** `/bin/linux-utils useradd`
**Help Text:**:

```bash
$ Usage:
  linux-utils useradd [flags] username

Flags:
  -b, --basedir string         directory to prepend to the new username in order to create the new user homedir. Ignored when -M (default "/home/")
  -c, --comment string         comment to add to new user. Usually a user name or long indentifier
  -e, --expiry string          if set, the day on which this account is to be disabled
  -G, --extra-groups strings   additional groups to add user to
  -g, --gid int                groupid, to set as this user's primary group. Must exist. If empty, a new group is created with the same name as the requested user (default -1)
  -h, --help                   help for useradd
  -d, --home string            directory to be used as homedir. If set, -d is ignored
  -M, --no-home                do not create a homedir
  -s, --shell string           login shell (can be changed later) (default "/bin/sh")
  -k, --skel string            a skeleton directory contains files and directories to be copied into the new homedir
  -r, --system                 create a system account (lower UID, no expiry, no homedir unless specified)
  -u, --uid int                uid to assign to user. The default is to use the next available (default -1)
```
