package types

// ...
import (
	"fmt"
	"strings"
)

// User: username:password:uid:gid:comment:home:shell
type User struct {
	Username string
	Password string // omitted (for now)
	UID      string
	GID      string
	Comment  string
	Home     string
	Shell    string
}

// ...
func NewUser(data string) (User, error) {
	var user User
	parts := strings.Split(data, ":")

	// ...
	if len(parts) > 7 {
		return user, fmt.Errorf("...")
	}

	// ...
	user = User{
		Username: getIndexValue(parts, 0, ""),
		Password: getIndexValue(parts, 1, ""),
		UID:      getIndexValue(parts, 2, ""),
		GID:      getIndexValue(parts, 3, ""),
		Comment:  getIndexValue(parts, 4, ""),
		Home:     getIndexValue(parts, 5, ""),
		Shell:    getIndexValue(parts, 6, ""),
	}

	return user, nil
}

// ...
func (user User) Create() error {
	// ...
	err := fmt.Errorf("could not create user, no supported command was found")

	// ...
	if path := commandExists("useradd"); path != "" {
		err = user.useradd(path)
	}

	return err
}

// ...
func (user User) Delete() error {
	// ...
	err := fmt.Errorf("could not delete user, no supported command was found")

	// ...
	if path := commandExists("userdel"); path != "" {
		err = user.userdel(path)
	}

	return err
}

// ...
func (user User) useradd(path string) error {
	args := []string{}

	// ...
	if user.Username == "" {
		return fmt.Errorf("username is required")
	}

	// ...
	if user.Password != "" {
		args = append(args, "--password")
		args = append(args, user.Password)
	}

	// ...
	if user.UID != "" {
		args = append(args, "--uid")
		args = append(args, user.UID)
	}

	// ...
	if user.GID != "" {
		args = append(args, "--gid")
		args = append(args, user.GID)
	}

	// ...
	if user.Comment != "" {
		args = append(args, "--comment")
		args = append(args, user.Comment)
	}

	// ...
	if user.Home != "" {
		args = append(args, "--home-dir")
		args = append(args, user.Home)
	} else {
		args = append(args, "--no-create-home")
	}

	// ...
	if user.Shell != "" {
		args = append(args, "--shell")
		args = append(args, user.Shell)
	}

	// the username to add
	args = append(args, user.Username)

	// ...
	_, exitCode := commandExecute(path, args...)
	switch exitCode {
	case 9:
		return fmt.Errorf("user '%s' already exists", user.Username)
	case 0:
		return nil
	default:
		return fmt.Errorf("unknown useradd error: %d", exitCode)
	}
}

// ...
func (user User) userdel(path string) error {
	args := []string{}

	// ...
	if user.Username == "" {
		return fmt.Errorf("username is required")
	}

	// remove home directory and mail spool
	args = append(args, "--remove")

	// force removal of files, even if not owned by user
	args = append(args, "--force")

	// the username to delete
	args = append(args, user.Username)

	// ...
	_, exitCode := commandExecute(path, args...)
	switch exitCode {
	case 6:
		return fmt.Errorf("user '%s' does not exist", user.Username)
	case 0:
		return nil
	default:
		return fmt.Errorf("unknown userdel error: %d", exitCode)
	}
}
