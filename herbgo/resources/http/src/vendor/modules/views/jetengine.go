package views

import (
	"modules/app"

	"github.com/herb-go/herb/ui/render"
	"github.com/herb-go/herb/ui/render/engines/jet"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
)

var initJetViews = func() {
	oc := render.NewOptionCommon()
	oc.Engine = jet.Engine
	oc.ViewRoot = util.Resource("/template.jet")
	Render.Init(oc)
	config.RegisterLoaderAndWatch(util.Resource("/template.jet/views.toml"), func(path string) {
		option := render.ViewsOptionCommon{}
		tomlconfig.MustLoad(path, &option)
		if app.Development.Debug {
			option.DevelopmentMode = true
		}
		Render.MustInitViews(option)
	}).Load()
	// jet.Engine.RegisterFunc("date", dateFormat)

}
