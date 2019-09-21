package main

import (
	"fmt"
	"os"
	"path"

	"github.com/herb-go/herbgo/name"
)

type modelDatasourceModule struct {
	help   string
	output bool
}

func (m *modelDatasourceModule) Init() {
	if !Args.Parsed() {
		Args.BoolVar(&m.output, "output", false,
			`create output datasource instead of model datasource. 
`)
		ParseArgs()
	}
}
func (m *modelDatasourceModule) Cmd() string {
	return "modeldatasource"
}
func (m *modelDatasourceModule) Name() string {
	return "Model datasource"
}
func (m *modelDatasourceModule) Description() string {
	return "Create new model dataseource componet code for model"
}
func (m *modelDatasourceModule) Help() string {
	m.Init()
	return m.help
}
func (m *modelDatasourceModule) Exec(a ...string) {
	Intro()
	m.Init()
	args := Args.Args()

	if len(args) >= 1 {
		folder := InAppFolderOrExit()
		n := name.MustNewOrExit(args...)
		task := RenderTasks{}
		var goFilePath, goTmplPath string
		if !m.output {
			modelPath := path.Join(folder, "src", "vendor", "modules", n.LowerPath("models", n.Lower+".go"))
			if !FileExists(modelPath) {
				fmt.Printf("File \"%s\" doesnt exist.\nInstalling  file \"%s\"failed.\n", modelPath, n.Raw)
				os.Exit(2)
			}
			goFilePath = path.Join(folder, "src", "vendor", "modules", n.LowerPath("models", n.Lower+"modeldatasource.go"))
			goTmplPath = path.Join(LibPath, "resources", "template", "model", "modeldatasource.go.tmpl")
		} else {
			outputPath := path.Join(folder, "src", "vendor", "modules", n.LowerPath("outputs", n.Lower+"output.go"))
			if !FileExists(outputPath) {
				fmt.Printf("File \"%s\" doesnt exist.\nInstalling  file \"%s\"failed.\n", outputPath, n.Raw)
				os.Exit(2)
			}
			goFilePath = path.Join(folder, "src", "vendor", "modules", n.LowerPath("outputs", n.Lower+"outputdatasource.go"))
			goTmplPath = path.Join(LibPath, "resources", "template", "model", "modeloutputdatasource.go.tmpl")

		}
		task.Add(goTmplPath, goFilePath, n)
		failed := task.Check()
		if failed != "" {
			fmt.Printf("File \"%s\" exists.\nInstalling  file \"%s\"failed.\n", failed, n.Raw)
			os.Exit(2)
		}
		task.Run()
		fmt.Printf("%s member data code created.\n", n.Title)

		return
	}
	fmt.Println(m.help)
}

var modeldatasource = &modelDatasourceModule{
	help: `Usage herbgo modeldatasource < -output>[name].
Create new   model datasouroce.
 In non-output mode,file below will be created:
	src/vendor/modules/[parent]/models/[name]modeldatasource.go
 In output mode,file below will be created:
	src/vendor/modules/[parent]/outputs/[name]outputsource.go
`}
