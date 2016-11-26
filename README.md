# cuddy

...

## Flags

|Flag|Type|Description|Default|
|---|---|---|---|
|-cuddy-group-label|string|Label containing group config|io.cuddy.group|
|-cuddy-user-label|string|Label containing user config|io.cuddy.user|
|-docker-api-version|float|Docker API version|1.24|
|-docker-cert-path|string|Path to TLS files||
|-docker-host|string|Daemon socket to connect to|unix:///var/run/docker.sock|
|-docker-tls-verify|bool|Use TLS and verify the remote||

## Labels

### `io.cuddy.group`
Content of the group label is expected to be `name:password:gid:members` - if you recognize this syntax, it's because you will find it in `/etc/group`.

Note: `members` will be ommited.

### `io.cuddy.user`
Content of the user label is expected to be `username:password:uid:gid:comment:home:shell` - and yes, you might already have guessed it, it's just like `/etc/passwd`.

## Examples
The most basic example will create a user `johndoe` (uid: `1001`), and a group `johndoe` (gid: `1001`)
```bash
docker run -d -p 80:80 -l "io.cuddy.user=johndoe::1001" nginx
```
This properly wouldn't make much sense if you are running a single instance, but let's assume you are running a service of 10 (or more) replicas, it would be rather trivial to create the same user on all nodes running or service, so this would make your day:
```bash
docker service create -p 80:80 -l "io.cuddy.group=thedoes::1010" -l "io.cuddy.user=johndoe::1001:1010" nginx
```
Now all nodes running a replica will have the user `johndoe` (uid: `1001`) created, and added to the group `thedoes` (gid: `1010`). Scaling the service to 5 and the nodes no longer running a replica will have the user and group removed.
