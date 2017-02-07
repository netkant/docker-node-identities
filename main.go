package main

// ...
import (
	"context"
	"fmt"
	"os"

	dtypes "github.com/docker/docker/api/types"
	devents "github.com/docker/docker/api/types/events"
	dfilters "github.com/docker/docker/api/types/filters"
	dclient "github.com/docker/docker/client"
	ctypes "github.com/urlund/docker-node-identities/types"
)

// ...
func main() {
	// ...
	if client, err := dclient.NewEnvClient(); err != nil {
		fmt.Printf("%s\n", err)
	} else {
		// ...
		go sync(client)

		// ...
		listen(client)
	}
}

// ...
func sync(client *dclient.Client) {

	// ...
	containers, err := getContainers(client)

	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// ...
	for _, container := range containers {
		if err := create(container.Labels); err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}

// ...
func listen(client *dclient.Client) {
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
			if err != nil {
				fmt.Printf("%s\n", err)
			}
		case e := <-messages:
			switch e.Action {
			case "start": // always triggered when a container starts running
				if container, err := getContainer(client, e.ID); err == nil {
					if err := create(container.Config.Labels); err != nil {
						fmt.Printf("%s\n", err)
					}
				} else {
					fmt.Printf("%s\n", err)
				}
			case "die": // other actions would be "kill" or "stop", but "die" is always triggered when a container exits
				if container, err := getContainer(client, e.ID); err == nil {
					if err := delete(client, container.Config.Labels); err != nil {
						fmt.Printf("%s\n", err)
					}
				} else {
					fmt.Printf("%s\n", err)
				}
			}
		}
	}
}

// ...
func parseUser(labels map[string]string) (ctypes.User, error) {
	var user ctypes.User

	data, exists := labels[CUDDY_USER_LABEL]
	if !exists {
		return user, fmt.Errorf("user label does not exist")
	}

	//
	return ctypes.NewUser(data)
}

// ...
func parseGroup(labels map[string]string) (ctypes.Group, error) {
	var group ctypes.Group

	data, exists := labels[CUDDY_GROUP_LABEL]
	if !exists {
		return group, fmt.Errorf("group label does not exist")
	}

	//
	return ctypes.NewGroup(data)
}

// ...
func create(labels map[string]string) error {
	group, err := parseGroup(labels)
	if err == nil {
		if err := group.Create(); err != nil {
			return err
		} else {
			fmt.Printf("group '%s' was created\n", group.Name)
		}
	}

	if user, err := parseUser(labels); err == nil {
		if group.GID != "" {
			user.GID = group.GID
		}
		if err := user.Create(); err != nil {
			return err
		} else {
			fmt.Printf("user '%s' was created\n", user.Username)
		}
	} else {
		return err
	}

	return nil
}

// ...
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
				return fmt.Errorf("user '%s' was not deleted (defined in another label)", user.Username)
			}
		}

		if err := user.Delete(); err != nil {
			return err
		} else {
			fmt.Printf("user '%s' was deleted\n", user.Username)
		}
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
				return fmt.Errorf("group '%s' was not deleted (defined in another label)", group.Name)
			}
		}

		if err := group.Delete(); err != nil {
			return err
		} else {
			fmt.Printf("group '%s' was deleted\n", group.Name)
		}
	}

	return nil
}

// ...
func getContainers(client *dclient.Client) ([]dtypes.Container, error) {
	// ...
	filters := dfilters.NewArgs()
	filters.Add("label", CUDDY_USER_LABEL)

	containerListOptions := dtypes.ContainerListOptions{
		Filters: filters,
	}

	// ...
	return client.ContainerList(context.Background(), containerListOptions)
}

// ...
func getContainer(client *dclient.Client, id string) (dtypes.ContainerJSON, error) {
	return client.ContainerInspect(context.Background(), id)
}
