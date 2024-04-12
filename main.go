package main

import (
	"log"

	"github.com/cloyop/locker/cmd"
	"github.com/cloyop/locker/pkg"
)

func main() {
	if !pkg.Config() {
		log.Fatal("Missing path please set: \n'export LOCKER_PATH=path/to/lockerDir' \n'export PATH=$PATH:$LOCKER_PATH'\n")
	}
	if pkg.ShouldInit() {
		cmd.InitPath()
	} else {
		cmd.NormalPath()
	}
}
