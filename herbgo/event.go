package main

import (
	"fmt"
	"os"
	"path"

	"github.com/herb-go/herbgo/name"
)

var installEvent = func(n *name.Name) bool {
	folder := InAppFolderOrExit()
	goFilePath := path.Join(folder, "src", "vendor", "modules", "appevents", n.Lower+".go")
	goFileTmplPath := path.Join(LibPath, "resources", "template", "event", "event.go.tmpl")
	task := RenderTasks{}

	task.Add(goFileTmplPath, goFilePath, n)
	failed := task.Check()
	if failed != "" {
		fmt.Printf("File \"%s\" exists.\nInstalling event  \"%s\"failed.\n", failed, n.Raw)
		return false
	}
	task.Run()
	fmt.Printf("Event  \"%s\"  created.\n", n.Title)
	return true
}

type eventModule struct {
	help string
}

func (m *eventModule) Init() {
	if !Args.Parsed() {
		ParseArgs()
	}
}
func (m *eventModule) Cmd() string {
	return "event"
}
func (m *eventModule) Name() string {
	return "Event"
}
func (m *eventModule) Description() string {
	return "Create app event file."
}
func (m *eventModule) Help() string {
	return m.help
}
func (m *eventModule) Exec(a ...string) {
	Intro()
	m.Init()
	args := Args.Args()
	var n *name.Name
	if len(args) == 0 {
		fmt.Println("No event  name given.")
		os.Exit(2)
	}
	n = name.MustNewOrExit(args...)
	if !installEvent(n) {
		os.Exit(2)
	}
}

var event = &eventModule{
	help: `Usage herbgo event <name>.
Create event files.
File below will be created:
	src/vendor/modules/appevents/<name>.go
`}
