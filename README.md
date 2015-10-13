# isosceles
Remote Development Tool

This is a tool built out of a desire to work locally in an environment set up for remote development. I got tired of the lag of a mounted sshfs volume. It's main and only feature is "active-sync", currently - but I hope to add more utilities to it as needed.

To set up projects, add an .isosceles file to your root folder on OS X. An example is located in the example_config.isosceles file in this repo and at the end of the README. 

**Features:**
* Watches an entire directory recursively for changes that match a specific pattern
* Kicks off an rsync when a trigger is detected.
* Pools triggers that happen during the cooldown period (set in the config), to keep from repeating useless syncs.
* rsync flags are fully customizable in the config
* Desktop Notifications can be enabled, with or without sound, for when triggers are processed.

**Coming Up:**
* tail of remote logs - WIP

**CLI Menu:**

![screenshot1](screenshots/help.png)

**Listing All Projects:**

![screenshot1](screenshots/active-projects.png)

**Active Sync:**

![screenshot1](screenshots/active-sync.png)

**Desktop Notifications:**

![screenshot1](screenshots/desktop-notification.png)

**Example configuration: (goes in ~/.isosceles)**

```

# isosceles config
###############################################################################################

[project "MacGruber1"]

    # Enabled this project in active-sync
    #################################################################
    enabled = true

    # Perform Initial Sync when active-sync is started?
    #################################################################
    initial-sync = true

    # Host and Folders Information
    #################################################################
    host = "host.name.com"
    local-folder = "/Users/USER/isosceles/PROJECT/"
    remote-folder = "/PROJECT/"

    # Watch Pattern for file change triggers
    #################################################################
    watch-pattern = "(.php|.html|.css|.htm|.js)"

    # Rsync Arguments, joined in order
    #################################################################
    rsync-arg = "-l"
    rsync-arg = "-r"
    rsync-arg = "-O"
    rsync-arg = "--dry-run"
    rsync-arg = "--stats"
    rsync-arg = "--progress"
    rsync-arg = "--delete"
    rsync-arg = "--no-owner"
    rsync-arg = "--no-group"
    rsync-arg = "--exclude=.git"
    rsync-arg = "--exclude=.git_ignore"

    # Wait perioed (in seconds) between concurrent syncs, to allow for changes
    # to batch together if there are a lot of files
    #################################################################
    cooldown = 1

    # Open the browser when active-sync is turned on
    #################################################################
    open-browser = true
    url = "http://www.host.name.com"

    # Mac OSX Notifications when sync triggers are executed
    #################################################################
    desktop-notify = true
    desktop-notify-sound = true

    # Enable/Disable log tails - todo
    #################################################################
    tail-error = true
    tail-access = true
    tail-extra = false

    # Log paths
    #################################################################
    error-log = "/var/log/nginx/*.access.log"
    access-log = "/var/log/nginx/*.error.log"
    extra-log = ""

