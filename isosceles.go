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
	app.Usage = "Remote Development Tool"
	app.Version = "1.1"
	app.Commands = []cli.Command{
		{
			Name:      "active-sync",
			ShortName: "as",
			Usage:     "Actively syncs all configured local folders to their remote",
			Action: func(c *cli.Context) error {
				log("Reading configuration file...", nil)
				cfg, err := config.ReadConfig()
				if err != nil {
					log("Config File", err)
				}

				log("Enabled Projects:", nil)
				cfg.ListEnabledProjects()
				log("Starting Active Sync... press q + return to exit", nil)
				active_sync.StartActiveSync(cfg) // blocking
				return nil
			},
		},
		{
			Name:      "all-projects",
			ShortName: "ap",
			Usage:     "List all configured projects",
			Action: func(c *cli.Context) error {
				cfg, err := config.ReadConfig()
				if err != nil {
					log("Config File", err)
				}
				cfg.ListAllProjects()
				return nil
			},
		},
		{
			Name:      "enabled-projects",
			ShortName: "ep",
			Usage:     "List all enabled projects",
			Action: func(c *cli.Context) error {
				cfg, err := config.ReadConfig()
				if err != nil {
					log("Config File", err)
				}
				cfg.ListEnabledProjects()
				return nil
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
