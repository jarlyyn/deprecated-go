package tomlconfig

import (
	"bytes"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/herb-go/util/config"
)

//Load load toml file and unmarshaler to interface.
//Return any error if rasied
func Load(path string, v interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config.NewError(path, err)
	}
	err = toml.Unmarshal(data, v)
	if err != nil {
		return config.NewError(path, err)
	}
	return nil
}

//MustLoad load toml file and unmarshaler to interface.
//Panic if  any error rasied
func MustLoad(path string, v interface{}) {
	err := Load(path, v)
	if err != nil {
		panic(err)
	}
}

//Save save interface to toml file
//Return any error if rasied
func Save(path string, v interface{}) error {
	buffer := bytes.NewBuffer([]byte{})
	err := toml.NewEncoder(buffer).Encode(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, buffer.Bytes(), 0655)
}

//MustSave save interface to toml file
//Panic if  any error rasied
func MustSave(path string, v interface{}) {
	err := Save(path, v)
	if err != nil {
		panic(err)
	}
}
