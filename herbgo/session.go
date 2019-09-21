package main

import (
	"fmt"
	"os"
	"path"

	"github.com/herb-go/herbgo/name"
)

var installSession = func(n *name.Name) bool {
	folder := InAppFolderOrExit()
	configPath := path.Join(folder, "config", n.LowerWithParentDotSeparated+".toml")
	configExamplePath := path.Join(folder, "system", "config.examples", n.LowerWithParentDotSeparated+".toml")
	goFilePath := path.Join(folder, "src", "vendor", "modules", "app", n.LowerWithParentDotSeparated+".go")
	goModuleFilePath := path.Join(folder, "src", "vendor", "modules", n.LowerPath(n.Lower+".go"))
	configTmplPath := path.Join(LibPath, "resources", "template", "session", "session.toml.tmpl")
	goFileTmplPath := path.Join(LibPath, "resources", "template", "session", "session.go.tmpl")
	goModuleFileTmplPath := path.Join(LibPath, "resources", "template", "session", "session.modules.go.tmpl")
	task := RenderTasks{}

	task.Add(configTmplPath, configPath, n)
	task.Add(configTmplPath, configExamplePath, n)
	task.Add(goFileTmplPath, goFilePath, n)
	task.Add(goModuleFileTmplPath, goModuleFilePath, n)
	failed := task.Check()
	if failed != "" {
		fmt.Printf("File \"%s\" exists.\nInstalling session module \"%s\"failed.\n", failed, n.Raw)
		return false
	}
	task.Run()
	fmt.Printf("Session module \"%s\" config files created.\n", n.Title)
	return true
}

type sessionModule struct {
	help string
}

func (m *sessionModule) Init() {
	if !Args.Parsed() {
		ParseArgs()
	}
}
func (m *sessionModule) Cmd() string {
	return "session"
}
func (m *sessionModule) Name() string {
	return "Session"
}
func (m *sessionModule) Description() string {
	return "Create session module and config files."
}
func (m *sessionModule) Help() string {
	return m.help
}
func (m *sessionModule) Exec(a ...string) {
	Intro()
	m.Init()
	args := Args.Args()
	var n *name.Name
	if len(args) == 0 {
		fmt.Println("No session module name given.\"session\" is used")
		n = name.MustNewOrExit("session")
	} else {
		n = name.MustNewOrExit(args...)
	}
	if !installSession(n) {
		os.Exit(2)
	}
}

var session = &sessionModule{
	help: `Usage herbgo session <name>.
Create session module and config files.
Default name is "session".
File below will be created:
	config/<name>.toml
	system/config.examples/<name>.toml
	src/vendor/modules/app/<name>.go
	src/vendor/modules/<name>/<name>.go
`}
