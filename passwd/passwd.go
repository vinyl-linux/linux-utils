package passwd

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
	// File points to the passwd file
	File = "/etc/passwd"

	// SystemMin/ Max provide a range for system account UIDs
	SystemMin = 2 // ensure we don't blat away root for any reason <3
	SystemMax = 999

	// UserMin/ Max provide a range for user account UIDs
	UserMin = 1000
	UserMax = 10000

	passwdTempl = template.Must(template.New("row").Parse("{{ .Username }}:{{ .Password }}:{{ .UID }}:{{ .GUID }}:{{ .Comment }}:{{ .Home }}:{{ .LoginShell }}\n"))
)

// Passwd holds all of the entries from the passwd file
type Passwd []Entry

// Entry represents a line from /etc/passwd
type Entry struct {
	Username   string
	Password   string
	UID        int
	GID        int
	Comment    string
	Home       string
	LoginShell string
}

// Read will read passwd.File, parse each line into a set of Entry types,
// and return
func Read() (p Passwd, err error) {
	p = make(Passwd, 0)

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
		if len(components) != 7 {
			return p, fmt.Errorf("unexpected number of components in line %d: expected 7, receieved %d", lineno, len(components))
		}

		entry := Entry{
			Username:   components[0],
			Password:   components[1],
			Comment:    components[4],
			Home:       components[5],
			LoginShell: components[6],
		}

		entry.UID, err = strconv.Atoi(components[2])
		if err != nil {
			return p, err
		}

		entry.GID, err = strconv.Atoi(components[3])
		if err != nil {
			return p, err
		}

		p.Add(entry)
	}

	return
}

// Add adds the provided Entry to Passwd.
//
// NOTE: this does *not* write to disk. Call
// p.Write() when ready to write to disk
func (p *Passwd) Add(e Entry) {
	(*p) = append((*p), e)
}

// Write serialises Passwd back to /etc/passwd
func (p Passwd) Write() (err error) {
	// write to a buffer first, in case something breaks and we hose our system
	var b bytes.Buffer

	for _, entry := range p {
		err = passwdTempl.Execute(&b, entry)
		if err != nil {
			return
		}
	}

	// os.Create will truncate if file does not exist, and create if not
	passwdFile, err := os.Create(File)
	if err != nil {
		return
	}

	defer passwdFile.Close()

	_, err = io.Copy(passwdFile, &b)

	return
}

// UserExists returns true if a user by the specified name already exists
func (p Passwd) UserExists(s string) bool {
	for _, entry := range p {
		if entry.Username == s {
			return true
		}
	}

	return false
}

// UIDExists returns true if a user is already registered with this uid
func (p Passwd) UIDExists(u int) bool {
	for _, entry := range p {
		if entry.UID == u {
			return true
		}
	}

	return false
}

// NextUID returns the next available UID. If system is true, meaning we're creating a
// system account, then choose the next ID from between passwd.SystemMin and passwd.SystemMax
//
// Otherwise, choose the nextID from between passwd.UserMin and passwd.UserMax.
//
// Return an error if there are no UIDs available
func (p Passwd) NextUID(system bool) (i int, err error) {
	if system {
		i = p.nextUID(SystemMin, SystemMax)
	} else {
		i = p.nextUID(UserMin, UserMax)
	}

	if p.UIDExists(i) {
		err = fmt.Errorf("there are no more available UIDs")
	}

	return
}

func (p Passwd) nextUID(min, max int) (highest int) {
	highest = min
	for _, entry := range p {
		if entry.UID < max && entry.UID > min && entry.UID > highest {
			highest = entry.UID + 1
		}
	}

	return
}
