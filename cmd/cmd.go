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
		params := processInput(pkg.ScanLine("Locker do: "))
		if params.err != nil {
			fmt.Println(params.err)
			continue
		}
		if params.exit {
			m.Exit()
			break
		}
		if params.clear {
			pkg.ClearTerminal()
			continue
		}
		if params.save {
			m.OnlySave()
			continue
		}
		if params.ls {
			ListAction(params, m)
			continue
		}
		m.NeedPin()
		m.LastActionDone()
		switch params.action {
		case "rm":
			ActionDelete(params, m)
			continue
		case "get":
			ActionGet(params, m)
			continue
		case "set":
			ActionSet(params, m)
			continue
		}
	}
}
func processInput(input string) *CMDParams {
	p := &CMDParams{}
	input = strings.ToLower(input)
	args := strings.Split(input, " ")
	size := len(args)
	if size == 0 {
		p.err = fmt.Errorf("nothing to do")
		return p
	}
	p.action = args[0]
	if !validAction(&p.action) {
		p.err = fmt.Errorf("first param must be ( get | set | rm | ls | save | clear | exit  ) recieved %v", p.action)
		return p
	}
	if toSave(&p.action) {
		p.save = true
		return p
	}
	if toExit(&p.action) {
		p.exit = true
		return p
	}
	if clearConsole(&p.action) {
		p.clear = true
		return p
	}
	if listData(&p.action) {
		if size > 1 {
			p.name = args[1]
		}
		if size > 2 {
			p.kvs.Key = args[2]
		}
		p.ls = true
		return p
	}
	if size < 2 {
		p.err = fmt.Errorf("missing params")
		return p
	}
	p.name = args[1]

	if p.action == "set" {
		if size <= 2 {
			p.err = fmt.Errorf("missing value to set")
			return p
		}
		if args[size-1] == "-" {
			p.kvs.Key = strings.Join(args[2:size-1], " ")
			return p
		}
		p.kvs.Key = args[2]
		if size > 3 {
			joined := strings.Join(args[3:], " ")
			p.kvs.Value = joined
		}
	}
	return p
}
