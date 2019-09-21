package main

import (
	"fmt"
	"os"
	"path"

	"github.com/herb-go/herbgo/name"
)

type constantModuel struct {
	help  string
	watch bool
}

func (m *constantModuel) Init() {
	if !Args.Parsed() {
		Args.BoolVar(&m.watch, "watch", false,
			`Whether reload config after file changed.
`)

		ParseArgs()
	}
}
func (m *constantModuel) Cmd() string {
	return "constant"
}
func (m *constantModuel) Name() string {
	return "Sonstant"
}
func (m *constantModuel) Description() string {
	return "Create new constant file and code"
}
func (m *constantModuel) Help() string {
	return m.help
}
func (m *constantModuel) Exec(a ...string) {
	Intro()
	m.Init()
	args := Args.Args()

	if len(args) >= 1 {
		folder := InAppFolderOrExit()
		n := name.MustNewOrExit(args...)
		task := RenderTasks{}
		systemPath := path.Join(folder, "system", "constants", n.LowerWithParentDotSeparated+".toml")
		goFilePath := path.Join(folder, "src", "vendor", "modules", "app", n.LowerWithParentDotSeparated+".go")
		configTmplPath := path.Join(LibPath, "resources", "template", "newconfigfile.tmpl")
		var goFileTmplPath string
		if m.watch {
			goFileTmplPath = path.Join(LibPath, "resources", "template", "newsystemwatchgo.tmpl")
		} else {
			goFileTmplPath = path.Join(LibPath, "resources", "template", "newsystemgo.tmpl")
		}
		task.Add(configTmplPath, systemPath, n)
		task.Add(goFileTmplPath, goFilePath, n)
		failed := task.Check()
		if failed != "" {
			fmt.Printf("File \"%s\" exists.\nInstalling create constant \"%s\"failed.\n", failed, n.Raw)
			os.Exit(2)
		}
		task.Run()
		fmt.Printf("%s config files created.\n", n.Title)

		return
	}
	fmt.Println(m.help)

}

var constant = &constantModuel{
	help: `Usage herbgo constant [name].
Create new constant file and code.
File below will be created:
	system/constants/[name].toml
	src/vendor/modules/app/[name].go
`}
