package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
	"github.com/vinyl-linux/linux-utils/passwd"
)

var (
	pwd passwd.Passwd
)

func Useradd(user string) (err error) {
	// does user exist?
	if pwd.UserExists(user) {
		return fmt.Errorf("user %s already exists", user)
	}

	// is uid in use?
	if uid != -1 {
		if pwd.UIDExists(uid) {
			return fmt.Errorf("uid %d is already in use", uid)
		}
	}

	uid, err = pwd.NextUID(system)
	if err != nil {
		return
	}

	// is g set, and does it exist?
	if gid > -1 {
		if !grp.GIDExists(gid) {
			return fmt.Errorf("gid %d doesn't correspond to an actual group", gid)
		}
	} else {
		// if not, create
		// can we have the same gid as uid?
		if !grp.GIDExists(uid) {
			gid = uid
		}

		err = Groupadd(user)
		if err != nil {
			return
		}
	}

	// add user to primary group
	//
	// Note: gid gets set in Groupadd... globals feel fragile, but
	// hopefull this comment makes this less magic
	err = AddToGroup(gid, user)
	if err != nil {
		return
	}

	// add user to supplementary groups
	err = AddToGroups(groups, user)
	if err != nil {
		return
	}

	if home == "" {
		home = filepath.Join(basedir, user)
	}

	entry := passwd.Entry{
		Username:   user,
		Password:   "x",
		UID:        uid,
		GID:        gid,
		Comment:    comment,
		Home:       home,
		LoginShell: shell,
	}

	pwd.Add(entry)
	defer func() {
		err1 := pwd.Write()
		if err != nil && err1 != nil {
			err = fmt.Errorf("errors occurred: %e, %e", err, err1)
		}

		if err1 != nil {
			err = err1
		}
	}()

	if !system {
		// if dir exists, break!
		_, err = os.Stat(home)
		if os.IsNotExist(err) {
			err = nil
		}

		if err != nil {
			return
		}

		// create dir
		err = os.MkdirAll(home, 0700)
		if err != nil {
			return
		}

		if skel != "" {
			err = copy.Copy(skel, home)
		}
	}

	return
}
