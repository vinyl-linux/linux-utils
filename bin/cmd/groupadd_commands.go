package cmd

import (
	"fmt"

	"github.com/vinyl-linux/linux-utils/group"
)

var (
	grp group.Group
)

func Groupadd(name string) (err error) {
	if gid > -1 {
		if grp.GIDExists(gid) {
			return fmt.Errorf("gid %d is already in use", gid)
		}
	} else {
		gid, err = grp.NextGID()
		if err != nil {
			return
		}
	}

	if grp.GroupExists(name) {
		return fmt.Errorf("group %s already exists", name)
	}

	grp.Add(group.Entry{
		Name:     name,
		Password: "x",
		GID:      gid,
	})

	return grp.Write()
}

// AddToGroup adds a user to a group. It assumes users are valid
// mainly because it doesn't really matter if they're not
func AddToGroup(gid int, user string) (err error) {
	err = grp.AddUser(gid, user)
	if err != nil {
		return
	}

	return grp.Write()
}

// AddToGroups adds a user to a group by name.
//
// It also assumes a user is valid. It only writes to
// disk when all groups are validated
func AddToGroups(names []string, user string) (err error) {
	if len(names) == 0 {
		return
	}

	var g group.Entry
	for _, name := range names {
		g, err = grp.ByName(name)
		if err != nil {
			return
		}

		err = grp.AddUser(g.GID, user)
		if err != nil {
			return
		}
	}

	return grp.Write()
}
