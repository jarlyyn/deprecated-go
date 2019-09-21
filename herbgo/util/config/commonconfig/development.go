package commonconfig

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

//ErrAppIsInInitializingMode error rasied if app is in initializing mode
var ErrAppIsInInitializingMode = errors.New("app is in initializing mode.you should turn 'Initializing' in 'config/development.toml' to false")

//ErrAppInitialized error raised if app is initialized
var ErrAppInitialized = errors.New("app is Initialized,should not be initialized again")

//ErrInitializingEnvIsSet err raied if initializing env is set.
//Use ErrInitializingEnvIsSet("envname") to create new error.
type ErrInitializingEnvIsSet string

func (e ErrInitializingEnvIsSet) Error() string {
	return fmt.Sprintf("Initializing Env %s is set", string(e))
}

//DevelopmentConfig program development staus config
type DevelopmentConfig struct {
	//Debug if app is in debug mod.
	Debug bool
	//Testing if app is in testing mode
	Testing bool
	//Profiling if app is in profiling mode
	Profiling bool
	//Benchmarking if app is in benchmarking mode
	Benchmarking bool
	//Initializing if app is Initializing.
	//App should panic before serveing if Initializing setted to True.
	Initializing         bool
	initializers         []*initializer
	initializingEnvs     map[string]bool
	usedInitializingEnvs map[string]string
}

//MustNotInitializing panic if Initializing is true
func (c *DevelopmentConfig) MustNotInitializing() {
	c.loadEnvs()
	if c.Initializing {
		fmt.Println("App is in initializing mode.")
		c.checkInitializers()
		if len(c.initializingEnvs) > 0 {
			fmt.Println("Avaliable initializing envs is listed below:")
			for k := range c.initializingEnvs {
				fmt.Print("  ")
				fmt.Println(k)
			}
		}
		panic(ErrAppIsInInitializingMode)
	}
	if len(c.usedInitializingEnvs) > 0 {
		used := []string{}
		for k := range c.usedInitializingEnvs {
			used = append(used, k)
		}
		panic(ErrInitializingEnvIsSet(strings.Join(used, ",")))
	}
}
func (c *DevelopmentConfig) addInitializer(i *initializer) {
	if c.initializers == nil {
		c.initializers = []*initializer{}
	}
	c.initializers = append(c.initializers, i)
}
func (c *DevelopmentConfig) loadEnvs() {
	c.usedInitializingEnvs = map[string]string{}
	for k := range c.initializingEnvs {
		env := os.Getenv(k)
		if env != "" {
			c.usedInitializingEnvs[k] = env
		}
	}
}

func (c *DevelopmentConfig) checkInitializers() {
	for _, v := range c.initializers {
		for _, env := range v.envs {
			if c.usedInitializingEnvs[env] != "" {
				if v.handler() {
					c.PanicInitialized()
				}
				break
			}
		}
	}
}

// OnEnv declare env name list which will be used by initializer.
func (c *DevelopmentConfig) OnEnv(envs ...string) *RegisteredEnvs {
	return &RegisteredEnvs{
		envs:   envs,
		config: c,
	}
}

// GetInitializeEnv get registered env by name.
func (c *DevelopmentConfig) GetInitializeEnv(name string) string {
	return c.usedInitializingEnvs[name]
}

//RegisteredEnvs registered env name list
type RegisteredEnvs struct {
	envs   []string
	config *DevelopmentConfig
}

//ThenInitalize set initalizer to registered env list.
func (e *RegisteredEnvs) ThenInitalize(handler func() bool) []string {
	i := &initializer{
		envs:    e.envs,
		handler: handler,
	}
	e.config.addInitializer(i)
	if e.config.initializingEnvs == nil {
		e.config.initializingEnvs = map[string]bool{}
	}
	for _, v := range e.envs {
		e.config.initializingEnvs[v] = true
	}
	return e.envs
}

type initializer struct {
	envs    []string
	handler func() bool
}

//PanicInitialized panic a app initialized error.
func (c *DevelopmentConfig) PanicInitialized() {
	panic(ErrAppInitialized)
}
