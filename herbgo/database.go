package main

import (
	"fmt"
	"os"
	"path"

	"github.com/herb-go/herbgo/name"
)

var installDatabase = func(n *name.Name) bool {
	folder := InAppFolderOrExit()
	configPath := path.Join(folder, "config", n.LowerWithParentDotSeparated+".toml")
	configExamplePath := path.Join(folder, "system", "config.examples", n.LowerWithParentDotSeparated+".toml")
	goFilePath := path.Join(folder, "src", "vendor", "modules", "app", n.LowerWithParentDotSeparated+".go")
	goModuleFilePath := path.Join(folder, "src", "vendor", "modules", n.LowerPath(n.Lower+".go"))
	configTmplPath := path.Join(LibPath, "resources", "template", "database", "database.toml.tmpl")
	goFileTmplPath := path.Join(LibPath, "resources", "template", "database", "database.go.tmpl")
	goModuleFileTmplPath := path.Join(LibPath, "resources", "template", "database", "database.modules.go.tmpl")
	task := RenderTasks{}

	task.Add(configTmplPath, configPath, n)
	task.Add(configTmplPath, configExamplePath, n)
	task.Add(goFileTmplPath, goFilePath, n)
	task.Add(goModuleFileTmplPath, goModuleFilePath, n)
	failed := task.Check()
	if failed != "" {
		fmt.Printf("File \"%s\" exists.\nInstalling database module \"%s\"failed.\n", failed, n.Raw)
		return false
	}
	task.Run()
	fmt.Printf("Database module \"%s\" config files created.\n", n.Title)
	return true
}

type databaseModule struct {
	help    string
	parents string
}

func (m *databaseModule) Init() {
	if !Args.Parsed() {
		ParseArgs()
	}
}
func (m *databaseModule) Cmd() string {
	return "database"
}
func (m *databaseModule) Name() string {
	return "Database"
}
func (m *databaseModule) Description() string {
	return "Create database module and config files."
}
func (m *databaseModule) Help() string {
	return m.help
}
func (m *databaseModule) Exec(a ...string) {
	Intro()
	m.Init()
	args := Args.Args()
	var n *name.Name
	if len(args) == 0 {
		fmt.Println("No database module name given.\"database\" is used")
		n = name.MustNewOrExit("database")
	} else {
		n = name.MustNewOrExit(args...)
	}
	if !installDatabase(n) {
		os.Exit(2)
	}
}

var database = &databaseModule{
	help: `Usage herbgo databse <name>.
Create database module and config files.
Default name is "database".
File below will be created:
	config/<name>.toml
	system/confg.examples/<name>.toml
	src/vendor/modules/app/<name>.go
	src/vendor/modules/<name>/<name>.go
`}
