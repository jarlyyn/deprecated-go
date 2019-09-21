package main

import (
	"io"
	"os"
	"path"
)

var defaultFileMode = os.ModePerm

func CopyFolder(src string, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, info.Mode())
	if err != nil {
		return err
	}

	folder, _ := os.Open(src)

	items, err := folder.Readdir(-1)
	if err != nil {
		return err
	}
	for _, obj := range items {

		if obj.IsDir() {
			err = CopyFolder(path.Join(src, obj.Name()), path.Join(dst, obj.Name()))
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(path.Join(src, obj.Name()), path.Join(dst, obj.Name()))
			if err != nil {
				return err
			}
		}

	}
	return nil

}

func CopyFile(src string, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	dpath := path.Dir(dst)
	_, err = os.Stat(dpath)
	if os.IsNotExist(err) {
		spath := path.Dir(src)
		fileinfo, err := os.Stat(spath)
		if err != nil {
			return err
		}
		err = os.MkdirAll(dpath, fileinfo.Mode())
		if err != nil {
			return err
		}
	}
	df, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer df.Close()

	_, err = io.Copy(df, sf)
	if err != nil {
		return err
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, info.Mode())
	return nil

}

func FileExists(f string) bool {
	_, err := os.Stat(f)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	panic(err)
}

func MustBeFolder(apppath string) bool {
	mode, err := os.Stat(apppath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}
	return mode.IsDir()
}
