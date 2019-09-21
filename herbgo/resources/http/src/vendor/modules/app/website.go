package app

import (
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/commonconfig"
	"github.com/herb-go/util/config/tomlconfig"
)

var Website = commonconfig.WebsiteConfig{}

func init() {
	config.RegisterLoader(util.Constants("/website.toml"), func(configpath string) {
		util.Must(tomlconfig.Load(configpath, &Website))
	})
}
