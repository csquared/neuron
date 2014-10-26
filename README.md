# neuron

<img src="http://www.amrita.edu/sites/default/files/neuron-black-and-white-hi.png" width="300px" />

`neuron` is a UNIX process wrapper that uses configuration
in ectd to run commands and manage ENVs for production services.

Neuron watches etcd for changes and either restarts your process or exits
gracefully when they occur. This allows ENV or command changes
to propagate to your running processes without deploying new code.

## usage

Neuron works well with 12factor apps that expect their config to come from environment.
In development or staging, `neuron -r` can be useful to quickly react to ENV or command
changes.

In production, allowing neuron to crash your processes assumes you have a process manager
responsible for monitoring processes.

```
> neuron
   ____  ___  __  ___________  ____
  / __ \/ _ \/ / / / ___/ __ \/ __ \
 / / / /  __/ /_/ / /  / /_/ / / / /
/_/ /_/\___/\__,_/_/   \____/_/ /_/

Usage of neuron:
  -cmd="": name of cmd key
  -e=".env": .env location for import
  -env="default": name of env dir
  -etcd="http://localhost:4001": url of etcd
  -p="Procfile": procfile location for import
  -r=false: restart instead of crashing
```

Run `neuron import` if you already have a `Procfile` and `.env` file:

```
~/projects/go/src/github.com/csquared/neuron (master*)$ neuron import

   ____  ___  __  ___________  ____
  / __ \/ _ \/ / / / ___/ __ \/ __ \
 / / / /  __/ /_/ / /  / /_/ / / / /
/_/ /_/\___/\__,_/_/   \____/_/ /_/

action=import procfile=Procfile envfile=.env
action=import-procfile process=web
action=import-procfile process=worker
action=import-env-var key=WEB_URL
action=import-env-var key=FOO
```

Given the following data in etcd:

    /services/foo-service/envs/dev/PORT = "5000"
    /services/foo-service/envs/dev/DATABASE_URL = "postgres:///foo-service"
    /services/foo-service/processes/web = "bundle exec puma -p $PORT -w 2 -t 12:16"

A call to

    > neuron -env=dev -cmd=web

in the directory `foo-service` is like calling:

    /bin/sh -c "bundle exec puma -p $PORT -w 2 -t 12:16"

in that directory with an ENV of

    PORT="5000"
    DATABASE_URL="postgres:///foo-service"

## options

### -cmd

Takes the fully qualified key in etcd or shorthand

### -env

Takes the fully qualified directory in etcd or shorthand

### -etcd

URL of etcd

### -r

restart process - useful for development or tuning params

if this is not set, neuron crashes and assumes your process
manager will reboot it. when your system restarts the
neuron process it will have the new ENV

## other commands

### neuron import

Loads Procfile and .env files into etcd

#### -p

Name of Procfile

#### -e

Name of .env file

### neuron bootstrap

Creates directories and single web process in etcd
