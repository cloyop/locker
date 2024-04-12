package storage

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type pinPerm struct {
	RequestForIt int
	Using        bool
	pin          string
}
type StoreData map[string]KeyValueStore

type KeyValueStore struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Prints
func (s *StoreData) PrintList() {
	var list = "All data:\n\n"
	for key, value := range *s {
		if value.Value != "" {
			list += fmt.Sprintf("[ %v -> %v : %v ]\n", key, value.Key, value.Value)
		} else {
			list += fmt.Sprintf("[ %v -> %v ]\n", key, value.Key)
		}
	}
	fmt.Println(list)
}
func (s *StoreData) PrintListNames() {
	var list = "All data:\n\n"
	for key := range *s {
		list += fmt.Sprintf("%v\n", key)
	}
	fmt.Println(list)
}
func (s *StoreData) PrintInFile() {
	fn := strings.Split(strings.Replace(time.Now().String(), " ", "_", 1), ".")[0]
	f, err := os.OpenFile(fn+".json", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	fmt.Println("File Created")
	if err = json.NewEncoder(f).Encode(s); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Data Writted on", f.Name())
}
func (kvs *KeyValueStore) Print() {
	if kvs.Value != "" {
		fmt.Printf("[ %v : %v ]\n", kvs.Key, kvs.Value)
	} else {
		fmt.Printf("[ %v ]\n", kvs.Key)
	}
}

// Encoding
func (s *StoreData) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(s)
}
func (s *StoreData) Decode(data *[]byte) error {
	b := new(bytes.Buffer)
	if _, err := b.Write(*data); err != nil {
		return err
	}
	return gob.NewDecoder(b).Decode(s)
}

func (p *pinPerm) PinIsValid(pin string) bool {
	return p.pin == pin
}
func (p *pinPerm) SetPin(pin string) {
	p.pin = pin
}
