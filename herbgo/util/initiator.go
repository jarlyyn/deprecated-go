package util

import (
	"fmt"
	"sort"
)

var Initiators = map[string]modulelist{}

func RegisterInitiator(module string, Name string, handler func()) {
	var position string
	lines := GetStackLines(8, 9)
	if len(lines) == 1 {
		position = fmt.Sprintf("%s\r\n", lines[0])
	}
	if _, ok := Initiators[module]; ok == false {

		Initiators[module] = []Module{}
	}
	Initiators[module] = append(Initiators[module], Module{Name: Name, Handler: handler, Position: position})
}
func InitOrderByName(module string) {
	var initiators modulelist
	var ok bool
	if initiators, ok = Initiators[module]; ok == false {
		return
	}
	sort.Sort(initiators)
	for k := range initiators {
		if Debug || ForceDebug {
			fmt.Printf("Herb-go util debug: Init module function %s/%s.\r\n", module, initiators[k].Name)
			if initiators[k].Position != "" {
				fmt.Print(initiators[k].Position)
			}
		}
		initiators[k].Handler()
	}
}
