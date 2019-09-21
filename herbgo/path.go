package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var goPath string

const liburl = "github.com/herb-go/herbgo/"

var LibPath string

func Resources(p ...string) string {
	return path.Join(LibPath, "resources", path.Join(p...))
}
func checkGoPath() string {
	p := os.Getenv("GOPATH")
	goPath = strings.Split(p, ":")[0]
	if goPath == "" {
		return "Go path env is empty.Please set GOPATH env in yonr shell config."
	}
	LibPath = path.Join(goPath, "src", liburl)
	mode, err := os.Stat(LibPath)
	if err != nil && os.IsNotExist(err) {
		return fmt.Sprintf("Folder \"%s\" does not exist.\nYou should use \"go get -u %s\" to install", LibPath, LibPath)
	}
	if !mode.IsDir() {
		return fmt.Sprintf("Path \"%s\" is not a folder.", LibPath)
	}
	return ""
}
