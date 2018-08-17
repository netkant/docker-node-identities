package types

// ...
import (
	"fmt"
	"os/user"
	"strings"
)

// Group is parsed from "name:password:gid:members"
type Group struct {
	Name     string
	Password string
	GID      string
	Members  string // omitted in create
}

// NewGroup ...
func NewGroup(data string) (Group, error) {
	var t Group
	parts := strings.Split(data, ":")

	// ...
	if len(parts) > 4 {
		return t, fmt.Errorf("group string should only consist of 4 parts")
	}

	// ...
	t = Group{
		Name:     getIndexValue(parts, 0, ""),
		Password: getIndexValue(parts, 1, ""),
		GID:      getIndexValue(parts, 2, ""),
		Members:  getIndexValue(parts, 3, ""),
	}

	return t, nil
}

// Create ...
func (t Group) Create() error {
	// ...
	err := fmt.Errorf("could not create group, no supported command was found")

	// ...
	if t.Name == "" {
		return fmt.Errorf("name is required")
	}

	// ...
	if _, err := user.LookupGroup(t.Name); err == nil {
		return fmt.Errorf("group '%s' already exists", t.Name)
	}

	// ...
	if t.GID != "" {
		if _, err := user.LookupGroupId(t.GID); err == nil {
			return fmt.Errorf("group with id '%s' already exists", t.GID)
		}
	}

	// ...
	if path := commandExists("groupadd"); path != "" {
		err = t.groupadd(path)
	}

	return err
}

// Delete ...
func (t Group) Delete() error {
	// ...
	err := fmt.Errorf("could not create group, no supported command was found")

	// ...
	if path := commandExists("groupdel"); path != "" {
		err = t.groupdel(path)
	}

	return err
}

// groupadd ...
func (t Group) groupadd(path string) error {
	args := []string{}

	// ...
	if t.Password != "" {
		args = append(args, "--password")
		args = append(args, t.Password)
	}

	// ...
	if t.GID != "" {
		args = append(args, "--gid")
		args = append(args, t.GID)
	}

	// the name to add
	args = append(args, t.Name)

	// ...
	_, exitCode := commandExecute(path, args...)
	switch exitCode {
	case 0:
		return nil
	default:
		return fmt.Errorf("unknown groupadd error: %d", exitCode)
	}
}

// ...
func (t Group) groupdel(path string) error {
	args := []string{}

	// ...
	if t.Name == "" {
		return fmt.Errorf("name is required")
	}

	// the name to delete
	args = append(args, t.Name)

	// ...
	_, exitCode := commandExecute(path, args...)
	switch exitCode {
	case 6:
		return fmt.Errorf("group '%s' does not exist", t.Name)
	case 0:
		// clog.info("group created")
		return nil
	default:
		return fmt.Errorf("unknown userdel error: %d", exitCode)
	}
}
