package views

import (
	"github.com/herb-go/herb/ui/render"

	"github.com/herb-go/util"
)

//Modulename module name used in initing and debuging.
const Modulename = "200View"

//Render html templete render
var Render = render.New()

var ViewsInitiator func()

func init() {
	util.RegisterModule(Modulename, func() {
		ViewsInitiator()
	})
}
