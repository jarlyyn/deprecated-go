package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/herb-go/herb/model/sql/db"
	"github.com/herb-go/herbgo/column"
	_ "github.com/herb-go/herbgo/column/mysqlcolumn"
	_ "github.com/herb-go/herbgo/column/sqlitecolumn"
	"github.com/herb-go/herbgo/name"
	"github.com/herb-go/util/config/tomlconfig"
	_ "github.com/mattn/go-sqlite3"
)

type queryConfirmed struct {
	WithInsert  bool
	WithUpdate  bool
	WithFind    bool
	WithFindAll bool
}

var installQuery = func(conn db.Database, query string, database string, module string, table ...string) bool {
	confirmed := queryConfirmed{}
	n := name.MustNewOrExit(table...)
	if module == "" {
		module = n.LowerWithParent
	}
	mc, err := column.New(conn, database, n.Raw)
	if err != nil {
		panic(err)
	}
	qn, _ := name.MustNew(false, query)
	data := map[string]interface{}{
		"Name":      n,
		"Columns":   mc,
		"Module":    module,
		"Query":     query,
		"QueryName": qn,
		"Confirmed": &confirmed,
	}
	folder := InAppFolderOrExit()
	modulepath := path.Join(folder, "src", "vendor", "modules", module)
	if !FileExists(modulepath) || !FileExists(filepath.Join(modulepath, "models", n.Lower+".go")) {
		fmt.Println("Model " + module + "(" + modulepath + ")" + " not exists.")
		os.Exit(2)
	}
	goModuleFilePath := path.Join(folder, "src", "vendor", "modules", module, "models", n.Lower+"query"+qn.Lower+".go")
	goModuleFileTmplPath := path.Join(LibPath, "resources", "template", "model", "query.go.tmpl")
	goModuleFieldsFilePath := path.Join(folder, "src", "vendor", "modules", module, "models", n.Lower+"queryfields"+qn.Lower+".go")
	goModuleFieldsFileTmplPath := path.Join(LibPath, "resources", "template", "model", "queryfields.go.tmpl")
	goOutputFilePath := path.Join(folder, "src", "vendor", "modules", module, "outputs", n.Lower+qn.Lower+"output.go")
	goOutputFileTmplPath := path.Join(LibPath, "resources", "template", "model", "queryoutput.go.tmpl")
	task := RenderTasks{}
	task.Add(goModuleFileTmplPath, goModuleFilePath, data)
	task.Add(goModuleFieldsFileTmplPath, goModuleFieldsFilePath, data)
	if mc.CanCreate() {
		confirmed.WithInsert = MustConfirm("Do you want to create insert query?")
	}
	if mc.HasPrimayKey() {
		confirmed.WithUpdate = MustConfirm("Do you want to create update query?")
	}
	confirmed.WithFind = MustConfirm("Do you want to create find query?")
	confirmed.WithFindAll = MustConfirm("Do you want to create find all query?")
	if MustConfirm("Do you want to create query output model?") {
		task.Add(goOutputFileTmplPath, goOutputFilePath, data)
	}
	failed := task.Check()
	if failed != "" {
		fmt.Printf("File \"%s\" exists.\nInstalling model query \"%s\"failed.\n", failed, n.Raw)
		return false
	}

	task.Run()
	fmt.Printf("model query \"%s\" config files created.\n", n.Title)
	return true
}

type queryModule struct {
	help     string
	database string
	query    string
	module   string
}

func (m *queryModule) Init() {
	if !Args.Parsed() {
		if !Args.Parsed() {
			Args.StringVar(&m.database, "database", "database",
				`database module name. 
	`)
			Args.StringVar(&m.query, "name", "",
				`query name.
`)
			Args.StringVar(&m.module, "module", "",
				`module name where the model files will be installed to.
	`)
			ParseArgs()
		}
	}
}
func (m *queryModule) Cmd() string {
	return "query"
}
func (m *queryModule) Name() string {
	return "Query"
}
func (m *queryModule) Description() string {
	return "Create model query config files and code."
}
func (m *queryModule) Help() string {
	m.Init()
	return m.help
}
func (m *queryModule) Exec(a ...string) {
	Intro()
	m.Init()
	args := Args.Args()
	if len(args) == 0 {
		fmt.Println("No model table name given.")
		os.Exit(2)
	}
	if m.query == "" {
		fmt.Println("No  query name given.")
		os.Exit(2)
	}
	conn := db.New()
	c := db.Config{}
	tomlconfig.MustLoad("./config/"+m.database+".toml", &c)
	conn.SetDriver(c.Driver)
	c.Apply(conn)

	if !installQuery(conn, m.query, m.database, m.module, args...) {
		os.Exit(2)
	}
}

var query = &queryModule{
	help: `Usage herbgo query -name <name> <modelname>.
Create model query module and config files.
File below will be created:
	src/vendor/modules/<modelname>/query<name>.go
`}
