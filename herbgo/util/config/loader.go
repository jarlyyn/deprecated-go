package config

import (
	"fmt"
	"sync"

	"github.com/herb-go/util"
)

var Debug = false

type Loader struct {
	Path     string
	Loader   func(path string)
	Position string
}

func (l *Loader) Load() {
	if Debug || util.ForceDebug {
		fmt.Printf("Herb-go util debug: Load config \"%s\"", l.Path)
		if l.Position != "" {
			fmt.Print(l.Position)
		}

	}
	l.Loader(l.Path)

}

var registeredLoader = []Loader{}

var Lock sync.RWMutex

func RegisterLoader(path string, loader func(path string)) {
	var position string
	lines := util.GetStackLines(8, 9)
	if len(lines) == 1 {
		position = fmt.Sprintf("%s\r\n", lines[0])
	}
	l := Loader{Path: path, Loader: loader, Position: position}
	registeredLoader = append(registeredLoader, l)
}

type LoaderWatcher struct {
	Loader   *Loader
	Watcher  *WatcherManager
	Callback func(Event)
}

func (l *LoaderWatcher) Load() {
	l.Loader.Load()
}
func RegisterLoaderAndWatch(path string, loader func(path string)) *LoaderWatcher {
	var position string
	lines := util.GetStackLines(8, 9)
	if len(lines) == 1 {
		position = fmt.Sprintf("%s\r\n", lines[0])
	}
	l := Loader{Path: path, Loader: loader, Position: position}
	registeredLoader = append(registeredLoader, l)
	callback := func(event Event) {
		if event.IsWrite() || event.IsCreate() {
			l.Load()
		}
	}
	Watcher.On(path, callback)
	return &LoaderWatcher{Loader: &l, Watcher: Watcher, Callback: callback}
}
func LoadAll() {
	defer Lock.RUnlock()
	Lock.RLock()
	for _, v := range registeredLoader {
		v.Load()
	}
}
