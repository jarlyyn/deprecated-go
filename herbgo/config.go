package main

import (
	"fmt"
	"os"
	"path"

	"github.com/herb-go/herbgo/name"
)

type configModule struct {
	help  string
	watch bool
}

func (m *configModule) Init() {
	if !Args.Parsed() {
		Args.BoolVar(&m.watch, "watch", false,
			`Whether reload config after file changed.
`)
		ParseArgs()
	}
}
func (m *configModule) Cmd() string {
	return "config"
}
func (m *configModule) Name() string {
	return "Config"
}
func (m *configModule) Description() string {
	return "Create new config file and code"
}
func (m *configModule) Help() string {
	m.Init()
	return m.help
}
func (m *configModule) Exec(a ...string) {
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
		configTmplPath := path.Join(LibPath, "resources", "template", "newconfigfile.tmpl")
		var goTmplPath string
		if m.watch {
			goTmplPath = path.Join(LibPath, "resources", "template", "newconfigwatchgo.tmpl")
		} else {
			goTmplPath = path.Join(LibPath, "resources", "template", "newconfiggo.tmpl")
		}
		task.Add(configTmplPath, configPath, n)
		task.Add(configTmplPath, configExamplePath, n)
		task.Add(goTmplPath, goFilePath, n)
		failed := task.Check()
		if failed != "" {
			fmt.Printf("File \"%s\" exists.\nInstalling  config \"%s\"failed.\n", failed, n.Raw)
			os.Exit(2)
		}
		task.Run()
		fmt.Printf("%s config files created.\n", n.Title)

		return
	}
	fmt.Println(m.help)
}

var config = &configModule{
	help: `Usage herbgo config [name].
Create new config file and code.
File below will be created:
	config/[name].toml
	system/config.examples/[name].toml
	src/vendor/modules/app/[name].go
`}
