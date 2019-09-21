package main

import (
	"modules/app"

	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
)

//Must panic if any error rasied
var Must = util.Must

func loadConfigs() {
	//Uncomment next line to print config loading log .
	//config.Debug = true
	config.Lock.RLock()
	app.LoadConfigs()
	config.Lock.RUnlock()
}
func initModules() {
	util.InitModulesOrderByName()
	//Put Your own init code here.
}

//Main app run func.
var run = func() {
	//Replace next line "errFuncWhenRunFuncNotRewrited()" with your own app run function
	errFuncWhenRunFuncNotRewrited()
}

func main() {
	// Set app root path.
	//Default rootpah is "src/../"
	//You can set os env  "HerbRoot" to overwrite this setting while starting app.
	util.RootPath = ""
	defer util.Recover()
	util.UpdatePaths()
	util.MustChRoot()
	loadConfigs()
	initModules()
	util.RegisterDataFolder() //Auto created appdata folder if not exists
	util.MustLoadRegisteredFolders()
	app.Development.MustNotInitializing()
	run()
}
