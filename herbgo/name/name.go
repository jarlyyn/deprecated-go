package name

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

type Name struct {
	Raw                         string
	Parents                     string
	ParentsList                 []string
	Title                       string
	Lower                       string
	Camel                       string
	Pascal                      string
	LowerWithParentDotSeparated string
	LowerWithParent             string
	PascalWithParents           string
}

func (n *Name) LowerPath(filename ...string) string {
	fmt.Println(path.Join(n.ParentsList...))
	return path.Join(path.Join(n.ParentsList...), n.Lower, path.Join(filename...))
}

func fieldsep(r rune) bool {
	return r == ' ' || r == '_' || r == '-'
}

func MustNewOrExitWithoutParents(s ...string) *Name {
	n, r := MustNew(true, s...)
	if r == false {
		os.Exit(2)
	}
	return n
}
func MustNewOrExit(s ...string) *Name {
	n, r := MustNew(true, s...)
	if r == false {
		os.Exit(2)
	}
	return n
}
func MustNew(withparents bool, s ...string) (*Name, bool) {
	all := strings.Join(s, " ")
	if all == "" {
		return &Name{}, true
	}
	list := strings.Split(all, "/")
	plist := list[0 : len(list)-1]
	parentsList := []string{}
	for _, v := range plist {
		if v != "" {
			parentsList = append(parentsList, v)
		}
	}
	r := list[len(list)-1]
	s = strings.FieldsFunc(r, fieldsep)
	var match bool
	var err error
	if withparents {
		match, err = regexp.MatchString("^[a-zA-Z][/0-9a-zA-Z\\s\\_\\-]*$", r)
	} else {
		match, err = regexp.MatchString("^[a-zA-Z][0-9a-zA-Z\\s\\_\\-]*$", r)
	}
	if err != nil {
		panic(err)
	}
	if !match {
		fmt.Printf("Name \"%s\" is not available.\nOnly alphanumeric character (0-9,a-z,A-Z) \"-\"_\"and space are allowed in name.\n", r)
		return nil, false
	}
	n := &Name{
		Raw:         r,
		ParentsList: parentsList,
		Parents:     strings.Join(parentsList, "/"),
	}

	n.Title = strings.ToUpper(n.Raw[0:1]) + n.Raw[1:]
	if len(s) == 0 {
		return n, true
	}
	if len(n.ParentsList) > 0 {
		for _, v := range n.ParentsList {
			if commonInitialisms[strings.ToUpper(v)] {
				n.PascalWithParents = n.PascalWithParents + strings.ToUpper(v)
			} else {
				n.PascalWithParents = n.PascalWithParents + strings.ToUpper(v[0:1]) + v[1:]
			}
		}
	}
	for _, v := range s {
		if commonInitialisms[strings.ToUpper(v)] {
			n.Pascal = n.Pascal + strings.ToUpper(v)
		} else {
			n.Pascal = n.Pascal + strings.ToUpper(v[0:1]) + v[1:]
		}
	}
	n.PascalWithParents = n.PascalWithParents + n.Pascal
	n.Camel = s[0][0:1] + n.Pascal[1:]
	n.Lower = strings.ToLower(n.Camel)
	n.LowerWithParentDotSeparated = n.Lower
	n.LowerWithParent = n.Lower
	if len(n.ParentsList) > 0 {
		n.LowerWithParentDotSeparated = strings.Join(n.ParentsList, ".") + "." + n.LowerWithParentDotSeparated
		n.LowerWithParent = strings.Join(n.ParentsList, "/") + "/" + n.LowerWithParent
	}
	return n, true
}
