package main

import (
	"fmt"
	"os"
	"path"

	"github.com/herb-go/herbgo/name"
)

var installModule = func(name *name.Name, level string) bool {
	folder := InAppFolderOrExit()
	goModuleFilePath := path.Join(folder, "src", "vendor", "modules", name.LowerPath("init.go"))
	goModuleFileTmplPath := path.Join(LibPath, "resources", "template", "module", "module.go.tmpl")
	task := RenderTasks{}
	task.Add(goModuleFileTmplPath, goModuleFilePath, map[string]interface{}{"Name": name, "Level": level})
	failed := task.Check()
	if failed != "" {
		fmt.Printf("File \"%s\" exists.\nInstalling module \"%s\"failed.\n", failed, name.Raw)
		return false
	}
	task.Run()
	fmt.Printf("Module \"%s\" created.\n", name.Title)
	return true
}

type moduleModule struct {
	help  string
	level string
}

func (m *moduleModule) Init() {
	if !Args.Parsed() {
		Args.StringVar(&m.level, "level", "900",
			`Module prefix.All modules are sorted by prefix when loading.
		`)
		ParseArgs()
	}
}
func (m *moduleModule) Cmd() string {
	return "module"
}
func (m *moduleModule) Name() string {
	return "Module"
}
func (m *moduleModule) Description() string {
	return "Create new module and config files."
}
func (m *moduleModule) Help() string {
	return m.help
}
func (m *moduleModule) Exec(a ...string) {
	Intro()
	m.Init()
	args := Args.Args()
	var n *name.Name
	if len(args) == 0 {
		fmt.Println("No  moddatabaseule name given.")
		m.Help()
		os.Exit(2)
	} else {
		n = name.MustNewOrExit(args...)
		if !installModule(n, m.level) {
			os.Exit(2)
		}
	}
}

var custommodule = &moduleModule{
	help: `Usage herbgo module [name].
Create module.
File below will be created:
	src/vendor/modules/<name>/<name>.go
`}
