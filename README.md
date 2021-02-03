# Vinyl Linux Utils

This repo holds linux binaries, designed for Vinyl Linux.

These are designed to be static, well tested, and with strict memory and type safety.

## `useradd`, `groupadd`

A port of `useradd` and `groupadd` which handles all of the expected functionality, such as:

1. Parsing `/etc/passwd`, `/etc/shadow`, `/etc/group`
1. CRUD operations on above

`groupadd` is expected to be symlinked to `useradd` (or vice versa, or both symlinked elsewhere- it doesn't strictly matter)
