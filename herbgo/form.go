package main

import (
	"fmt"
	"os"
	"path"

	"github.com/herb-go/herbgo/name"
)

type formConfig struct {
	Name *name.Name
}

var installForm = func(c *formConfig) bool {
	folder := InAppFolderOrExit()
	task := RenderTasks{}

	goFormFilePath := path.Join(folder, "src", "vendor", "modules", c.Name.LowerPath("forms"), c.Name.Lower+"form.go")
	goFormFileTmplPath := path.Join(LibPath, "resources", "template", "form", "form.go.tmpl")
	if MustConfirm("Do will want  to create form validating action?") {
		goActionFilePath := path.Join(folder, "src", "vendor", "modules", c.Name.LowerPath("actions"), c.Name.Lower+"action.go")
		goActionFileTmplPath := path.Join(LibPath, "resources", "template", "form", "action.go.tmpl")
		task.Add(goActionFileTmplPath, goActionFilePath, c)
	}

	task.Add(goFormFileTmplPath, goFormFilePath, c)
	failed := task.Check()
	if failed != "" {
		fmt.Printf("File \"%s\" exists.\nInstalling form  \"%s\"failed.\n", failed, c.Name.Raw)
		return false
	}
	task.Run()
	fmt.Printf("Form file \"%s\"  created.\n", c.Name.Title)
	return true
}

type formModule struct {
	help string
}

func (m *formModule) Init() {
	if !Args.Parsed() {
		ParseArgs()
	}
}
func (m *formModule) Cmd() string {
	return "form"
}
func (m *formModule) Name() string {
	return "Form"
}
func (m *formModule) Description() string {
	return "Create form."
}
func (m *formModule) Help() string {
	m.Init()
	return m.help
}
func (m *formModule) Exec(a ...string) {
	Intro()
	m.Init()
	args := Args.Args()
	var c = &formConfig{}
	if len(args) == 0 {
		fmt.Println("No form  name given")
		os.Exit(2)
	}
	c.Name = name.MustNewOrExit(args...)

	if !installForm(c) {
		os.Exit(2)
	}
}

var form = &formModule{
	help: `Usage herbgo form  <options> <name>.
	Create form file.
`}
