package util

var warnings = map[string][]string{}

func SetWarning(module string, info ...string) {
	if len(info) == 0 {
		delete(warnings, module)
		return
	}
	warnings[module] = info
}

func Warnings() map[string][]string {
	return warnings
}

func DelWarning(module string) {
	SetWarning(module)
}

func HasWarning() bool {
	return len(warnings) > 0
}
