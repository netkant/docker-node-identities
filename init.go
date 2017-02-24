package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// ...
var (
	USER_LABEL         string
	GROUP_LABEL        string
	DOCKER_API_VERSION float64
	DOCKER_CERT_PATH   string
	DOCKER_TLS_VERIFY  bool
	VERSION            bool
	DOCKER_HOST        string
)

func init() {
	// ...
	if _default := os.Getenv("GROUP_LABEL"); _default == "" {
		GROUP_LABEL = "docker.node.identities.group"
	}

	// ...
	if _default := os.Getenv("USER_LABEL"); _default == "" {
		USER_LABEL = "docker.node.identities.user"
	}

	// ...
	if _default := os.Getenv("DOCKER_API_VERSION"); _default == "" {
		DOCKER_API_VERSION = 1.24
	}

	// ...
	if _default := os.Getenv("DOCKER_CERT_PATH"); _default == "" {
		DOCKER_CERT_PATH = ""
	}

	// ...
	if _default := os.Getenv("DOCKER_HOST"); _default == "" {
		DOCKER_HOST = "unix:///var/run/docker.sock"
	}

	// ...
	if _default := os.Getenv("DOCKER_TLS_VERIFY"); _default == "" {
		DOCKER_TLS_VERIFY = false
	}

	// ...
	flag.StringVar(&GROUP_LABEL, "group-label", GROUP_LABEL, "Label containing group config")
	flag.StringVar(&USER_LABEL, "user-label", USER_LABEL, "Label containing user config")
	flag.Float64Var(&DOCKER_API_VERSION, "docker-api-version", DOCKER_API_VERSION, "Docker API version")
	flag.StringVar(&DOCKER_CERT_PATH, "docker-cert-path", DOCKER_CERT_PATH, "Path to TLS files")
	flag.StringVar(&DOCKER_HOST, "docker-host", DOCKER_HOST, "Daemon socket to connect to")
	flag.BoolVar(&DOCKER_TLS_VERIFY, "docker-tls-verify", DOCKER_TLS_VERIFY, "Use TLS and verify the remote")
	flag.BoolVar(&VERSION, "version", VERSION, "Show version")
	flag.Parse()

	if VERSION {
		fmt.Println("version: 1.0.2")
		os.Exit(0)
	}

	// ...
	os.Setenv("DOCKER_API_VERSION", strconv.FormatFloat(DOCKER_API_VERSION, 'f', 2, 64))
	os.Setenv("DOCKER_CERT_PATH", DOCKER_CERT_PATH)
	os.Setenv("DOCKER_HOST", DOCKER_HOST)
	os.Setenv("DOCKER_TLS_VERIFY", strconv.FormatBool(DOCKER_TLS_VERIFY))
}
