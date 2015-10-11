package main

import (
	"fmt"
	"os"

	"github.com/murdinc/cli"
	"github.com/murdinc/isosceles/active_sync"
	"github.com/murdinc/isosceles/config"
)

func main() {

	app := cli.NewApp()
	app.Name = "isosceles"
	app.Usage = "Remote Development Tool Roll"
	app.Version = "1.0"
	app.Commands = []cli.Command{
		{
			Name:        "active-sync",
			ShortName:   "as",
			Example:     "active-sync",
			Description: "Actively syncs all configured local folders to their remote",
			Action: func(c *cli.Context) {
				log("Reading configuration file...", nil)
				cfg, err := config.ReadConfig()
				if err != nil {
					log("Config File", err)
				}

				log("Enabled Projects:", nil)
				cfg.ListEnabledProjects()
				log("Starting Active Sync..", nil)
				active_sync.StartActiveSync(cfg) // blocking
			},
		},
		{
			Name:        "all-projects",
			ShortName:   "ap",
			Example:     "all-projects",
			Description: "List all configured projects",
			Action: func(c *cli.Context) {
				cfg, err := config.ReadConfig()
				if err != nil {
					log("Config File", err)
				}
				cfg.ListAllProjects()
			},
		},
		{
			Name:        "enabled-projects",
			ShortName:   "ep",
			Example:     "enabled-projects",
			Description: "List all enabled projects",
			Action: func(c *cli.Context) {
				cfg, err := config.ReadConfig()
				if err != nil {
					log("Config File", err)
				}
				cfg.ListEnabledProjects()
			},
		},
	}

	app.Run(os.Args)
}

// Log Function
////////////////..........
func log(kind string, err error) {
	if err == nil {
		fmt.Printf("%s\n", kind)
	} else {
		fmt.Printf("[ERROR - %s]: %s\n", kind, err)
	}
}
