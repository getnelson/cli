# nelson-cli

[![Build Status](https://travis.oncue.verizon.net/iptv/nelson-cli.svg?token=Lp2ZVD96vfT8T599xRfV&branch=master)](https://travis.oncue.verizon.net/iptv/nelson-cli)

## Getting Started

Installing and using the client is super easy:

1. Download the latest release from [the nexus](http://nexus.oncue.verizon.net/nexus/content/groups/internal/verizon/inf/nelson/cli/)
2. Stuff it into your $PATH (e.g. `/usr/local/bin`)
3. Make it executable `sudo chmod +x /usr/local/bin/nelson`)

Then you're ready to use the CLI. The first command you should execute after install is `login` which allows you to securely interact with the remote *Nelson* service. To do this, you just need to follow these steps:

1. [Obtain a Github personal access token](https://help.github.com/articles/creating-an-access-token-for-command-line-use/)
2. Set the Github token into your environment: `export GITHUB_TOKEN=XXXXXXXXXXXXXXXX`
3. `nelson login nelson.yourcompany.com`, then you're ready to start using the other commands! If you're running the *Nelson* service insecurely - without SSL - then you need to pass the `--disable-tls` flag to the login command.

> ⛔ Note that currently the Nelson client can only be logged into *one* remote *Nelson* service at a time. ⛔

The below set of commands are the currently implemented set:

```
# fully explicit login
$ nelson login --token 1f3f3f3f3 nelson.yourdomain.com

# read token from environment variable GITHUB_TOKEN, explicit host
$ nelson login nelson.yourdomain.com

# read token from env var GITHUB_TOKEN and host from NELSON_ADDR
$ nelson login

# for testing with a local server, you can do:
$ nelson login --disable-tls --token 1f3f3f3f3 nelson.local:9000

# list the available nelson datacenters
$ nelson datacenters list

# just an alias for the above
$ nelson dcs list

# show the deployment log for a given deployment id
$ nelson stacks fs 02481438b432

# display the current user information
$ nelson whoami

# redeploy a very specific deployment id
$ nelson stacks redeploy b8ff485a0306

# show you the units deployed in a given datacenter
$ nelson units list --datacenter us-west-2

# list the available clean policies on this remote nelson
$ nelson system cleanup-policies
```

The following commands are currently being developed:

```
# list the available workflows
$ nelson system workflows

$ nelson datacenter inspect <arg>

# inspect a specific unit; showing dependencies and crap
$ nelson units inspect howdy-batch-0.3

# list me all the deployments, in all datacenters for this unit name
$ nelson stacks list --unit howdy-batch-0.3

# inspect a very specific deployment
$ nelson stacks inspect b8ff485a0306

# list the workflows availabe in the remote nelson
$ nelson workflows inspect --type job quasar

```

## Development


1. `brew install go` - install the Go programming language:
1. `go get https://getgb.io/` - instal the `gb` build tool
1. `go get https://github.com/codeskyblue/fswatch` - install `fswatch` so we can do continous compilation
1. `alias fswatch="$GOPATH/bin/fswatch"
1. `fswatch`

This should give continous compilation without the tedious need to constantly restart `gb build`
