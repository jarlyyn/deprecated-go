package util

import (
	"fmt"
	"os"
)

//EnvForceDebugMode env field to set force demog mode
const EnvForceDebugMode = "HerbDebug"

//EnvRootPath env field to set root path
const EnvRootPath = "HerbRoot"

//ForceDebug force useing debug mode
var ForceDebug bool

//IgnoreEnv ignore os env settings.
var IgnoreEnv = false

func init() {
	RootPath = os.Getenv(EnvRootPath)
	fmt.Println(RootPath)
	if IgnoreEnv == false && os.Getenv(EnvForceDebugMode) != "" {
		ForceDebug = true
		Debug = true
	}
}
