package views

import (
	"modules/app"

	"github.com/herb-go/herb/ui/render"
	"github.com/herb-go/herb/ui/render/engines/gotemplate"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
)

var initTmplViews = func() {
	oc := render.NewOptionCommon()
	oc.Engine = gotemplate.Engine
	oc.ViewRoot = util.Resource("/template.tmpl")
	Render.Init(oc)
	config.RegisterLoaderAndWatch(util.Resource("/template.tmpl/views.toml"), func(path string) {
		option := render.ViewsOptionCommon{}
		tomlconfig.MustLoad(path, &option)
		if app.Development.Debug {
			option.DevelopmentMode = true
		}
		Render.MustInitViews(option)
	}).Load()
	// gotemplate.Engine.RegisterFunc("date", dateFormat)

}
