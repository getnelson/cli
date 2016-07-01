# nelson-cli

[![Build Status](https://travis.oncue.verizon.net/iptv/nelson-cli.svg?token=Lp2ZVD96vfT8T599xRfV&branch=master)](https://travis.oncue.verizon.net/iptv/nelson-cli)

Command line client for the Nelson deployment API

```
IMPLEMENTED

# fully explicit login
$ nelson login --token 1f3f3f3f3 nelson.yourdomain.com

# read token from environment variable GITHUB_TOKEN, explicit host
$ nelson login nelson.yourdomain.com

# read token from env var GITHUB_TOKEN and host from NELSON_ADDR
$ nelson login

# for testing with a local server, you can do:
$ nelson login --disable-tls nelson.local:9000

# list the available nelson datacenters
$ nelson datacenter list

# just an alias for the above
$ nelson dc list

# show the deployment log for a given deployment id
$ nelson stack fs 1234

# display the current user information
$ nelson whoami

# redeploy a very specific deployment id
$ nelson stack redeploy 1234

TODO

$ nelson datacenter inspect <arg>

# show you the units deployed in a given datacenter
$ nelson unit list --datacenter us-west-2

# inspect a specific unit; showing dependencies and crap
$ nelson unit inspect howdy-batch-0.3

# list me all the deployments, in all datacenters for this unit name
$ nelson stack list --unit howdy-batch-0.3

# inspect a very specific deployment
$ nelson stack inspect 1234

```

# Development

1. Install `gb`: https://getgb.io/
1. Install https://github.com/codeskyblue/fswatch
1. `alias fswatch="$GOPATH/bin/fswatch"
1. `fswatch`

This should give continous compilation without the tedious need to constantly restart `gb build`