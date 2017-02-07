# docker-node-identities
A tool that will create/remove users and groups on a docker node, if required by a container.

## Flags

```
$ docker-node-identities -help
Usage of docker-node-identities:
-docker-api-version float
      Docker API version (default 1.24)
-docker-cert-path string
      Path to TLS files
-docker-host string
      Daemon socket to connect to (default "unix:///var/run/docker.sock")
-docker-tls-verify
      Use TLS and verify the remote
-group-label string
      Label containing group config (default "docker.node.identities.group")
-user-label string
      Label containing user config (default "docker.node.identities.user")
-version
	  Show version
```

## Labels

### Group
Content of the group label is expected to be `name:password:gid:members` - if you recognize this syntax, it's because you will find it in `/etc/group`.

Note: `members` will be ommited.

### User
Content of the user label is expected to be `username:password:uid:gid:comment:home:shell` - and yes, you might already have guessed it, it's just like `/etc/passwd`.

## Install on docker node
To keep things running, you can use systemd:

```
wget -O /usr/local/bin/docker-node-identities https://github.com/urlund/docker-node-identities/releases/download/1.0.1/docker-node-identities-linux-amd64
wget -O /etc/systemd/system/docker-node-identities.service  https://raw.githubusercontent.com/urlund/docker-node-identities/master/docker-node-identities.service
chmod +x /usr/local/bin/docker-node-identities
service docker-node-identities start
```

Now `docker-node-identities` should be running and ready to keep your users/groups in sync.

## Examples
The most basic example will create a user `johndoe` (uid: `1001`), and a group `johndoe` (gid: `1001`)

```bash
docker run -d -p 80:80 -l "docker.node.identities.user=johndoe::1001" nginx
```

This properly wouldn't make much sense if you are running a single instance, but let's assume you are running a service of 10 (or more) replicas, it would be rather trivial to create the same user on all nodes running or service, so this will make your day:

```bash
docker service create -p 80:80 -l "docker.node.identities.group=thedoes::1010" -l "docker.node.identities.user=johndoe::1001:1010" nginx
```

Now all nodes running a replica will have the user `johndoe` (uid: `1001`) created, and added to the group `thedoes` (gid: `1010`). Scaling the service to 5 and the nodes no longer running a replica will have the user and group removed.
