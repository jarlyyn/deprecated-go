package util

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
func mustPath(path string, err error) string {
	if err != nil {
		panic(err)
	}
	return path
}

var RootPath string
var ResouresPath string
var AppDataPath string
var ConfigPath string
var SystemPath string
var ConstantsPath string
var UpdatePaths = func() error {
	if RootPath == "" {
		RootPath = filepath.Join(filepath.Dir(mustPath(os.Executable())), "../")
	}
	ResouresPath = path.Join(RootPath, "resources")
	AppDataPath = path.Join(RootPath, "appdata")
	ConfigPath = path.Join(RootPath, "config")
	SystemPath = path.Join(RootPath, "system")
	ConstantsPath = path.Join(RootPath, "system", "constants")
	return nil
}

var MustChRoot = func() {
	Must(os.Chdir(RootPath))
}

func SetConfigPath(paths ...string) {
	ConfigPath = path.Join(paths...)
}
func MustGetWD() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}

func joinPath(p string, filepath ...string) string {
	return path.Join(p, path.Join(filepath...))
}
func Resource(filepaths ...string) string {
	return joinPath(ResouresPath, filepaths...)
}
func Config(filepaths ...string) string {
	return joinPath(ConfigPath, filepaths...)
}
func AppData(filepaths ...string) string {
	return joinPath(AppDataPath, filepaths...)
}
func System(filepaths ...string) string {
	return joinPath(SystemPath, filepaths...)
}
func Constants(filepaths ...string) string {
	return joinPath(ConstantsPath, filepaths...)
}

var QuitChan = make(chan int)
var SignalChan = make(chan os.Signal)
var LeaveMessage = "Bye."

func WaitingQuit() {
	signal.Notify(SignalChan, os.Interrupt, os.Kill)
	select {
	case <-SignalChan:
		close(QuitChan)
	case <-QuitChan:
	}
	fmt.Println("Quiting ...")
}
func Bye() {
	if LeaveMessage != "" {
		fmt.Println(LeaveMessage)
	}
}
func Quit() {
	defer func() {
		recover()
	}()
	QuitChan <- 1
}

var LoggerMaxLength = 5
var LoggerIgnoredErrorsChecker = []func(error) bool{
	func(err error) bool {
		oe, ok := err.(*net.OpError)
		if ok {
			if oe.Err == syscall.EPIPE || oe.Err == syscall.ECONNRESET {
				return true
			}
			se, ok := oe.Err.(*os.SyscallError)
			if ok && (se.Err == syscall.EPIPE || se.Err == syscall.ECONNRESET) {
				return true
			}
		}
		return false
	},
}

var IsErrorIgnored = func(err error) bool {
	for k := range LoggerIgnoredErrorsChecker {
		if LoggerIgnoredErrorsChecker[k](err) {
			return true
		}
	}
	return false
}
var RegisterLoggerIgnoredErrorsChecker = func(f func(error) bool) {
	LoggerIgnoredErrorsChecker = append(LoggerIgnoredErrorsChecker, f)
}

func init() {
	err := UpdatePaths()
	if err != nil {
		panic(err)
	}
}
