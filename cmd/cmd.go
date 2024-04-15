package cmd

import (
	"fmt"
	"strings"

	"github.com/cloyop/locker/pkg"
	"github.com/cloyop/locker/storage"
)

func CmdLoop(m *storage.Metadata) {
	pkg.ClearTerminal()
	for {
		action, name, params, err := processInput(pkg.ScanLine("Locker do: "))
		if err != nil {
			fmt.Println(err)
			continue
		}
		if action == "exit" {
			m.Exit()
			break
		}
		if action == "clear" {
			pkg.ClearTerminal()
			continue
		}
		if action == "save" {
			if m.OnlySave() {
				fmt.Println("Changes Save Succesfully")
				m.ChangesMade(false)
			}
			continue
		}
		m.NeedPin()
		if action == "ls" {
			actionList(name, params, m)
			continue
		}
		if name == "" {
			fmt.Println("Missing Name: '<action> <name> <key> <value>'")
			continue
		}
		switch action {
		case "rm":
			actionRemove(name, m)
			continue
		case "get":
			actionGet(name, params, m)
			continue
		case "set":
			actionSet(name, params, m)
			continue
		}
	}
}
func processInput(input string) (action, name string, params []string, err error) {
	args := strings.Split(input, " ")
	size := len(args)
	action = args[0]
	if !validAction(&action) {
		err = fmt.Errorf("first param must be ( get | set | rm | ls | save | clear | exit  ) recieved %v", action)
		return
	}
	if size > 1 {
		name = strings.ToLower(args[1])
	}
	if size > 2 {
		params = args[2:]
	}
	return
}
