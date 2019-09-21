package main

import (
	"fmt"
	"os"
)

const Version = "0.2"
const AppName = "HerbGo cli tool"

func Intro() {
	fmt.Printf("%s version %s \n", AppName, Version)
}
func Exec(args ...string) {
	if len(args) > 0 {
		m := Modules.Module(args[0])
		if m != nil {
			m.Exec(args[1:]...)
			return
		}
	}
	help.Exec()
}
func main() {
	errmsg := checkGoPath()
	if errmsg != "" {
		fmt.Println(errmsg)
		os.Exit(1)
	}
	args := os.Args[1:]
	Exec(args...)
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
