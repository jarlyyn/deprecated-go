package main

import (
	"fmt"
	"os"
	"path"
)

const argValueAppTypeCommon = "common"
const argValueAppTypeAPI = "api"
const argValueAppTypeMinimum = "minimum"

const argValueAppTemplateTmpl = "tmpl"
const argValueAppTemplateJet = "jet"

func CheckIsHerbAppFolder(appPath string) bool {
	if !FileExists(path.Join(appPath, "src", "main.go")) {
		return false
	}
	if !MustBeFolder(path.Join(appPath, "config")) {
		return false
	}
	if !MustBeFolder(path.Join(appPath, "resources")) {
		return false
	}
	if !MustBeFolder(path.Join(appPath, "system")) {
		return false
	}
	if !MustBeFolder(path.Join(appPath, "src", "vendor", "modules", "app")) {
		return false
	}
	return true
}

func InAppFolderOrExit() string {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if !(CheckIsHerbAppFolder(d)) {
		fmt.Printf("Current path \"%s\" is not a herb app root folder.\r\n", d)
		os.Exit(2)
	}
	return d
}
func createApp(appPath string) {
	if FileExists(appPath) {
		fmt.Printf("\"%s\" exists.\nCreate app fail.\n", appPath)
		os.Exit(2)
	}
	err := CopyFolder(Resources("skeleton"), appPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("App installed in \"%s\"\n", appPath)
}

func installHTTP(appPath string) (successed bool) {
	task := CopyTasks{}
	templatePath := Resources("http")
	task.Add(
		path.Join(templatePath, "config", "http.toml"),
		path.Join(appPath, "config", "http.toml"),
	)
	task.Add(
		path.Join(templatePath, "config", "http.toml"),
		path.Join(appPath, "system", "config.examples", "http.toml"),
	)
	task.Add(
		path.Join(templatePath, "config", "csrf.toml"),
		path.Join(appPath, "config", "csrf.toml"),
	)
	task.Add(
		path.Join(templatePath, "config", "csrf.toml"),
		path.Join(appPath, "system", "config.examples", "csrf.toml"),
	)
	task.Add(
		path.Join(templatePath, "system", "constants", "assets.toml"),
		path.Join(appPath, "system", "constants", "assets.toml"),
	)
	srcPath := path.Join(appPath, "src")
	SrcSrcPath := path.Join(templatePath, "src")
	task.Add(
		path.Join(SrcSrcPath, "http.go"),
		path.Join(srcPath, "http.go"),
	)
	modulesPath := path.Join(srcPath, "vendor", "modules")
	SrcModulesPath := path.Join(SrcSrcPath, "vendor", "modules")
	task.Add(
		path.Join(SrcModulesPath, "app", "http.go"),
		path.Join(modulesPath, "app", "http.go"),
	)
	task.Add(
		path.Join(SrcModulesPath, "app", "csrf.go"),
		path.Join(modulesPath, "app", "csrf.go"),
	)
	task.Add(
		path.Join(SrcModulesPath, "app", "assets.go"),
		path.Join(modulesPath, "app", "assets.go"),
	)

	task.Add(
		path.Join(SrcModulesPath, "messages", "forms.go"),
		path.Join(modulesPath, "messages", "forms.go"),
	)
	task.Add(
		path.Join(SrcModulesPath, "messages", "messages.go"),
		path.Join(modulesPath, "messages", "messages.go"),
	)
	task.Add(
		path.Join(SrcModulesPath, "middlewares", "middlewares.go"),
		path.Join(modulesPath, "middlewares", "middlewares.go"),
	)
	task.Add(
		path.Join(SrcModulesPath, "middlewares", "csrf.go"),
		path.Join(modulesPath, "middlewares", "csrf.go"),
	)
	task.Add(
		path.Join(SrcModulesPath, "routers", "api.go"),

		path.Join(modulesPath, "routers", "api.go"),
	)
	task.Add(
		path.Join(SrcModulesPath, "routers", "assests.go"),

		path.Join(modulesPath, "routers", "assests.go"),
	)
	task.Add(
		path.Join(SrcModulesPath, "routers", "routers.go"),

		path.Join(modulesPath, "routers", "routers.go"),
	)
	routersRoutersGoPath := path.Join(srcPath, "main.go")
	ReplaceLine(routersRoutersGoPath,
		`//Replace next line "errFuncWhenRunFuncNotRewrited()" with your own app run function`,
		"	//Run app as http server.",
	)
	ReplaceLine(routersRoutersGoPath,
		`errFuncWhenRunFuncNotRewrited()`,
		"	RunHTTP()",
	)

	failed := task.Check()
	if failed != "" {
		fmt.Printf("File \"%s\" exists.\nInstalling http moudles failed.\n", failed)
		return false
	}
	task.Run()
	fmt.Println("Http modules installed.")
	return true
}
func installJetEngine(appPath string) (successed bool) {
	task := CopyTasks{}
	templatePath := Resources("http")
	task.Add(
		path.Join(templatePath, "resources", "template.jet", "views.toml"),
		path.Join(appPath, "resources", "template.jet", "views.toml"),
	)
	task.Add(
		path.Join(templatePath, "resources", "template.jet", "layouts", "main.jet"),
		path.Join(appPath, "resources", "template.jet", "layouts", "main.jet"),
	)
	task.Add(
		path.Join(templatePath, "resources", "template.jet", "views", "index.jet"),
		path.Join(appPath, "resources", "template.jet", "views", "index.jet"),
	)
	srcPath := path.Join(appPath, "src")
	SrcSrcPath := path.Join(templatePath, "src")

	modulesPath := path.Join(srcPath, "vendor", "modules")
	SrcModulesPath := path.Join(SrcSrcPath, "vendor", "modules")
	task.Add(
		path.Join(SrcModulesPath, "views", "jetengine.go"),
		path.Join(modulesPath, "views", "jetengine.go"),
	)
	failed := task.Check()
	if failed != "" {
		fmt.Printf("File \"%s\" exists.\nInstalling jet engine moudles failed.\n", failed)
		return false
	}
	task.Run()
	viewsInitGoPath := path.Join(modulesPath, "views", "init.go")

	ReplaceLine(viewsInitGoPath,
		"var ViewsInitiator func()",
		"var ViewsInitiator = initJetViews",
	)
	routersRoutersGoPath := path.Join(modulesPath, "routers", "routers.go")
	ReplaceLine(routersRoutersGoPath,
		`//Router.GET("/").Use(HTMLMiddlewares()...).HandleFunc(actions.IndexAction)`,
		"	Router.GET(\"/\").\n		Use(HTMLMiddlewares()...).\n		HandleFunc(actions.IndexAction)",
	)
	fmt.Println("Render jet engine installed.")
	return true
}
func installTmplEngine(appPath string) (successed bool) {
	task := CopyTasks{}
	templatePath := Resources("http")
	task.Add(
		path.Join(templatePath, "resources", "template.tmpl", "views.toml"),
		path.Join(appPath, "resources", "template.tmpl", "views.toml"),
	)
	task.Add(
		path.Join(templatePath, "resources", "template.tmpl", "layouts", "main.tmpl"),
		path.Join(appPath, "resources", "template.tmpl", "layouts", "main.tmpl"),
	)
	task.Add(
		path.Join(templatePath, "resources", "template.tmpl", "views", "index.tmpl"),
		path.Join(appPath, "resources", "template.tmpl", "views", "index.tmpl"),
	)
	srcPath := path.Join(appPath, "src")
	SrcSrcPath := path.Join(templatePath, "src")

	modulesPath := path.Join(srcPath, "vendor", "modules")
	SrcModulesPath := path.Join(SrcSrcPath, "vendor", "modules")
	task.Add(
		path.Join(SrcModulesPath, "views", "tmplengine.go"),
		path.Join(modulesPath, "views", "tmplengine.go"),
	)
	failed := task.Check()
	if failed != "" {
		fmt.Printf("File \"%s\" exists.\nInstalling go template engine moudles failed.\n", failed)
		return false
	}
	task.Run()
	viewsInitGoPath := path.Join(modulesPath, "views", "init.go")

	ReplaceLine(viewsInitGoPath,
		"var ViewsInitiator func()",
		"var ViewsInitiator = initTmplViews",
	)
	routersRoutersGoPath := path.Join(modulesPath, "routers", "routers.go")
	ReplaceLine(routersRoutersGoPath,
		`//Router.GET("/").Use(HTMLMiddlewares()...).HandleFunc(actions.ActionWebsiteIndex)`,
		"	Router.GET(\"/\").\n		Use(HTMLMiddlewares()...).\n		HandleFunc(actions.ActionWebsiteIndex)",
	)
	fmt.Println("Render go template engine installed.")
	return true
}
func installWebsite(appPath string) (successed bool) {
	task := CopyTasks{}
	templatePath := Resources("http")
	task.Add(
		path.Join(templatePath, "system", "constants", "website.toml"),
		path.Join(appPath, "system", "constants", "website.toml"),
	)
	task.Add(
		path.Join(templatePath, "resources", "errorpages", "404.html"),
		path.Join(appPath, "resources", "errorpages", "404.html"),
	)
	task.Add(
		path.Join(templatePath, "resources", "errorpages", "500.html"),
		path.Join(appPath, "resources", "errorpages", "500.html"),
	)
	srcPath := path.Join(appPath, "src")
	SrcSrcPath := path.Join(templatePath, "src")

	modulesPath := path.Join(srcPath, "vendor", "modules")
	SrcModulesPath := path.Join(SrcSrcPath, "vendor", "modules")
	task.Add(
		path.Join(SrcModulesPath, "actions", "website.go"),
		path.Join(modulesPath, "actions", "website.go"),
	)
	task.Add(
		path.Join(SrcModulesPath, "middlewares", "errorpages.go"),
		path.Join(modulesPath, "middlewares", "errorpages.go"),
	)
	task.Add(
		path.Join(SrcModulesPath, "app", "website.go"),
		path.Join(modulesPath, "app", "website.go"),
	)
	task.Add(
		path.Join(SrcModulesPath, "routers", "html.go"),
		path.Join(modulesPath, "routers", "html.go"),
	)
	task.Add(
		path.Join(SrcModulesPath, "views", "init.go"),
		path.Join(modulesPath, "views", "init.go"),
	)
	task.Add(
		path.Join(SrcModulesPath, "views", "views.go"),
		path.Join(modulesPath, "views", "views.go"),
	)
	failed := task.Check()
	if failed != "" {
		fmt.Printf("File \"%s\" exists.\nInstalling http moudles failed.\n", failed)
		return false
	}
	task.Run()
	routersRoutersGoPath := path.Join(modulesPath, "routers", "routers.go")

	ReplaceLine(routersRoutersGoPath,
		`//"modules/actions"`,
		`	"modules/actions"`,
	)
	ReplaceLine(routersRoutersGoPath,
		`//var RouterHTML = newHTMLRouter()`,
		"	var RouterHTML = newHTMLRouter()",
	)
	ReplaceLine(routersRoutersGoPath,
		`//Router.StripPrefix("/page").Use(HTMLMiddlewares()...).Handle(RouterHTML)`,
		"	Router.StripPrefix(\"/page\").\r\n		Use(HTMLMiddlewares()...).\r\n		Handle(RouterHTML)",
	)
	fmt.Println("Http modules installed.")
	return true
}

var appTypes = choices{
	newChoice("1", "minimum : minimum app.Usually used in non-http app."),
	newChoice("2", "api: http api app.Created with basic api routers and middlewares."),
	newChoice("3", "common : common website.."),
}
var appTypeValues = map[string]string{
	"1": "minimum",
	"2": "api",
	"3": "common",
}

type newAppModule struct {
	help           string
	apptype        string
	templateEngine string
}

func (m *newAppModule) Init() {
	if !Args.Parsed() {
		Args.StringVar(&m.templateEngine, "template", argValueAppTemplateTmpl,
			`App template engine type,avaliable value:
			tmpl:golang template engine
			jet: jet template engine.
`)
		ParseArgs()
	}

}
func (m *newAppModule) Cmd() string {
	return "new"
}
func (m *newAppModule) Name() string {
	return "New App"
}
func (m *newAppModule) Description() string {
	return "Create new app."
}
func (m *newAppModule) Help() string {
	m.Init()
	return m.help
}

func (m *newAppModule) Exec(args ...string) {
	var defaultType = "2"
	m.Init()
	fmt.Println(Args.Args())
	if Args.NArg() != 1 {
		fmt.Println(m.Help())
		Args.PrintDefaults()
		return
	}
	t := MustChoose("Please choose your app type", appTypes, defaultType)
	m.apptype = appTypeValues[t]
	if m.apptype == "" {
		m.apptype = appTypeValues[defaultType]
	}
	dst := Args.Arg(0)
	var templateEngineInstaller func(string) bool
	switch m.templateEngine {
	case argValueAppTemplateTmpl:
		templateEngineInstaller = installTmplEngine
	case argValueAppTemplateJet:
		templateEngineInstaller = installJetEngine
	default:
		fmt.Println("Unknow template engine type " + m.templateEngine)
		os.Exit(2)
	}
	switch m.apptype {
	case argValueAppTypeMinimum:
		createApp(dst)
		return
	case argValueAppTypeAPI:
		createApp(dst)
		if !installHTTP(dst) {
			os.Exit(2)
		}
		return
	case argValueAppTypeCommon:
		createApp(dst)
		installHTTP(dst)
		installWebsite(dst)
		if !templateEngineInstaller(dst) {
			os.Exit(2)
		}
	default:
		fmt.Println("Unknow app type " + m.apptype)
		os.Exit(2)
	}
}

var newApp = &newAppModule{
	help: `Usage herbgo new <options> [path] .
Create new app in given path.
`}
