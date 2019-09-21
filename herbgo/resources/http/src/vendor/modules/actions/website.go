package actions

import "net/http"
import "modules/views"

var IndexAction = func(w http.ResponseWriter, r *http.Request) {
	data := views.NewRenderData("Index")
	data["Data"] = "data"
	views.ViewIndex.MustRender(w, data)
}
