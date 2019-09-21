package main

import (
	"fmt"
)

func Help(m Module) {
	fmt.Println(m.Help())
	Args.PrintDefaults()
}

type HelpModule string

func (m HelpModule) Cmd() string {
	return "help"
}
func (m HelpModule) Name() string {
	return "Help"
}
func (m HelpModule) Description() string {
	return "Display help information"
}
func (m HelpModule) Help() string {
	return string(m)
}
func (m HelpModule) Exec(args ...string) {
	Intro()
	if len(args) == 1 {
		module := Modules.Module(args[0])
		if module == nil {
			fmt.Printf("Command %s not found.\nUse herbgo list to show all modules.\n", args[0])
			return
		}
		Help(module)
		return
	}
	fmt.Println(string(m))
}

var help = HelpModule(`Usage herbgo <command> [<args>].
Command list:
  help         :show herbgo helps.
  list         :list all usable modules.
  new          :create new empty app. 
  <module cmd> :execute module with args.
Use 'herbgo help <command>' or 'herbgo help <module cmd>' to see detailed help.`)
