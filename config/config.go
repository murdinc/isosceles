package config

import (
	"fmt"
	"io/ioutil"
	"os/user"

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
	Initial_Sync         bool
	CoolDown             int
	Watch_Pattern        string
	Rsync_Arg            []string
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

	return &config, nil
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
 {{ if $project.Enabled }}{{ ansi "fggreen"}}✓ {{else}}{{ ansi "fgred"}}X {{ end }}{{ansi ""}}{{ ansi "underscore"}}{{ ansi "bright" }}{{ ansi "fgwhite"}}[{{ $name }}]{{ ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}                    Host: {{ ansi ""}}{{ ansi "fgcyan"}}{{ $project.Host }}{{ ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}            Local Folder: {{ ansi ""}}{{ ansi "fgcyan"}}{{ $project.Local_Folder }}{{ ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}           Remote Folder: {{ ansi ""}}{{ ansi "fgcyan"}}{{ $project.Remote_Folder }}{{ ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}                     URL: {{ ansi ""}}{{ ansi "fgcyan"}}{{ $project.URL }}{{ ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}            Initial Sync: {{ if $project.Initial_Sync }}{{ ansi "fggreen"}}✓ {{else}}{{ ansi "fgred"}}X {{ end }}{{ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}            Open Browser: {{ if $project.Open_Browser }}{{ ansi "fggreen"}}✓ {{else}}{{ ansi "fgred"}}X {{ end }}{{ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}          Desktop Notify: {{ if $project.Desktop_Notify }}{{ ansi "fggreen"}}✓ {{else}}{{ ansi "fgred"}}X {{ end }}{{ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}    Desktop Notify Sound: {{ if $project.Desktop_Notify }}{{ ansi "fggreen"}}✓ {{else}}{{ ansi "fgred"}}X {{ end }}{{ansi ""}}

    {{ ansi "bright"}}{{ ansi "fgwhite"}}               Error Log: {{ ansi ""}}{{ ansi "fgcyan"}}{{ $project.Error_Log }}{{ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}              Access Log: {{ ansi ""}}{{ ansi "fgcyan"}}{{ $project.Access_Log }}{{ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}               Extra Log: {{ ansi ""}}{{ ansi "fgcyan"}}{{ $project.Extra_Log }}{{ansi ""}}

    {{ ansi "bright"}}{{ ansi "fgwhite"}}          Tail Error Log: {{ if $project.Tail_Error }}{{ ansi "fggreen"}}✓ {{else}}{{ ansi "fgred"}}X {{ end }}{{ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}         Tail Access Log: {{ if $project.Tail_Access }}{{ ansi "fggreen"}}✓ {{else}}{{ ansi "fgred"}}X {{ end }}{{ansi ""}}
    {{ ansi "bright"}}{{ ansi "fgwhite"}}          Tail Extra Log: {{ if $project.Tail_Extra }}{{ ansi "fggreen"}}✓ {{else}}{{ ansi "fgred"}}X {{ end }}{{ansi ""}}

{{ ansi "fgwhite"}}------------------------------------------------------------------------------------------------
{{ ansi ""}}
{{ end }}
`

// Log Function
////////////////..........
func log(kind string, err error) {
	if err == nil {
		fmt.Printf("%s\n", kind)
	} else {
		detail := err.Error()
		terminal.ShowErrorMessage(fmt.Sprintf("ERROR - %s", kind), detail)
	}
}
