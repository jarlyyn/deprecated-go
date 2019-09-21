package util

import (
	"fmt"
	"log"
	"sort"
)

type Module struct {
	Name     string
	Handler  func()
	Position string
}

var Debug = false

func DebugPrintln(args ...interface{}) {
	if Debug {
		fmt.Println(args...)
	}
}

type modulelist []Module

func (m modulelist) Len() int {
	return len(m)
}
func (m modulelist) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m modulelist) Less(i, j int) bool {
	return m[i].Name < m[j].Name
}

var Modules = modulelist{}

func RegisterModule(Name string, handler func()) Module {
	var position string
	lines := GetStackLines(8, 9)
	if len(lines) == 1 {
		position = fmt.Sprintf("%s\r\n", lines[0])
	}
	m := Module{Name: Name, Handler: handler, Position: position}
	Modules = append(Modules, m)
	return m
}

func InitModulesOrderByName() {
	MustLoadRegisteredFolders()
	sort.Sort(Modules)
	for k := range Modules {
		if Debug || ForceDebug {
			fmt.Println("Herb-go util debug: Init module " + Modules[k].Name)
			if Modules[k].Position != "" {
				fmt.Print(Modules[k].Position)
			}
		}
		Modules[k].Handler()
	}
	if Debug || ForceDebug {
		SetWarning("Util", "Debug mode enabled.")
	}
	if HasWarning() {
		output := "Warning:\r\n"
		for k, v := range warnings {
			output = output + "  " + k + ":\r\n"
			for _, wv := range v {
				output = output + "    " + wv + "\r\n"
			}
		}
		log.Print(output)
	}
}
