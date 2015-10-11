package active_sync

import (
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
	ch       chan syncRequest
	cooldown int
}

type syncRequest struct {
	ProjectName   string
	Local_Folder  string
	Remote_Folder string
	Host          string
	URL           string
	FileName      string
	CoolDown      int
	Args          []string
}

func StartActiveSync(cfg *config.IsoscelsConfig) {

	// Create a pipeline
	p := goauto.NewPipeline("ActiveSync™", goauto.Silent)
	p.OSX = true
	defer p.Stop()

	for project, meta := range cfg.Project {

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

		ch := make(chan syncRequest)
		go syncQueue(ch)

		// Create a workflow
		wf := goauto.NewWorkflow(NewSyncTask(project, meta, ch, meta.CoolDown))

		// Add a file pattern to match
		if err := wf.WatchPattern(meta.Watch_Pattern); err != nil {
			panic(err)
		}

		// Add workflow to pipeline
		p.Add(wf)

		rsyncArgs := append(meta.Rsync_Arg, meta.Local_Folder, fmt.Sprintf("%s:%s", meta.Host, meta.Remote_Folder))

		// Run an initial sync
		if meta.Initial_Sync == true {
			logdim("  ╚═══ Running Initial Sync...", nil)
			ch <- syncRequest{
				ProjectName:   project,
				Local_Folder:  meta.Local_Folder,
				Remote_Folder: meta.Remote_Folder,
				Host:          meta.Host,
				URL:           meta.URL,
				FileName:      "Initial Sync",
				CoolDown:      meta.CoolDown,
				Args:          rsyncArgs,
			}
		}

	}
	p.Start()
}

func NewSyncTask(name string, project *config.Project, ch chan syncRequest, cooldown int) goauto.Tasker {
	return &syncTask{name: name, project: project, ch: ch, cooldown: cooldown}
}

func (gt *syncTask) Run(info *goauto.TaskInfo) (err error) {
	info.Buf.Reset()

	trimedFileName := strings.TrimPrefix(info.Src, gt.project.Local_Folder)
	info.Target = fmt.Sprintf("%s%s", gt.project.Remote_Folder, trimedFileName)
	logdim(fmt.Sprintf("[%s] File modified: %s", gt.name, trimedFileName), nil)

	rsyncArgs := append(gt.project.Rsync_Arg, gt.project.Local_Folder, fmt.Sprintf("%s:%s", gt.project.Host, gt.project.Remote_Folder))

	gt.ch <- syncRequest{
		ProjectName:   gt.name,
		Local_Folder:  gt.project.Local_Folder,
		Remote_Folder: gt.project.Remote_Folder,
		Host:          gt.project.Host,
		URL:           gt.project.URL,
		FileName:      trimedFileName,
		CoolDown:      gt.project.CoolDown,
		Args:          rsyncArgs,
	}

	return nil

}

func syncQueue(ch chan syncRequest) {
	fileCount := 0
	lastSync := time.Now()
	sr := syncRequest{}

	for {

		select {
		case sr = <-ch:
			fileCount++
		default:
			currentCount := fileCount
			now := time.Now()
			td := now.Sub(lastSync).Seconds()
			if td > float64(sr.CoolDown) && fileCount > 0 {
				gocmd := exec.Command("rsync", sr.Args...)

				log(fmt.Sprintf("[%s] Starting rsync...", sr.ProjectName), nil)
				logdim(fmt.Sprintf("  ╚═══ cmd: rsync %s", strings.Join(sr.Args, " ")), nil)

				err := gocmd.Run()

				noteStr := fmt.Sprintf("Trigger: %s", sr.FileName)

				if fileCount > 1 {
					noteStr = fmt.Sprintf("Completed Sync of [%d] trigger(s).", fileCount)
				}

				note := gosxnotifier.NewNotification(noteStr)
				note.Title = sr.ProjectName
				if err == nil {
					note.Subtitle = "File Sync Complete"
					note.Sound = gosxnotifier.Bottle
					note.Link = sr.URL
					note.AppIcon = "images/logo.png"
				} else {
					note.Subtitle = "File Sync Failure!"
					note.Sound = gosxnotifier.Sosumi
					note.AppIcon = "images/logo-failure.png"
					log("runSync", err)
				}

				err = note.Push()
				if err != nil {
					log("Error with Desktop Notification!", err)
				}

				log(fmt.Sprintf("[%s] Completed Sync of [%d] trigger(s).", sr.ProjectName, fileCount), nil)
				lastSync = time.Now()
				fileCount = fileCount - currentCount
			}
		}

	}
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
