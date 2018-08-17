package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	clog "github.com/urlund/docker-node-identities/log"
)

// ...
var (
	DockerAPIVersion float64
	DockerCertPath   string
	DockerHost       string
	DockerTLSVerify  bool
	Debug            bool
	GroupLabel       string
	UserLabel        string
	Version          bool
)

func init() {
	// ...
	if _default := os.Getenv("GROUP_LABEL"); _default == "" {
		GroupLabel = "docker.node.identities.group"
	}

	// ...
	if _default := os.Getenv("USER_LABEL"); _default == "" {
		UserLabel = "docker.node.identities.user"
	}

	// ...
	if _default := os.Getenv("DOCKER_API_VERSION"); _default == "" {
		DockerAPIVersion = 1.24
	}

	// ...
	if _default := os.Getenv("DOCKER_CERT_PATH"); _default == "" {
		DockerCertPath = ""
	}

	// ...
	if _default := os.Getenv("DOCKER_HOST"); _default == "" {
		DockerHost = "unix:///var/run/docker.sock"
	}

	// ...
	if _default := os.Getenv("DOCKER_TLS_VERIFY"); _default == "" {
		DockerTLSVerify = false
	}

	// ...
	flag.StringVar(&GroupLabel, "group-label", GroupLabel, "Label containing group config")
	flag.StringVar(&UserLabel, "user-label", UserLabel, "Label containing user config")
	flag.Float64Var(&DockerAPIVersion, "docker-api-version", DockerAPIVersion, "Docker API version")
	flag.StringVar(&DockerCertPath, "docker-cert-path", DockerCertPath, "Path to TLS files")
	flag.StringVar(&DockerHost, "docker-host", DockerHost, "Daemon socket to connect to")
	flag.BoolVar(&DockerTLSVerify, "docker-tls-verify", DockerTLSVerify, "Use TLS and verify the remote")
	flag.BoolVar(&Version, "version", Version, "Show version")
	flag.BoolVar(&Debug, "debug", Debug, "Show debug info")
	flag.Parse()

	if Version {
		fmt.Println("version: 1.0.4")
		os.Exit(0)
	}

	clog.Settings = &clog.LogSettings{
		Debug: Debug,
	}

	// ...
	os.Setenv("DOCKER_API_VERSION", strconv.FormatFloat(DockerAPIVersion, 'f', 2, 64))
	os.Setenv("DOCKER_CERT_PATH", DockerCertPath)
	os.Setenv("DOCKER_HOST", DockerHost)
	os.Setenv("DOCKER_TLS_VERIFY", strconv.FormatBool(DockerTLSVerify))
}
