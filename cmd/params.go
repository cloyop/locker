package cmd

import "github.com/cloyop/locker/storage"

type CMDParams struct {
	action string
	name   string
	kvs    storage.KeyValueStore
	ls     bool
	exit   bool
	clear  bool
	save   bool
	err    error
}

func toSave(c *string) bool {
	return *c == "save"
}
func toExit(c *string) bool {
	return *c == "exit"
}
func clearConsole(c *string) bool {
	return *c == "clear"
}
func listData(c *string) bool {
	return *c == "ls"
}
func validAction(action *string) bool {
	if *action != "get" &&
		*action != "set" &&
		*action != "rm" &&
		*action != "ls" &&
		*action != "exit" &&
		*action != "save" &&
		*action != "clear" {
		return false
	}
	return true
}
