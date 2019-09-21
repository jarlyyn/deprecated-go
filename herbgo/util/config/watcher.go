package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/herb-go/util"
)

// Op describes a set of file operations.
type Op uint32

var (
	Create = Op(fsnotify.Create)
	Write  = Op(fsnotify.Write)
	Remove = Op(fsnotify.Remove)
	Rename = Op(fsnotify.Rename)
	Chmod  = Op(fsnotify.Chmod)
)

type Event struct {
	e *fsnotify.Event
}

func (e *Event) Path() string {
	return e.e.Name
}

func (e *Event) IsCreate() bool {
	return e.e.Op&fsnotify.Create == fsnotify.Create
}

func (e *Event) IsWrite() bool {
	return e.e.Op&fsnotify.Write == fsnotify.Write
}

func (e *Event) IsRemove() bool {
	return e.e.Op&fsnotify.Remove == fsnotify.Remove
}

func (e *Event) IsReName() bool {
	return e.e.Op&fsnotify.Rename == fsnotify.Rename
}

func (e *Event) IsChmod() bool {
	return e.e.Op&fsnotify.Chmod == fsnotify.Chmod
}

type WatcherManager struct {
	*fsnotify.Watcher
	registeredFuncs map[string][]func(event Event)
}

func (w *WatcherManager) On(path string, callback func(event Event)) {
	if w.registeredFuncs[path] == nil {
		w.registeredFuncs[path] = []func(event Event){callback}
	} else {
		w.registeredFuncs[path] = append(w.registeredFuncs[path], callback)
	}
	w.Add(path)
}
func (w *WatcherManager) Start() {
	go func() {
		for {
			select {
			case event := <-w.Watcher.Events:
				fns := w.registeredFuncs[event.Name]
				for _, k := range fns {
					defer util.Recover()
					k(Event{&event})
				}
			case err := <-w.Watcher.Errors:
				util.LogError(err)
			}
		}
	}()

}

func NewWatcherManager() (*WatcherManager, error) {
	watecher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil
	}
	w := &WatcherManager{
		Watcher:         watecher,
		registeredFuncs: map[string][]func(event Event){},
	}
	w.Start()
	return w, nil
}

var Watcher *WatcherManager

func init() {
	var err error
	Watcher, err = NewWatcherManager()
	if err != nil {
		panic(err)
	}
}
