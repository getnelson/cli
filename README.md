# nelson-cli

[![Build Status](https://travis.oncue.verizon.net/iptv/nelson-cli.svg?token=Lp2ZVD96vfT8T599xRfV&branch=master)](https://travis.oncue.verizon.net/iptv/nelson-cli)

## Getting Started

If you just want to use nelson-cli, then run the following:

```
curl -GL https://is.gd/5C93hC | bash
```

This script will download and install the latest version and put it on your `$PATH`. We do not endorse piping scripts from the wire to `bash`, and you should read the script before executing the command. It will:

1. Fetch the latest version from Nexus
2. Verify the SHA1 sum
3. Extract the tarball
4. Copy nelson to `/usr/local/bin/nelson`

It is safe to rerun this script to keep nelson-cli current. If you have the source code checked out locally, you need only execute: `scripts/install-nelson-cli` to install the latest version of nelson-cli.  

Then you're ready to use the CLI. The first command you should execute after install is `login` which allows you to securely interact with the remote *Nelson* service. To do this, you just need to follow these steps:

1. [Obtain a Github personal access token](https://help.github.com/articles/creating-an-access-token-for-command-line-use/)
2. Set the Github token into your environment: `export GITHUB_TOKEN=XXXXXXXXXXXXXXXX`
3. `nelson login nelson.yourcompany.com`, then you're ready to start using the other commands! If you're running the *Nelson* service insecurely - without SSL - then you need to pass the `--disable-tls` flag to the login command.

> ⛔ Note that currently the Nelson client can only be logged into *one* remote *Nelson* service at a time. ⛔

The below set of commands are the currently implemented set:

### System Operations

```
# display the current user information
$ nelson whoami

# list the available clean policies on this remote nelson
$ nelson system cleanup-policies

# fully explicit login
$ nelson login --token 1f3f3f3f3 nelson.yourdomain.com

# read token from environment variable GITHUB_TOKEN, explicit host
$ nelson login nelson.yourdomain.com

# read token from env var GITHUB_TOKEN and host from NELSON_ADDR
$ nelson login

# for testing with a local server, you can do:
$ nelson login --disable-tls --token 1f3f3f3f3 nelson.local:9000
```

### Datacenter Operations

```
# list the available nelson datacenters
$ nelson datacenters list

# just an alias for the above
$ nelson dcs list
```

### Unit Operations

```
# show the units deployed in a given datacenter
$ nelson units list --namespaces dev --datacenters us-west-2

# show the units availabe in several datacenters
$ nelson units list --namespaces dev --datacenters us-west-2,us-west-1

# show the units availabe in all datacenters for a given namespace
$ nelson units list --namespaces dev

# show the units availabe in all datacenters for a given namespace and status
$ nelson units list --namespaces dev --statuses deploying,ready,deprecated

# show the units that have been terminated by nelson in a given namespace
$ nelson units list --namespaces dev --statuses terminated

# deprecate a specific unit and feature version
$ nelson units deprecate --unit foo --version 1.2

```

### Stack Operations

```
# show the stacks deployed in a given datacenter
$ nelson stacks list --namespaces dev --datacenters us-west-2

# show the stacks availabe in several datacenters
$ nelson stacks list --namespaces dev --datacenters us-west-2,us-west-1

# show the stacks availabe in all datacenters for a given namespace
$ nelson stacks list --namespaces dev

# show the stacks availabe in all datacenters for a given namespace and status
$ nelson stacks list --namespaces dev --statuses deploying,ready,deprecated

# show the stacks that have been terminated by nelson in a given namespace
$ nelson stacks list --namespaces dev --statuses terminated

# redeploy a very specific deployment id
$ nelson stacks redeploy b8ff485a0306

# inspect a very specific deployment
$ nelson stacks inspect b8ff485a0306

# show the deployment log for a given deployment id
$ nelson stacks fs 02481438b432
```

### Loadbalancer Operations

```
# list the loadbalancers
nelson lbs list -ns dev -d us-east-1
nelson lbs list -ns dev

# remove a loadbalancer
nelson lbs down 04dsq452xvq

# create a new loadbalancer
nelson lbs up --name howdy-lb --major-version 1 --datacenter us-east-1 --namespace dev
nelson lbs up -n howdy-lb -mv 1 -d us-east-1 -ns dev
```

The following commands are currently being developed:

```
$ nelson stack deploy --manual

# list the workflows availabe in the remote nelson
$ nelson workflows inspect --type job quasar

# inspect a specific unit; showing dependencies and crap
$ nelson units inspect howdy-batch-0.3

# list me all the deployments, in all datacenters for this unit name
$ nelson stacks list --unit howdy-batch-0.3

```

## Development


1. `brew install go` - install the Go programming language:
1. `go get github.com/constabulary/gb` - install the `gb` build tool
1. `go get github.com/codeskyblue/fswatch` - install `fswatch` so we can do continous compilation
1. `alias fswatch="$GOPATH/bin/fswatch"`
1. `fswatch`

This should give continous compilation without the tedious need to constantly restart `gb build`
