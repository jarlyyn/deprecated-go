package main

import (
	"fmt"
	"os"
	"path"

	"github.com/herb-go/herbgo/name"
)

type routerModule struct {
	help string
}

func (m *routerModule) Init() {
	if !Args.Parsed() {
		ParseArgs()
	}
}
func (m *routerModule) Cmd() string {
	return "router"
}
func (m *routerModule) Name() string {
	return "Router"
}
func (m *routerModule) Description() string {
	return "Create new router"
}
func (m *routerModule) Help() string {
	m.Init()
	return m.help
}
func (m *routerModule) Exec(a ...string) {
	Intro()
	m.Init()
	args := Args.Args()

	if len(args) >= 1 {
		folder := InAppFolderOrExit()
		n := name.MustNewOrExit(args...)
		task := RenderTasks{}
		goFilePath := path.Join(folder, "src", "vendor", "modules", "routers", n.LowerWithParentDotSeparated+".go")
		goTmplPath := path.Join(LibPath, "resources", "template", "router", "router.go.tmpl")
		task.Add(goTmplPath, goFilePath, n)
		failed := task.Check()
		if failed != "" {
			fmt.Printf("File \"%s\" exists.\n Creating router \"%s\"failed.\n", failed, n.Raw)
			os.Exit(2)
		}
		task.Run()
		fmt.Printf("Router %s  created.\n", n.Title)

		return
	}
	fmt.Println(m.help)
}

var router = &routerModule{
	help: `Usage herbgo router [name].
Create new router.
File below will be created:
	routers/[name].go
`}
