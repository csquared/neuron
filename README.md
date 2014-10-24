# neuron

<img src="http://www.amrita.edu/sites/default/files/neuron-black-and-white-hi.png" width="300px" />

a process wrapper that pulls ENV and a command from etcd, then
watches etcd for changes and either restarts your process or exits
gracefully


Given the following data in etcd:

    /services/foo-service/envs/dev/PORT = "5000"
    /services/foo-service/envs/dev/DATABASE_URL = "postgres:///foo-service"
    /services/foo-service/processes/web = "bundle exec puma -p $PORT -w 2 -t 12:16"


A call to 

    neuron -env=dev -cmd=web

is like calling:

    /bin/sh -c "bundle exec puma -p $PORT -w 2 -t 12:16"

with an ENV of

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
