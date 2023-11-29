package helper

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ariefsn/go-resik/logger"
	"github.com/joho/godotenv"
)

type envApp struct {
	Name string
	Host string
	Port string
}

type env struct {
	App  envApp
	Mode string
}

func (e *env) IsDebug() bool {
	return strings.ToLower(e.Mode) == "debug"
}

type envValue struct {
	value    string
	fallback interface{}
}

func (e envValue) String() string {
	if e.value == "" && e.fallback != nil {
		return fmt.Sprintf("%s", e.fallback)
	}
	return e.value
}

func (e envValue) Int() int {
	if e.value == "" && e.fallback != nil {
		return e.fallback.(int)
	}
	v, err := strconv.Atoi(e.value)
	if err != nil {
		logger.Error(err)
	}
	return v
}

func (e envValue) Bool() bool {
	v, err := strconv.ParseBool(e.value)
	if err != nil {
		logger.Error(err)
	}
	return v
}

var _env *env

func fromEnv(key string, fallback ...interface{}) envValue {
	var fb interface{}
	if len(fallback) > 0 {
		fb = fallback[0]
	}
	return envValue{
		value:    os.Getenv(key),
		fallback: fb,
	}
}

func InitEnv(envFile ...string) {
	err := godotenv.Load(envFile...)
	if err != nil {
		logger.Warning(err.Error())
	}
	_env = &env{
		App: envApp{
			Name: fromEnv("APP_NAME", "RESIK ARCH").String(),
			Host: fromEnv("APP_HOST", "0.0.0.0").String(),
			Port: fromEnv("APP_PORT", "6001").String(),
		},
		Mode: fromEnv("MODE", "DEBUG").String(),
	}
}

func Env() *env {
	if _env == nil {
		InitEnv()
	}

	return _env
}
