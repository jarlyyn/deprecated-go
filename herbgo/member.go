package main

import (
	"fmt"
	"os"
	"path"

	"github.com/herb-go/herbgo/name"
)

type memberConfig struct {
	Name              *name.Name
	InstallSession    bool
	InstallSQLUser    bool
	InstallCache      bool
	DatabaseInstalled bool
}

var installMember = func(c *memberConfig) bool {
	folder := InAppFolderOrExit()
	goModuleFilePath := path.Join(folder, "src", "vendor", "modules", c.Name.LowerPath(c.Name.Lower+".go"))
	goModuleFileTmplPath := path.Join(LibPath, "resources", "template", "member", "member.modules.go.tmpl")
	goMiddlewareFilePath := path.Join(folder, "src", "vendor", "modules", "middlewares", c.Name.LowerWithParentDotSeparated+".go")
	goMiddlewareFileTmplPath := path.Join(LibPath, "resources", "template", "member", "middleware.go.tmpl")
	task := RenderTasks{}

	task.Add(goModuleFileTmplPath, goModuleFilePath, c)
	task.Add(goMiddlewareFileTmplPath, goMiddlewareFilePath, c)
	failed := task.Check()
	if failed != "" {
		fmt.Printf("File \"%s\" exists.\nInstalling member module \"%s\"failed.\n", failed, c.Name.Raw)
		return false
	}
	task.Run()
	fmt.Printf("Member module \"%s\" config files created.\n", c.Name.Title)
	return true
}

type memberModule struct {
	help string
}

func (m *memberModule) Init() {
	if !Args.Parsed() {
		ParseArgs()
	}
}
func (m *memberModule) Cmd() string {
	return "member"
}
func (m *memberModule) Name() string {
	return "Member"
}
func (m *memberModule) Description() string {
	return "Create member module and config files."
}
func (m *memberModule) Help() string {
	return m.help
}
func (m *memberModule) Exec(a ...string) {
	Intro()
	m.Init()
	args := Args.Args()
	var c = &memberConfig{}
	if len(args) == 0 {
		fmt.Println("No member module name given.\"member\" is used")
		c.Name = name.MustNewOrExit("member")
	} else {
		c.Name = name.MustNewOrExit(args...)
	}
	c.InstallSession = MustConfirm("Do you want to install session module?")
	c.InstallCache = MustConfirm("Do you want to add member cache code?Otherwise you have to install member cache manually.")
	c.InstallSQLUser = MustConfirm("Do you want to install sqluser?Otherwise you have to install user modules manually.")
	if c.InstallSession {
		sessionName := name.MustNewOrExit(c.Name.Parents + "/" + c.Name.Lower + "/session")
		if !installSession(sessionName) {
			os.Exit(2)
		}
	}
	if c.InstallCache {
		cacheName := name.MustNewOrExit(c.Name.Parents + "/" + c.Name.Lower + "/cache")
		if !installCache(cacheName) {
			os.Exit(2)
		}
	}
	if c.InstallSQLUser {
		if FileExists(path.Join(InAppFolderOrExit(), "src", "vendor", "modules", "database", "database"+".go")) {
			c.DatabaseInstalled = true
		} else {
			if MustConfirm("Database module not found.\nDo you want to install database module?Otherwise you have to install user modules manually.") {
				if !installDatabase(name.MustNewOrExit("database")) {
					os.Exit(2)
				}
				c.DatabaseInstalled = true
			}
		}
	}
	if !installMember(c) {
		os.Exit(2)
	}
}

var member = &memberModule{
	help: `Usage herbgo member <name>.
Create member module and config files.
Default name is "member".
File below will be created:
	config/<name>.toml
	system/config.examples/<name>.toml
	src/vendor/modules/app/<name>.go
	src/vendor/modules/<name>/<name>.go
`}
