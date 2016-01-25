package active_sync

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/deckarep/gosx-notifier"
	"github.com/dshills/goauto"
	"github.com/murdinc/cli"
	"github.com/murdinc/isosceles/config"
	"github.com/toqueteos/webbrowser"
)

type syncTask struct {
	name     string
	project  *config.Project
	wait     int
	triggers int
}

func StartActiveSync(cfg *config.IsoscelsConfig) {

	for project, meta := range cfg.Project {

		// Create a pipeline
		p := goauto.NewPipeline("ActiveSync™", goauto.Silent)
		p.OSX = true
		defer p.Stop()

		if meta.Enabled != true {
			log(fmt.Sprintf("Skipping Project: [%s], since its disabled...", project), nil)
			continue
		}

		log(fmt.Sprintf("Setting up Project: [%s]", project), nil)
		logdim(fmt.Sprintf("  ╚═══ Recursively Watching Folder: [%s]...", meta.Local_Folder), nil)
		logdim(fmt.Sprintf("  ╚═══ Syncing to Remote Folder: [%s]...", meta.Remote_Folder), nil)
		logdim(fmt.Sprintf("  ╚═══ On Remote Host: [%s]...", meta.Host), nil)
		logdim(fmt.Sprintf("  ╚═══ With Cooldown Period: [%f] second(s).", float64(meta.CoolDown)), nil)

		// Open URL in Browser
		if meta.Open_Browser == true {
			logdim("  ╚═══ Opening Browser...", nil)
			webbrowser.Open(meta.URL)
		}

		// Watch directories recursively, ignoring hidden directories
		if err := p.WatchRecursive(meta.Local_Folder, goauto.IgnoreHidden); err != nil {
			panic(err)
		}

		// Create a workflow
		wf := goauto.NewWorkflow(NewSyncTask(project, meta))

		// Add a file pattern to match
		if err := wf.WatchPattern(meta.Watch_Pattern); err != nil {
			panic(err)
		}

		// Add workflow to pipeline
		p.Add(wf)

		go p.Start()

		// Run an initial sync
		if meta.Initial_Sync == true {
			logdim("  ╚═══ Running Initial Sync...", nil)
			wf.Run(&goauto.TaskInfo{Src: "Initial Sync"})
		}

	}

	for {
		reader := bufio.NewReader(os.Stdin)
		rune, _, _ := reader.ReadRune()
		if rune == 'q' {
			return
		}
	}

}

func NewSyncTask(name string, project *config.Project) goauto.Tasker {
	return &syncTask{name: name, project: project}
}

func (task *syncTask) Run(info *goauto.TaskInfo) (err error) {
	info.Buf.Reset()

	trimedFileName := strings.TrimPrefix(info.Src, task.project.Local_Folder)
	info.Target = fmt.Sprintf("%s%s", task.project.Remote_Folder, trimedFileName)
	logdim(fmt.Sprintf("[%s] File modified: %s", task.name, trimedFileName), nil)

	// If we aren't already waiting for a batch of files, start
	if task.wait < 1 {
		task.wait = task.project.CoolDown
		task.triggers = 1

		go func() {
			for task.wait > 0 {
				time.Sleep(time.Second)
				task.wait--
			}

			rsyncArgs := append(task.project.Rsync_Arg, task.project.Local_Folder, fmt.Sprintf("%s:%s", task.project.Host, task.project.Remote_Folder))

			gocmd := exec.Command("rsync", rsyncArgs...)

			logdim(fmt.Sprintf("[%s] Starting rsync...\n  ╚═══ cmd: rsync %s", task.name, strings.Join(rsyncArgs, " ")), nil)

			err := gocmd.Run()

			noteStr := fmt.Sprintf("Trigger: %s", trimedFileName)

			if task.triggers > 1 {
				noteStr = fmt.Sprintf("Completed Sync of [%d] trigger(s).", task.triggers)
			}

			note := gosxnotifier.NewNotification(noteStr)
			note.Title = task.name
			if err == nil {
				note.Subtitle = "File Sync Complete"

				if task.project.Desktop_Notify_Sound == true {
					note.Sound = gosxnotifier.Bottle
				}
				note.Link = task.project.URL
				note.AppIcon = "images/logo.png"
				log(fmt.Sprintf("[%s] Completed Sync of [%d] trigger(s).", task.name, task.triggers), nil)
			} else {
				note.Subtitle = "File Sync Failure!"

				if task.project.Desktop_Notify_Sound == true {
					note.Sound = gosxnotifier.Sosumi
				}
				note.AppIcon = "images/logo-failure.png"
				log("runSync", err)
				log(fmt.Sprintf("[%s] Failed Sync of [%d] trigger(s).", task.name, task.triggers), nil)
			}

			if task.project.Desktop_Notify == true {
				err = note.Push()
				if err != nil {
					log("Error with Desktop Notification!", err)

				}
			}

		}()

	} else {
		// Just reset our cooldown and increment triggers
		task.wait = task.project.CoolDown
		task.triggers++
	}

	return nil

}

// Log Functions
////////////////..........
func log(kind string, err error) {
	if err == nil {
		fmt.Printf("%s\n", kind)
	} else {
		fmt.Printf("[ERROR - %s]: %s\n", kind, err)
	}
}

func logdim(kind string, err error) {
	if err == nil {
		kind = kind + "\n"
		cli.PrintAnsi(`{{ ansi "dim"}}{{ . }}{{ ansi ""}}`, kind)
	} else {
		detail := err.Error()
		cli.ShowErrorMessage(fmt.Sprintf("ERROR - %s", kind), detail)
		os.Exit(1)
	}
}
