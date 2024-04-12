package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cloyop/locker/pkg"
	"github.com/cloyop/locker/storage"
)

func ActionSet(params *CMDParams, m *storage.Metadata) {
	if params.name == "-f" {
		f, err := os.Open(params.kvs.Key)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		if pkg.ScanLine("this action can overwrite existing values, want to continue? (Y/n): ") == "n" {
			return
		}
		err = m.SetFromJson(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Success writes")
		m.Data.PrintList()
		m.ChangesMade(true)
		return
	}
	item, exist := m.Data[params.name]
	if exist {
		fmt.Println("the current item exists:")
		item.Print()
		r := pkg.ScanLine("want to rewrite it? (Y/n):")
		if r == "n" {
			return
		}
	}
	if params.kvs.Key == "-" {
		params.kvs.Key = item.Key
	}
	if params.kvs.Value == "-" {
		params.kvs.Value = item.Value
	}
	m.Data[params.name] = params.kvs
	fmt.Printf("%v Writted succesfully \n", params.name)
	params.kvs.Print()
	m.ChangesMade(true)
}
func ActionGet(params *CMDParams, m *storage.Metadata) {
	if params.name == "-f" {
		if len(m.Data) == 0 {
			fmt.Println("No Data")
			return
		}
		m.Data.PrintInFile()
		return
	}
	kvs, exist := m.Data[params.name]
	if !exist {
		fmt.Printf("Item %v do not exists\n", params.name)
		return
	}
	kvs.Print()
}
func ActionDelete(params *CMDParams, m *storage.Metadata) {
	if params.name == "*" {
		if len(m.Data) == 0 {
			fmt.Println("no data to remove")
			return
		}
		m.Data = make(storage.StoreData)
		fmt.Println("all data removed Succesfully")
		m.ChangesMade(true)
		return
	}
	_, exist := m.Data[params.name]
	if !exist {
		fmt.Printf("Item %v do not exists\n", params.name)
		return
	}
	delete(m.Data, params.name)
	fmt.Printf("Item %v removed Successfully\n", params.name)
	m.ChangesMade(true)
}
func ListAction(params *CMDParams, m *storage.Metadata) {
	if len(m.Data) == 0 {
		fmt.Println("No Data")
		return
	}
	if params.name == "names" {
		m.Data.PrintListNames()
		return
	}
	if params.name == "find" {
		matches := storage.StoreData{}
		for name, kvs := range m.Data {
			if strings.Contains(name, params.kvs.Key) {
				matches[name] = kvs
			}
		}
		matches.PrintList()
		return
	}
	m.Data.PrintList()
}
