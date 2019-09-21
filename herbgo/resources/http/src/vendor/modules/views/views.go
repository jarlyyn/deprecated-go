package views

import (
	"modules/app"

	"github.com/herb-go/herb/ui/render"
)

func SkinPath() string {
	return app.Assets.URLPrefix + "/skin/"
}
func NewRenderData(title string, additionalRenderData ...render.Data) render.Data {
	data := render.Data{}
	data.Set("Name", app.Website.Name)
	data.Set("SkinPath", SkinPath())
	data.Set("BaseURL", app.HTTP.Config.BaseURL)
	data.Set("MetaKeywords", "")
	data.Set("MetaDescription", app.Website.Description)
	data.Set("Title", title)
	for _, v := range additionalRenderData {
		data.Merge(&v)
	}
	return data
}

var ViewIndex = Render.GetView("index")
