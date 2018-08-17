package types

// ...
import (
	"fmt"
	"os/user"
	"strings"
)

// User is parsed from  "username:password:uid:gid:comment:home:shell"
type User struct {
	Username string
	Password string // omitted (for now)
	UID      string
	GID      string
	Comment  string
	Home     string
	Shell    string
}

// NewUser ...
func NewUser(data string) (User, error) {
	var t User
	parts := strings.Split(data, ":")

	// ...
	if len(parts) > 7 {
		return t, fmt.Errorf("user string should only consist of 7 parts")
	}

	// ...
	t = User{
		Username: getIndexValue(parts, 0, ""),
		Password: getIndexValue(parts, 1, ""),
		UID:      getIndexValue(parts, 2, ""),
		GID:      getIndexValue(parts, 3, ""),
		Comment:  getIndexValue(parts, 4, ""),
		Home:     getIndexValue(parts, 5, ""),
		Shell:    getIndexValue(parts, 6, ""),
	}

	return t, nil
}

// Create ...
func (t User) Create() error {
	// ...
	err := fmt.Errorf("could not create user, no supported command was found")

	// ...
	if t.Username == "" {
		return fmt.Errorf("username is required")
	}

	// ...
	if _, err := user.Lookup(t.Username); err == nil {
		return fmt.Errorf("user '%s' already exists", t.Username)
	}

	// ...
	if t.UID != "" {
		if _, err := user.LookupId(t.UID); err == nil {
			return fmt.Errorf("user with id '%s' already exists", t.UID)
		}
	}

	// ...
	if t.GID == "" {
		t.GID = t.UID
	}

	// ...
	if _, err := user.LookupGroupId(t.GID); err != nil {
		return fmt.Errorf("group does not exist")
	}

	// ...
	if path := commandExists("useradd"); path != "" {
		err = t.useradd(path)
	}

	return err
}

// Delete ...
func (t User) Delete() error {
	// ...
	err := fmt.Errorf("could not delete user, no supported command was found")

	// ...
	if path := commandExists("userdel"); path != "" {
		err = t.userdel(path)
	}

	return err
}

// useradd ...
func (t User) useradd(path string) error {
	args := []string{}

	// ...
	if t.Password != "" {
		args = append(args, "--password")
		args = append(args, t.Password)
	}

	// ...
	if t.UID != "" {
		args = append(args, "--uid")
		args = append(args, t.UID)
	}

	// ...
	if t.GID != "" {
		args = append(args, "--gid")
		args = append(args, t.GID)
	}

	// ...
	if t.Comment != "" {
		args = append(args, "--comment")
		args = append(args, t.Comment)
	}

	// ...
	if t.Home != "" {
		args = append(args, "--home-dir")
		args = append(args, t.Home)
	} else {
		args = append(args, "--no-create-home")
	}

	// ...
	if t.Shell != "" {
		args = append(args, "--shell")
		args = append(args, t.Shell)
	}

	// the username to add
	args = append(args, t.Username)

	// ...
	_, exitCode := commandExecute(path, args...)
	switch exitCode {
	case 0:
		return nil
	default:
		return fmt.Errorf("unknown useradd error: %d", exitCode)
	}
}

// ...
func (t User) userdel(path string) error {
	args := []string{}

	// ...
	if t.Username == "" {
		return fmt.Errorf("username is required")
	}

	// remove home directory and mail spool
	args = append(args, "--remove")

	// force removal of files, even if not owned by user
	args = append(args, "--force")

	// the username to delete
	args = append(args, t.Username)

	// ...
	_, exitCode := commandExecute(path, args...)
	switch exitCode {
	case 6:
		return fmt.Errorf("user '%s' does not exist", t.Username)
	case 0:
		return nil
	default:
		return fmt.Errorf("unknown userdel error: %d", exitCode)
	}
}
