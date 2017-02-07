package types

// ...
import (
	"fmt"
	"strings"
)

// Group: name:password:gid:members
type Group struct {
	Name     string
	Password string
	GID      string
	Members  string // omitted in create
}

// ...
func NewGroup(data string) (Group, error) {
	var group Group
	parts := strings.Split(data, ":")

	// ...
	if len(parts) > 4 {
		return group, fmt.Errorf("...")
	}

	// ...
	group = Group{
		Name:     getIndexValue(parts, 0, ""),
		Password: getIndexValue(parts, 1, ""),
		GID:      getIndexValue(parts, 2, ""),
		Members:  getIndexValue(parts, 3, ""),
	}

	return group, nil
}

// ...
func (group Group) Create() error {
	// ...
	err := fmt.Errorf("could not create group, no supported command was found")

	// ...
	if path := commandExists("groupadd"); path != "" {
		err = group.groupadd(path)
	}

	return err
}

// ...
func (group Group) Delete() error {
	// ...
	err := fmt.Errorf("could not create group, no supported command was found")

	// ...
	if path := commandExists("groupdel"); path != "" {
		err = group.groupdel(path)
	}

	return err
}

// ...
func (group Group) groupadd(path string) error {
	args := []string{}

	// ...
	if group.Name == "" {
		return fmt.Errorf("name is required")
	}

	// ...
	if group.Password != "" {
		args = append(args, "--password")
		args = append(args, group.Password)
	}

	// ...
	if group.GID != "" {
		args = append(args, "--gid")
		args = append(args, group.GID)
	}

	// the name to add
	args = append(args, group.Name)

	// ...
	_, exitCode := commandExecute(path, args...)
	switch exitCode {
	case 9:
		return fmt.Errorf("group '%s' already exists", group.Name)
	case 0:
		return nil
	default:
		return fmt.Errorf("unknown groupadd error: %d", exitCode)
	}
}

// ...
func (group Group) groupdel(path string) error {
	args := []string{}

	// ...
	if group.Name == "" {
		return fmt.Errorf("name is required")
	}

	// the name to delete
	args = append(args, group.Name)

	// ...
	_, exitCode := commandExecute(path, args...)
	switch exitCode {
	case 6:
		return fmt.Errorf("group '%s' does not exist", group.Name)
	case 0:
		return nil
	default:
		return fmt.Errorf("unknown userdel error: %d", exitCode)
	}
}
