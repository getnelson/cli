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

TODO

$ nelson dc inspect us-west-2

$ nelson unit inspect 1234

$ nelson dpl inspect 1234
$ nelson dpl replay 1234

```

# Development

1. Install `gb`: https://getgb.io/
1. Install https://github.com/codeskyblue/fswatch
1. `alias fswatch="$GOPATH/bin/fswatch"
1. `fswatch`

This should give continous compilation without the tedious need to constantly restart `gb build`