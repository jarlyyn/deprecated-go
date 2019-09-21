package main

import "fmt"

type ListModule string

func (m ListModule) Cmd() string {
	return "list"
}
func (m ListModule) Name() string {
	return "List"
}
func (m ListModule) Description() string {
	return "List all modules"
}
func (m ListModule) Help() string {
	return string(m)
}
func (m ListModule) Exec(args ...string) {
	Intro()
	for _, v := range Modules {
		fmt.Printf("%s : %s\n", v.Cmd(), v.Name())
		fmt.Println("  " + v.Description())
		fmt.Print("\n")
	}
}

var list = ListModule("List all modules.")
