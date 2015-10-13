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
	ch       chan syncRequest
	cooldown int
}

type syncRequest struct {
	ProjectName          string
	Local_Folder         string
	Remote_Folder        string
	Host                 string
	URL                  string
	FileName             string
	CoolDown             int
	Args                 []string
	Desktop_Notify       bool
	Desktop_Notify_Sound bool
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
				ProjectName:          project,
				Local_Folder:         meta.Local_Folder,
				Remote_Folder:        meta.Remote_Folder,
				Host:                 meta.Host,
				URL:                  meta.URL,
				FileName:             "Initial Sync",
				CoolDown:             meta.CoolDown,
				Args:                 rsyncArgs,
				Desktop_Notify:       meta.Desktop_Notify,
				Desktop_Notify_Sound: meta.Desktop_Notify_Sound,
			}
		}

		go p.Start()
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		rune, _, _ := reader.ReadRune()
		if rune == 'q' {
			return
		}
	}

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
		ProjectName:          gt.name,
		Local_Folder:         gt.project.Local_Folder,
		Remote_Folder:        gt.project.Remote_Folder,
		Host:                 gt.project.Host,
		URL:                  gt.project.URL,
		FileName:             trimedFileName,
		CoolDown:             gt.project.CoolDown,
		Args:                 rsyncArgs,
		Desktop_Notify:       gt.project.Desktop_Notify,
		Desktop_Notify_Sound: gt.project.Desktop_Notify_Sound,
	}

	return nil

}

func syncQueue(ch chan syncRequest) {
	fileCount := 0
	lastSync := time.Now()
	lastTrigger := lastSync
	sr := syncRequest{}

	for {

		select {
		case sr = <-ch:
			fileCount++
			lastTrigger = time.Now()
		default:
			currentCount := fileCount
			now := time.Now()
			lsTD := now.Sub(lastSync).Seconds()    // last sync time difference
			ltTD := now.Sub(lastTrigger).Seconds() // last trigger time difference

			if lsTD > float64(sr.CoolDown) && fileCount > 0 && ltTD > .5 {

				gocmd := exec.Command("rsync", sr.Args...)

				logdim(fmt.Sprintf("[%s] Starting rsync...\n  ╚═══ cmd: rsync %s", sr.ProjectName, strings.Join(sr.Args, " ")), nil)

				err := gocmd.Run()

				noteStr := fmt.Sprintf("Trigger: %s", sr.FileName)

				if fileCount > 1 {
					noteStr = fmt.Sprintf("Completed Sync of [%d] trigger(s).", fileCount)
				}

				note := gosxnotifier.NewNotification(noteStr)
				note.Title = sr.ProjectName
				if err == nil {
					note.Subtitle = "File Sync Complete"

					if sr.Desktop_Notify_Sound == true {
						note.Sound = gosxnotifier.Bottle
					}
					note.Link = sr.URL
					note.AppIcon = "images/logo.png"
					log(fmt.Sprintf("[%s] Completed Sync of [%d] trigger(s).", sr.ProjectName, fileCount), nil)
				} else {
					note.Subtitle = "File Sync Failure!"

					if sr.Desktop_Notify_Sound == true {
						note.Sound = gosxnotifier.Sosumi
					}
					note.AppIcon = "images/logo-failure.png"
					log("runSync", err)
					log(fmt.Sprintf("[%s] Failed Sync of [%d] trigger(s).", sr.ProjectName, fileCount), nil)
				}

				if sr.Desktop_Notify == true {
					err = note.Push()
					if err != nil {
						log("Error with Desktop Notification!", err)

					}
				}

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
