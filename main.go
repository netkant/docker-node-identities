package main

// ...
import (
	"context"
	"fmt"
	// "os"
	"time"

	dtypes "github.com/docker/docker/api/types"
	devents "github.com/docker/docker/api/types/events"
	dfilters "github.com/docker/docker/api/types/filters"
	dclient "github.com/docker/docker/client"
	clog "github.com/urlund/docker-node-identities/log"
	ctypes "github.com/urlund/docker-node-identities/types"
)

// main ...
func main() {
	for {
		if client, err := dclient.NewEnvClient(); err == nil {
			if _, err := client.Ping(context.Background()); err == nil {
				clog.Debug("connected to docker")

				// ...
				go sync(client)

				// ...
				if err := listen(client); err != nil {
					clog.Info("%s", err)
				}
			} else {
				clog.Info("%s", err)
				time.Sleep(5000 * time.Millisecond)
			}
		} else {
			clog.Info("%s", err)
		}
	}
}

// sync ...
func sync(client *dclient.Client) {
	// ...
	if containers, err := getContainers(client); err == nil {
		if len(containers) > 0 {
			clog.Debug("checking already running containers")
		}

		// ...
		for _, container := range containers {
			create(container.Labels)
		}
	} else {
		clog.Info("%s", err)
	}
}

// listen ...
func listen(client *dclient.Client) error {
	// make sure we are only listening to container events
	filters := dfilters.NewArgs()
	filters.Add("type", devents.ContainerEventType)

	// get event channel
	messages, errs := client.Events(context.Background(), dtypes.EventsOptions{
		Filters: filters,
	})

	for {
		select {
		case err := <-errs:
			return err
		case e := <-messages:
			switch e.Action {
			case "start":
				if container, err := getContainer(client, e.ID); err == nil {
					create(container.Config.Labels)
				}
			case "die", "kill", "stop":
				if container, err := getContainer(client, e.ID); err == nil {
					if err := delete(client, container.Config.Labels); err != nil {
						clog.Debug("container %s: %s", container.Name, err)
					}
				}
			}
		}
	}
}

// parseUser ...
func parseUser(labels map[string]string) (ctypes.User, error) {
	var user ctypes.User

	data, exists := labels[UserLabel]
	if !exists {
		return user, fmt.Errorf("user label '%s' does not exist", UserLabel)
	}

	//
	return ctypes.NewUser(data)
}

// parseGroup ...
func parseGroup(labels map[string]string) (ctypes.Group, error) {
	var group ctypes.Group

	data, exists := labels[GroupLabel]
	if !exists {
		return group, fmt.Errorf("group label '%s' does not exist", GroupLabel)
	}

	return ctypes.NewGroup(data)
}

// create ...
func create(labels map[string]string) {
	if err := createGroup(labels); err != nil {
		clog.Debug("%s", err)
	}

	if err := createUser(labels); err != nil {
		clog.Debug("%s", err)
	}
}

// createGroup ...
func createGroup(labels map[string]string) error {
	group, err := parseGroup(labels)
	if err == nil {
		if err := group.Create(); err != nil {
			return err
		}

		clog.Info("group '%s' was created", group.Name)
	}

	return err
}

// create ...
func createUser(labels map[string]string) error {
	user, err := parseUser(labels)
	if err == nil {
		if err := user.Create(); err != nil {
			return err
		}

		clog.Info("user '%s' was created", user.Username)
	}

	return err
}

// delete ...
func delete(client *dclient.Client, labels map[string]string) error {
	// ...
	containers, err := getContainers(client)
	if err != nil {
		return err
	}

	// ...
	if user, err := parseUser(labels); err == nil {
		// check if user is defined in another label, to cancel delete
		for _, container := range containers {
			_user, err := parseUser(container.Labels)
			if err != nil {
				return err
			}

			if _user.Username == user.Username {
				return fmt.Errorf("user '%s' was not deleted (also defined by another container)", user.Username)
			}
		}

		if err := user.Delete(); err != nil {
			return err
		}

		clog.Info("user '%s' was deleted", user.Username)
	} else {
		return err
	}

	if group, err := parseGroup(labels); err == nil {
		// check if group is defined in another label, to cancel delete
		for _, container := range containers {
			_group, err := parseGroup(container.Labels)
			if err != nil {
				return err
			}

			if _group.Name == group.Name {
				return fmt.Errorf("group '%s' was not deleted (also defined by another container)", group.Name)
			}
		}

		if err := group.Delete(); err != nil {
			return err
		}

		clog.Info("group '%s' was deleted", group.Name)
	}

	return nil
}

// getContainers ...
func getContainers(client *dclient.Client) ([]dtypes.Container, error) {
	// ...
	filters := dfilters.NewArgs()
	filters.Add("label", UserLabel)

	containerListOptions := dtypes.ContainerListOptions{
		Filters: filters,
	}

	return client.ContainerList(context.Background(), containerListOptions)
}

// getContainer ...
func getContainer(client *dclient.Client, id string) (dtypes.ContainerJSON, error) {
	return client.ContainerInspect(context.Background(), id)
}
