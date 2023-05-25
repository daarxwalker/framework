package framework

import (
	"os"
	"strings"
)

func Env() string {
	if NoEnv() {
		check(os.Setenv(env, "development"))
	}
	return os.Getenv(env)
}

func NoEnv() bool {
	return len(os.Getenv(env)) == 0
}

func Development() bool {
	return os.Getenv(env) == "development"
}

func Staging() bool {
	return os.Getenv(env) == "staging"
}

func Test() bool {
	return os.Getenv(env) == "test"
}

func Production() bool {
	return os.Getenv(env) == "production"
}

func EnvPath(path string) string {
	if Development() && !strings.Contains(path, "dev/") {
		return "dev/" + path
	}
	return path
}
