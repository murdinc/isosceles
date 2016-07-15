package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/deckarep/gosx-notifier"
	"github.com/murdinc/terminal"

	"gopkg.in/gcfg.v1"
)

type IsoscelsConfig struct {
	Project map[string]*Project
}

type Project struct {
	Enabled              bool
	Host                 string
	Local_Folder         string
	Remote_Folder        string
	URL                  string
	Error_Log            string
	Access_Log           string
	Extra_Log            string
	Tail_Error           bool
	Tail_Access          bool
	Tail_Extra           bool
	Open_Browser         bool
	Desktop_Notify       bool
	Desktop_Notify_Sound bool
	CoolDown             int
	Watch_Pattern        string
	Rsync_Arg            []string
	Valid                bool `ini:"-"`
}

// Read in a config file
func ReadConfig() (*IsoscelsConfig, error) {

	currentUser, _ := user.Current()

	configLocation := currentUser.HomeDir + "/.isosceles"
	configFile, _ := ioutil.ReadFile(configLocation)
	configString := string(configFile)

	config := IsoscelsConfig{}

	err := gcfg.ReadStringInto(&config, configString)
	if err != nil {
		return &config, err
	}

	errcnt := 0

Loop:
	for project, conf := range config.Project {
		fmt.Printf("Checking config for project: [%s]\n", project)

		// Add trailing slash for rsync if needed
		if conf.Local_Folder[len(conf.Local_Folder)-1] != os.PathSeparator {
			conf.Local_Folder = fmt.Sprintf("%s%s", conf.Local_Folder, string(os.PathSeparator))
		}

		// Check source folder
		source, err := os.Stat(conf.Local_Folder)
		if err != nil {
			fmt.Printf("  ╚═══ Local folder looks bad: [%s]\n", conf.Local_Folder)
			terminal.ErrorLine(err.Error())
			fmt.Println("")
			errcnt++
			continue Loop
		}

		if source.IsDir() {
			fmt.Printf("  ╚═══ Local folder looks good: [%s]\n", conf.Local_Folder)

		}

		// Add trailing slash for rsync if needed
		if conf.Remote_Folder[len(conf.Remote_Folder)-1] != os.PathSeparator {
			conf.Remote_Folder = fmt.Sprintf("%s%s", conf.Remote_Folder, string(os.PathSeparator))
		}

		fmt.Printf("  ╚═══ Remote folder looks good: [%s]\n", conf.Remote_Folder)

		conf.Valid = true
	}

	DesktopNotification("isosceles - configs loaded", fmt.Sprintf("Loaded %d configs, with %d error(s)", len(config.Project), errcnt))

	return &config, nil
}

func DesktopNotification(title, message string) {
	note := gosxnotifier.NewNotification(message)
	note.Title = title
	note.Sound = gosxnotifier.Bottle
	note.Push()
}

// List all projects
func (c *IsoscelsConfig) ListAllProjects() {
	terminal.PrintAnsi(ProjectsTemplate, c)
}

func (c *IsoscelsConfig) ListEnabledProjects() {
	p := make(map[string]*Project)

	for name, project := range c.Project {
		if project.Enabled == true {
			p[name] = project
		}
	}

	terminal.PrintAnsi(ProjectsTemplate, IsoscelsConfig{Project: p})
}

var ProjectsTemplate = `{{range $name, $project := .Project}}
 {{ if and $project.Enabled $project.Valid }}{{ ansi "fggreen"}}✓ {{else}}{{ ansi "fgred"}}X {{ end }}{{ansi ""}}{{ ansi "underscore"}}{{ ansi "bright" }}{{ ansi "fgwhite"}}[{{ $name }}]{{ ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}                    Host: {{ ansi ""}}{{ ansi "fgcyan"}}{{ $project.Host }}{{ ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}            Local Folder: {{ ansi ""}}{{ ansi "fgcyan"}}{{ $project.Local_Folder }}{{ ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}           Remote Folder: {{ ansi ""}}{{ ansi "fgcyan"}}{{ $project.Remote_Folder }}{{ ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}                     URL: {{ ansi ""}}{{ ansi "fgcyan"}}{{ $project.URL }}{{ ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}          Desktop Notify: {{ if $project.Desktop_Notify }}{{ ansi "fggreen"}}✓ {{else}}{{ ansi "fgred"}}X {{ end }}{{ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}    Desktop Notify Sound: {{ if $project.Desktop_Notify }}{{ ansi "fggreen"}}✓ {{else}}{{ ansi "fgred"}}X {{ end }}{{ansi ""}}


{{ ansi "fgwhite"}}------------------------------------------------------------------------------------------------
{{ ansi ""}}
{{ end }}
`
