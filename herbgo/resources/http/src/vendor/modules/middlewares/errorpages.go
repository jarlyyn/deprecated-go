package middlewares

import (
	"modules/views"
	"net/http"

	"github.com/herb-go/util"
)

func initMiddlewareErrorPage() {
	MiddlewareErrorPage.OnStatus(404, func(w http.ResponseWriter, r *http.Request, status int) {
		views.Render.MustHTMLFile(w, util.Resource("errorpages", "404.html"), 404)
	})
	MiddlewareErrorPage.OnError(func(w http.ResponseWriter, r *http.Request, status int) {
		views.Render.MustHTMLFile(w, util.Resource("errorpages", "500.html"), status)
	})
	MiddlewareErrorPage.IgnoreStatus(422)
	if util.Debug {
		MiddlewareErrorPage.IgnoreStatus(500)
	}
}

func init() {
	util.RegisterInitiator(ModuleName, "errorpage", initMiddlewareErrorPage)
}
