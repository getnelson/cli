# nelson-cli

[![Build Status](https://travis.oncue.verizon.net/iptv/nelson-cli.svg?token=Lp2ZVD96vfT8T599xRfV&branch=master)](https://travis.oncue.verizon.net/iptv/nelson-cli)

Command line client for the Nelson deployment API

```
$ nelson login

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