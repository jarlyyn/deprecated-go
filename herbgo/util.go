package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

func ReplaceLine(path string, from string, to string) {
	var found bool
	info, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return
		}
		panic(err)
	}
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var line string
	var lines = []string{}
	for err != io.EOF {
		line, err = r.ReadString(10)
		if strings.TrimSpace(line) == from {
			found = true
			if line != "" {
				line = to + "\r\n"
			}
		}
		lines = append(lines, line)
	}
	f.Close()
	if found {
		output := strings.Join(lines, "")
		err := ioutil.WriteFile(path, []byte(output), info.Mode())
		if err != nil {
			panic(err)
		}
		fmt.Printf("File \"%s\" updated.\r\n", path)
	}
}

type CopyTask struct {
	Src string
	Dst string
}

type CopyTasks []CopyTask

func (t *CopyTasks) Add(src string, dst string) {
	*t = append(*t, CopyTask{Src: src, Dst: dst})
}
func (t CopyTasks) Check() (failed string) {
	for _, v := range t {
		if FileExists(v.Dst) {
			return v.Dst
		}
	}
	return ""
}
func (t CopyTasks) Run() {
	var err error
	for _, v := range t {
		err = CopyFile(v.Src, v.Dst)
		if err != nil {
			panic(err)
		}
		fmt.Printf("File \"%s\" created.\n", v.Dst)
	}

}

type RenderTask struct {
	Src  string
	Dst  string
	Data interface{}
}

type RenderTasks []RenderTask

func (t *RenderTasks) Add(src string, dst string, data interface{}) {
	*t = append(*t, RenderTask{Src: src, Dst: dst, Data: data})
}
func (t RenderTasks) Check() (failed string) {
	for _, v := range t {
		if FileExists(v.Dst) {
			return v.Dst
		}
	}
	return ""
}
func (t RenderTasks) Run() {
	for _, v := range t {
		info, err := os.Stat(v.Src)
		if err != nil {
			panic(err)
		}
		dpath := path.Dir(v.Dst)
		_, err = os.Stat(dpath)
		if os.IsNotExist(err) {
			srcinfo, err := os.Stat(path.Dir(v.Src))
			if err != nil {
				panic(err)
			}
			err = os.MkdirAll(dpath, srcinfo.Mode())
			if err != nil {
				panic(err)
			}
		}

		Must(ioutil.WriteFile(v.Dst, MustRender(v.Src, v.Data), info.Mode()))
		fmt.Printf("File \"%s\" created.\n", v.Dst)
	}

}

func MustRender(tmplfile string, data interface{}) []byte {

	template := template.Must(template.ParseFiles(tmplfile))
	w := bytes.NewBuffer([]byte{})
	err := template.Execute(w, data)
	if err != nil {
		panic(err)
	}
	return w.Bytes()
}

func MustConfirm(msg string) bool {
	fmt.Println(msg + "(y/n)")
	var s string
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	if s == "y" || s == "yes" {
		return true
	}
	return false
}

type choice struct {
	key   string
	value string
}

func newChoice(key string, value string) *choice {
	return &choice{
		key:   key,
		value: value,
	}
}

type choices []*choice

func MustChoose(msg string, choices choices, defaultchoice string) string {
	if len(choices) > 0 {
		fmt.Println("Please choose.")
		if defaultchoice != "" {
			fmt.Printf("Default choice is %s .\n", defaultchoice)
		}
		fmt.Println(msg)
		for k := range choices {
			if choices[k].key == defaultchoice {
				fmt.Printf("*%s: %s\n", choices[k].key, choices[k].value)
			} else {
				fmt.Printf("%s: %s\n", choices[k].key, choices[k].value)
			}
		}
		var s string
		_, err := fmt.Scan(&s)
		if err != nil {
			panic(err)
		}
		s = strings.TrimSpace(s)
		s = strings.ToLower(s)
		for k := range choices {
			if strings.TrimSpace(strings.ToLower(choices[k].key)) == s {
				return choices[k].key
			}
		}
	}
	return defaultchoice
}
