# isosceles
> Remote Development Tool

[![Build Status](https://travis-ci.org/murdinc/isosceles.svg)](https://travis-ci.org/murdinc/isosceles)

## Intro
**isosceles** is a tool built out of a desire to work locally in an environment set up for remote development. It's basically a wrapper for rsync that watches for file changes.

To set up projects, add an .isosceles file to your user folder on OS X. An example is located in the example_config.isosceles file in this repo and at the end of the README.

## Features
* Watches an entire directory recursively for changes that match a specific pattern
* Kicks off an rsync when a trigger is detected.
* Pools triggers that happen during the cooldown period (set in the config), to keep from repeating useless syncs.
* rsync flags are fully customizable in the config
* Desktop Notifications can be enabled, with or without sound, for when triggers are processed.

## Installation
```
curl -s http://dl.sudoba.sh/get/isosceles | sh
```



## Development
These steps are only required if you wish to develop on the isosceles code, or compile it yourself.

1. Install Go (if you haven't already): https://golang.org/doc/install

2. Download and install isosceles:

  `$ go get -u github.com/murdinc/isosceles`

3. If `$GOPATH/bin` is not yet in your `$PATH` (check with `$ echo $PATH`):

  `$ export PATH="$PATH:$GOPATH/bin"`

4. Test it! `$ go test github.com/murdinc/isosceles/...`

```
$ go test github.com/murdinc/isosceles/...
ok      github.com/murdinc/isosceles    0.010s
```



## Screenshots

**CLI Menu:**

![screenshot1](screenshots/help.png)

**Listing All Projects:**

![screenshot1](screenshots/active-projects.png)

**Active Sync:**

![screenshot1](screenshots/active-sync.png)

**Desktop Notifications:**

![screenshot1](screenshots/desktop-notification.png)

**Example configuration: (goes in ~/.isosceles)**


## Example Config

```

[project "Project1"]

    # Enabled Flag
    enabled = true

    # Host and Folders
    host = "host.name.com"
    local-folder = "/Users/USER/isosceles/PROJECT/"
    remote-folder = "/PROJECT/"
    url = "http://www.host.name.com"

    # Watch Pattern
    watch-pattern = "(.php|.css|.js|.sass|.html)"

    # Rsync Arguments
    rsync-arg = "-l"
    rsync-arg = "-r"
    rsync-arg = "-O"
    #rsync-arg = "--dry-run"
    rsync-arg = "--stats"
    rsync-arg = "--progress"
    rsync-arg = "--delete"
    rsync-arg = "--no-owner"
    rsync-arg = "--no-group"
    rsync-arg = "--exclude=.git"
    rsync-arg = "--exclude=.git_ignore"

    # Wait between concurrent syncs, to allow for changes to batch together
    cooldown = 1 # second

    desktop-notify = true
    desktop-notify-sound = true

