package framework

import "os"

func root() string {
	dir, err := os.Getwd()
	check(err)
	return dir
}
