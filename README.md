# Nelson CLI

[![Build Status](https://travis-ci.org/Verizon/nelson-cli.svg?branch=master)](https://travis-ci.org/Verizon/nelson-cli)
[![Latest Release](https://img.shields.io/github/release/verizon/nelson-cli.svg)](https://github.com/Verizon/nelson-cli/releases)

A handy tool for interacting with your installation of [Nelson](https://github.com/Verizon/nelson)

## Getting Started

If you just want to use nelson-cli, then run the following:

```
curl -GqL https://raw.githubusercontent.com/Verizon/nelson-cli/master/scripts/install | bash
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

The below set of commands are the currently implemented set - node that for subcommands, both plural and singular command verbs work. For example `stacks` and `stack` are functionallty identical:

### Global Flags
```
# print debug output for network request
$ nelson --debug <command>

# print analogous curl command for network request
$ nelson --debug-curl <command>
```

### System Operations

```
# display the current user information
$ nelson whoami

# display info for the currently deployed Nelson
$ nelson system version

# list the available clean-up policies on this remote nelson
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

### Repository Operations

```
# enable a repo for use with nelson
$ nelson repo enable --owner octocat --repo helloworld

# disable a repo
$ nelson repo disable --owner octocat --repo helloworld
```

### Datacenter Operations

```
# list the available nelson datacenters
$ nelson datacenters list

# just an alias for the above
$ nelson dcs list
```

### Namespace Operations

```
# create an entirely new namespace
$ nelson namespace create --datacenter dc123 --namespace foobar
$ nelson ns create --dc dc123 --ns foobar
```

### Unit Operations

```
# show the units deployed in a given datacenter
$ nelson units list --namespaces dev --datacenters sacremento

# show the units available in several datacenters
$ nelson units list --namespaces dev --datacenters sacremento,nyc

# show the units available in all datacenters for a given namespace
$ nelson units list --namespaces dev

# show the units available in all datacenters for a given namespace and status
$ nelson units list --namespaces dev --statuses deploying,ready,deprecated

# show the units that have been terminated by nelson in a given namespace
$ nelson units list --namespaces dev --statuses terminated

# deprecate a specific unit and feature version
$ nelson units deprecate --unit foo --version 1.2

# deprecate a specific unit and feature version, and expire the unit right away
$ nelson units deprecate --no-grace --unit foo --version 1.2

# take a deployment from one namespace and commit it to the specified target namespace
$ nelson units commit --foo --version 1.2.3 --target qa

```

### Stack Operations

```
# show the stacks deployed in a given datacenter
$ nelson stacks list --namespaces dev --datacenters sacremento

# show the stacks available in several datacenters
$ nelson stacks list --namespaces dev --datacenters sacremento,nyc

# show the stacks available in all datacenters for a given namespace
$ nelson stacks list --namespaces dev

# show the stacks available in all datacenters for a given namespace and status
$ nelson stacks list --namespaces dev --statuses deploying,ready,deprecated

# show the stacks that have been terminated by nelson in a given namespace
$ nelson stacks list --namespaces dev --statuses terminated

# inspect a very specific deployment and show more detailed routing information
$ nelson stacks inspect b8ff485a0306

# redeploy a very specific deployment id.
# this spawns a new stack using the exact same container image
$ nelson stacks redeploy b8ff485a0306

# show the deployment log for a given deployment id
$ nelson stacks fs 02481438b432
$ nelson stacks logs 02481438b432

# show the current *runtime* status as seen by consul and nomad
$ nelson stacks runtime 02481438b432

# manually register a stack - only needed to inform Nelson about 
# something you setup out of band (e.g. some static database)
$ nelson stacks manual \
  --datacenter sacremento \
  --namespace dev \
  --service-type zookeeper \
  --version 3.4.6 \
  --hash fj5xlf53q \
  --description "some zookeeper" \
  --port 2181
```

### Loadbalancer Operations

```
# list the loadbalancers
nelson lbs list -ns dev -d sacremento
nelson lbs list -ns dev

# remove a loadbalancer
nelson lbs down 04dsq452xvq

# create a new loadbalancer
nelson lbs up --name howdy-lb --major-version 1 --datacenter us-east-1 --namespace dev
nelson lbs up -n howdy-lb -mv 1 -d us-east-1 -ns dev

# get info about a loadbalancer
nelson lbs inspect -guid 04dsq452xvq
```

## Lint operations

### Templates

Testing consul templates is tedious, because many require vault access and/or Nelson environment variables to render.  nelson-cli can render your consul-template in an environment similar to your container.  Specifically, it:

* creates a vault token with the same permissions as your unit and resources
* renders the template server-side with the vault token and a full set of NELSON environment variables for the dev namespace
* shows the error output, if any

Given this template at `application.cfg.template`:

```
{{with $ns := env "NELSON_ENV"}}
{{with secret (print "someroot/" $ns "/test/creds/howdy-http")}}
username={{.data.username}}
{{end}}
{{end}}
```

Lint it as `howdy-http` unit, using resource `test`:

```
$ nelson lint template -u <nelson-unit-name> -r <resource-name> -t <path-to>/application.cfg.template
template rendering failed
2017/02/15 18:54:15.496679 [INFO] consul-template v0.18.1 (9c62737)
2017/02/15 18:54:15.496716 [INFO] (runner) creating new runner (dry: true, once: true)
2017/02/15 18:54:15.497461 [INFO] (runner) creating watcher
2017/02/15 18:54:15.497809 [INFO] (runner) starting
2017/02/15 18:54:15.497884 [INFO] (runner) initiating run
2017/02/15 18:54:15.999977 [INFO] (runner) initiating run
Consul Template returned errors:
/consul-template/templates/nelson7713234105042928921.template: execute: template: :3:16: executing "" at <.data.username>: can't evaluate field data in type *dependency.Secret

Template linting failed.
```

Oops.  Line 3 of the template should be `.Data`, not `.data`.  Fix it and try again:

```
$ nelson lint template -u <nelson-unit-name> -r <resource-name> -t <path-to>/application.cfg.template
Template rendered successfully.
Rendered output discarded for security reasons.
```

The template rendered, but because we don't want to expose any secrets, we get a simple success message.  Congratulations.  Your template should now render correctly when deployed by Nelson.

## Development

1. `brew install go` - install the Go programming language:
1. create a directory to contain your go projects
1. in your .bashrc or .zshrc, add
    `export GOPATH=~/[path_to_go_directory]/go`
    `export PATH=$PATH:$GOPATH/bin`
    `ulimit -n 7168`
1. `make install-dev` - install the `gb` build tool
1. `alias fswatch="$GOPATH/bin/fswatch"`
1. `make watch`

This should give continuous compilation without the tedious need to constantly restart `gb build`
