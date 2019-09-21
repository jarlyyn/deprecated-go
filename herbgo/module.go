package main

//Module herbgo module interface.
type Module interface {
	Cmd() string
	Name() string
	Description() string
	Help() string
	Exec(args ...string)
}

//ModuleList type registered modules list
type ModuleList []Module

//Register register new moudle
func (l *ModuleList) Register(m Module) {
	*l = append(*l, m)
}

//Module get module form registered modules with gvien name.
func (l *ModuleList) Module(cmd string) Module {
	for _, v := range *l {
		if v.Cmd() == cmd {
			return v
		}
	}
	return nil
}

//Modules registered modules list
var Modules = ModuleList{}

func init() {
	Modules.Register(help)
	Modules.Register(list)
	Modules.Register(newApp)
	Modules.Register(config)
	Modules.Register(constant)
	Modules.Register(database)
	Modules.Register(session)
	Modules.Register(custommodule)
	Modules.Register(model)
	Modules.Register(cache)
	Modules.Register(member)
	Modules.Register(form)
	Modules.Register(router)
	Modules.Register(event)
	Modules.Register(api)
	Modules.Register(query)
	Modules.Register(modeldatasource)
}
