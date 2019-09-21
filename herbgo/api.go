package main

import (
	"fmt"
	"os"
	"path"

	"github.com/herb-go/herbgo/name"
)

type apiModule struct {
	help string
}

func (m *apiModule) Init() {
	if !Args.Parsed() {
		ParseArgs()
	}
}
func (m *apiModule) Cmd() string {
	return "api"
}
func (m *apiModule) Name() string {
	return "API"
}
func (m *apiModule) Description() string {
	return "Create new api file and code"
}
func (m *apiModule) Help() string {
	m.Init()
	return m.help
}
func (m *apiModule) Exec(a ...string) {
	Intro()
	m.Init()
	args := Args.Args()

	if len(args) >= 1 {
		folder := InAppFolderOrExit()
		n := name.MustNewOrExit(args...)
		task := RenderTasks{}
		configPath := path.Join(folder, "config", n.LowerWithParentDotSeparated+".toml")
		configExamplePath := path.Join(folder, "system", "config.examples", n.LowerWithParentDotSeparated+".toml")
		goFilePath := path.Join(folder, "src", "vendor", "modules", "app", n.LowerWithParentDotSeparated+".go")
		configTmplPath := path.Join(LibPath, "resources", "template", "api", "configfile.tmpl")
		var goTmplPath string
		goTmplPath = path.Join(LibPath, "resources", "template", "api", "configgo.tmpl")
		task.Add(configTmplPath, configPath, n)
		task.Add(configTmplPath, configExamplePath, n)
		task.Add(goTmplPath, goFilePath, n)
		failed := task.Check()
		if failed != "" {
			fmt.Printf("File \"%s\" exists.\nCreate  api \"%s\"failed.\n", failed, n.Raw)
			os.Exit(2)
		}
		task.Run()
		fmt.Printf("%s api files created.\n", n.Title)

		return
	}
	fmt.Println(m.help)
}

var api = &apiModule{
	help: `Usage herbgo api [name].
Create new api server and client config and go file.
File below will be created:
	config/[name].toml
	system/config.examples/[name].toml
	src/vendor/modules/app/[name].go
`}
