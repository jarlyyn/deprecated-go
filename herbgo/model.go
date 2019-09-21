package main

import (
	"fmt"
	"os"
	"path"

	_ "github.com/go-sql-driver/mysql"
	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herbgo/column"
	_ "github.com/herb-go/herbgo/column/mysqlcolumn"
	_ "github.com/herb-go/herbgo/column/sqlitecolumn"
	"github.com/herb-go/herbgo/name"
	"github.com/herb-go/util/config/tomlconfig"
	_ "github.com/mattn/go-sqlite3"
)

type modelConfirmed struct {
	CreateForm   bool
	CreateOutput bool
	CreateAction bool
	WithCreate   bool
	WithRead     bool
	WithUpdate   bool
	WithDelete   bool
	WithList     bool
	WithPager    bool
}

var installModel = func(conn db.Database, database string, module string, table ...string) bool {
	confirmed := modelConfirmed{}
	n := name.MustNewOrExit(table...)
	if module == "" {
		module = n.LowerWithParent
	}
	mc, err := column.New(conn, database, n.Raw)
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{
		"Name":      n,
		"Columns":   mc,
		"Module":    module,
		"Confirmed": &confirmed,
	}
	folder := InAppFolderOrExit()
	modulepath := path.Join(folder, "src", "vendor", "modules", module)
	if !FileExists(modulepath) {
		installModule(name.MustNewOrExit(module), "900")
	}
	goModuleQueriesPath := path.Join(folder, "src", "vendor", "modules", module, "models", n.Lower+"queries.go")
	goModuleQueriesTmplPath := path.Join(LibPath, "resources", "template", "model", "modelqueries.go.tmpl")
	goModuleFieldsFilePath := path.Join(folder, "src", "vendor", "modules", module, "models", n.Lower+"fields.go")
	goModuleFieldsFileTmplPath := path.Join(LibPath, "resources", "template", "model", "modelfields.go.tmpl")
	goModuleModelFilePath := path.Join(folder, "src", "vendor", "modules", module, "models", n.Lower+".go")
	goModuleModelFileTmplPath := path.Join(LibPath, "resources", "template", "model", "model.go.tmpl")
	task := RenderTasks{}
	task.Add(goModuleQueriesTmplPath, goModuleQueriesPath, data)
	task.Add(goModuleFieldsFileTmplPath, goModuleFieldsFilePath, data)
	task.Add(goModuleModelFileTmplPath, goModuleModelFilePath, data)
	if len(mc.PrimaryKeys) == 1 && mc.PrimaryKeys[0].ColumnType == "string" || mc.PrimaryKeys[0].ColumnType == "int" {
		goModuleFormFilePath := path.Join(folder, "src", "vendor", "modules", module, "forms", n.Lower+"form.go")
		goModuleFormFileTmplPath := path.Join(LibPath, "resources", "template", "model", "modelform.go.tmpl")
		goModuleActionFilePath := path.Join(folder, "src", "vendor", "modules", module, "actions", n.Lower+"action.go")
		goModuleActionFileTmplPath := path.Join(LibPath, "resources", "template", "model", "modelaction.go.tmpl")
		goModuleOutputFilePath := path.Join(folder, "src", "vendor", "modules", module, "outputs", n.Lower+"output.go")
		goModuleoutputFileTmplPath := path.Join(LibPath, "resources", "template", "model", "modeloutput.go.tmpl")
		if MustConfirm("Do you want to install standard \"CRUD\" components?") {
			confirmed.CreateAction = true
			confirmed.CreateForm = true
			confirmed.CreateOutput = true
			confirmed.WithCreate = true
			confirmed.WithRead = true
			confirmed.WithUpdate = true
			confirmed.WithDelete = true
			confirmed.WithList = true
			confirmed.WithPager = true
			task.Add(goModuleFormFileTmplPath, goModuleFormFilePath, data)
			task.Add(goModuleActionFileTmplPath, goModuleActionFilePath, data)
			task.Add(goModuleoutputFileTmplPath, goModuleOutputFilePath, data)
		} else {
			if mc.CanCreate() {
				confirmed.WithCreate = MustConfirm("Do you want to create model \"Create\" component?")
			}
			if mc.HasPrimayKey() {
				confirmed.WithRead = MustConfirm("Do you want to create model \"Read\" component?")
				confirmed.WithUpdate = MustConfirm("Do you want to create model \"Update\" component?")
				confirmed.WithDelete = MustConfirm("Do you want to create model \"Delete\" component?")
			}
			if confirmed.WithList = MustConfirm("Do you want to create model \"List\" component?"); confirmed.WithList {
				confirmed.WithPager = MustConfirm("Do you want to use pager for  \"List\" component?")
			}
			if confirmed.WithCreate || confirmed.WithRead || confirmed.WithUpdate || confirmed.WithDelete || confirmed.WithList {
				confirmed.CreateForm = MustConfirm("Do you want to create model forms?")
				task.Add(goModuleFormFileTmplPath, goModuleFormFilePath, data)
				if confirmed.CreateAction = MustConfirm("Do you want to create model actions?"); confirmed.CreateAction {
					task.Add(goModuleActionFileTmplPath, goModuleActionFilePath, data)
				}
			}
			if confirmed.CreateOutput = MustConfirm("Do you want to create model output class?"); confirmed.CreateOutput {
				task.Add(goModuleoutputFileTmplPath, goModuleOutputFilePath, data)
			}
		}
	}
	failed := task.Check()
	if failed != "" {
		fmt.Printf("File \"%s\" exists.\nInstalling model module \"%s\"failed.\n", failed, n.Raw)
		return false
	}

	task.Run()
	fmt.Printf("model module \"%s\" config files created.\n", n.Title)
	return true
}

type modelModule struct {
	help     string
	database string
	module   string
}

func (m *modelModule) Init() {
	if !Args.Parsed() {
		if !Args.Parsed() {
			Args.StringVar(&m.database, "database", "database",
				`database module name. 
	`)
			Args.StringVar(&m.module, "module", "",
				`module name where the model files will be installed to.
	`)
			ParseArgs()
		}
	}
}
func (m *modelModule) Cmd() string {
	return "model"
}
func (m *modelModule) Name() string {
	return "Model"
}
func (m *modelModule) Description() string {
	return "Create model module and config files."
}
func (m *modelModule) Help() string {
	m.Init()
	return m.help
}
func (m *modelModule) Exec(a ...string) {
	Intro()
	m.Init()
	args := Args.Args()
	if len(args) == 0 {
		fmt.Printf("No model table name given.")
		os.Exit(2)
	}
	conn := db.New()
	c := db.Config{}
	tomlconfig.MustLoad("./config/"+m.database+".toml", &c)
	conn.SetDriver(c.Driver)
	c.Apply(conn)

	if !installModel(conn, m.database, m.module, args...) {
		os.Exit(2)
	}
}

var model = &modelModule{
	help: `Usage herbgo model <name>.
Create model module and config files.
File below will be created:
	src/vendor/modules/<name>/<name>.go
`}
