package cmd

import (
	"fmt"
	"strings"

	"github.com/cloyop/locker/pkg"
	"github.com/cloyop/locker/storage"
)

func actionSet(name string, params []string, m *storage.Metadata) {
	size := len(params)
	kvs := storage.KeyValueStore{}
	if size < 1 {
		fmt.Printf("missing key, do: 'set %v <key> <value>' - value optional \n", name)
		return
	}
	kvs.Key = params[0]
	if name == "-f" {
		m.SetFromJson(kvs.Key)
		return
	}
	if size > 1 {
		kvs.Value = params[1]
		if size > 2 {
			if params[size-1] == "-" {
				kvs.Key = strings.Join(params[:size-1], " ")
				kvs.Value = ""
			} else {
				kvs.Value = strings.Join(params[1:], " ")
			}
		}
	}
	item, found := m.Data[name]
	if found {
		if pkg.ScanLine(fmt.Sprintf("Item %v already exists want to overwrite? (Y/n):", name)) == "n" {
			return
		}
		if kvs.Key == "-" {
			kvs.Key = item.Key
		}
		if kvs.Value == "" {
			kvs.Value = item.Value
		}
	}
	m.Data[name] = kvs
	fmt.Println("Sucessfully Saved item ->", name)
	kvs.Print()
	m.ChangesMade(true)
}
func actionGet(name string, params []string, m *storage.Metadata) {
	size := len(params)

	if len(m.Data) == 0 {
		fmt.Println("No Data")
		return
	}
	if name == "-f" {
		var fn string
		if size > 0 {
			fn = params[0]
		}
		m.Data.PrintInFile(fn)
		return
	}
	kvs, exist := m.Data[name]
	if !exist {
		fmt.Printf("Item %v do not exists\n", name)
		return
	}
	kvs.Print()
}
func actionRemove(name string, m *storage.Metadata) {
	if name == "*" {
		if len(m.Data) == 0 {
			fmt.Println("no data to remove")
			return
		}
		m.Data = storage.StoreData{}
		fmt.Println("all data removed Succesfully")
		m.ChangesMade(true)
		return
	}
	_, exist := m.Data[name]
	if !exist {
		fmt.Printf("Item %v do not exists\n", name)
		return
	}
	delete(m.Data, name)
	fmt.Printf("Item %v removed Successfully\n", name)
	m.ChangesMade(true)
}
func actionList(name string, params []string, m *storage.Metadata) {
	size := len(params)
	var key string
	if size > 0 {
		key = params[0]
	}
	if name == "names" {
		m.Data.PrintListNames()
		return
	}
	if name == "find" {
		matches := storage.StoreData{}
		for name, kvs := range m.Data {
			if strings.Contains(name, key) {
				matches[name] = kvs
			}
		}
		matches.PrintList()
		return
	}
	if name != "" {
		fmt.Printf("dont recognize param '%v'\n", name)
		return
	}
	m.Data.PrintList()
}
