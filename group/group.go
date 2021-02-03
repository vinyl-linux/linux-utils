package group

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/template"
)

var (
	// File points to the group file
	File = "/etc/group"

	groupTempl = template.Must(template.New("row").Parse("{{ .Name }}:{{ .Password }}:{{ .GID }}:{{ .UsersCSV }}\n"))
)

// Group holds all of the entries from the group file
type Group []Entry

// Entry represents a line from /etc/group
type Entry struct {
	Name     string
	Password string
	GID      int
	Users    []string
	UsersCSV string
}

// Read will read group.File, parse each line into a set of Entry types,
// and return
func Read() (p Group, err error) {
	p = make(Group, 0)

	var data *os.File
	data, err = os.Open(File)
	if err != nil {
		return
	}

	defer data.Close()

	buf := bufio.NewReader(data)
	lineno := 0
	for {
		lineno++

		line, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return p, err
		}

		if err == io.EOF {
			break
		}

		components := strings.Split(strings.TrimSpace(line), ":")
		if len(components) != 4 {
			return p, fmt.Errorf("unexpected number of components in line %d: expected 4, receieved %d", lineno, len(components))
		}

		entry := Entry{
			Name:     components[0],
			Password: components[1],
			Users:    strings.Split(components[3], ","),
		}

		entry.GID, err = strconv.Atoi(components[2])
		if err != nil {
			return p, err
		}

		p.Add(entry)
	}

	return
}

// Add adds the provided Entry to Group.
//
// NOTE: this does *not* write to disk. Call
// p.Write() when ready to write to disk
func (p *Group) Add(e Entry) {
	(*p) = append((*p), e)
}

// Write serialises Group back to /etc/group
func (p Group) Write() (err error) {
	// write to a buffer first, in case something breaks and we hose our system
	var b bytes.Buffer

	for _, entry := range p {
		err = groupTempl.Execute(&b, entry)
		if err != nil {
			return
		}
	}

	// os.Create will truncate if file does not exist, and create if not
	groupFile, err := os.Create(File)
	if err != nil {
		return
	}

	defer groupFile.Close()

	_, err = io.Copy(groupFile, &b)

	return
}

// GroupExists returns true if a group by the specified name already exists
func (p Group) GroupExists(s string) bool {
	for _, entry := range p {
		if entry.Name == s {
			return true
		}
	}

	return false
}

// GIDExists returns true if a user is already registered with this gid
func (p Group) GIDExists(u int) bool {
	for _, entry := range p {
		if entry.GID == u {
			return true
		}
	}

	return false
}

// NextGID returns the next available GID.
func (p Group) NextGID() (i int, err error) {
	i = -1
	for _, entry := range p {
		if entry.GID > i {
			i = entry.GID + 1
		}
	}

	return
}
