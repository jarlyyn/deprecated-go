package util

import "os"

var DefaultFolderMode = os.FileMode(0700)
var registeredFolders = [][]string{}

func RegisterDataFolder(folder ...string) string {
	registeredFolders = append(registeredFolders, folder)
	return AppData(folder...)
}

func MustLoadRegisteredFolders() {
	for _, v := range registeredFolders {
		folder := AppData(v...)
		_, err := os.Stat(folder)
		if err != nil {
			if os.IsNotExist(err) {

				err = os.MkdirAll(folder, DefaultFolderMode)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
