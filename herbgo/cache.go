package main

import (
	"fmt"
	"os"
	"path"

	"github.com/herb-go/herbgo/name"
)

var installCache = func(n *name.Name) bool {
	folder := InAppFolderOrExit()
	configPath := path.Join(folder, "config", n.LowerWithParentDotSeparated+".toml")
	configExamplePath := path.Join(folder, "system", "config.examples", n.LowerWithParentDotSeparated+".toml")
	goFilePath := path.Join(folder, "src", "vendor", "modules", "app", n.LowerWithParentDotSeparated+".go")
	goModuleFilePath := path.Join(folder, "src", "vendor", "modules", n.LowerPath(n.Lower+".go"))
	configTmplPath := path.Join(LibPath, "resources", "template", "cache", "cache.toml.tmpl")
	goFileTmplPath := path.Join(LibPath, "resources", "template", "cache", "cache.go.tmpl")
	goModuleFileTmplPath := path.Join(LibPath, "resources", "template", "cache", "cache.modules.go.tmpl")
	task := RenderTasks{}

	task.Add(configTmplPath, configPath, n)
	task.Add(configTmplPath, configExamplePath, n)
	task.Add(goFileTmplPath, goFilePath, n)
	task.Add(goModuleFileTmplPath, goModuleFilePath, n)
	failed := task.Check()
	if failed != "" {
		fmt.Printf("File \"%s\" exists.\nInstalling cache module \"%s\"failed.\n", failed, n.Raw)
		return false
	}
	task.Run()
	fmt.Printf("Cache module \"%s\" config files created.\n", n.Title)
	return true
}

type cacheModule struct {
	help string
}

func (m *cacheModule) Init() {
	if !Args.Parsed() {
		ParseArgs()
	}
}
func (m *cacheModule) Cmd() string {
	return "cache"
}
func (m *cacheModule) Name() string {
	return "Cache"
}
func (m *cacheModule) Description() string {
	return "Create cache module and config files."
}
func (m *cacheModule) Help() string {
	return m.help
}
func (m *cacheModule) Exec(a ...string) {
	Intro()
	m.Init()
	args := Args.Args()
	var n *name.Name
	if len(args) == 0 {
		fmt.Println("No cache module name given.\"cache\" is used")
		n = name.MustNewOrExit("cache")
	} else {
		n = name.MustNewOrExit(args...)
	}
	if !installCache(n) {
		os.Exit(2)
	}
}

var cache = &cacheModule{
	help: `Usage herbgo cache <name>.
Create cache module and config files.
Default name is "cache".
File below will be created:
	config/<name>.toml
	system/config.examples/<name>.toml
	src/vendor/modules/app/<name>.go
	src/vendor/modules/<name>/<name>.go
`}
