package storage

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/cloyop/locker/pkg"
)

type Metadata struct {
	ChangesMade bool
	LastAction  time.Time
	Password    string
	Pin         string
	Data        StoreData
}

func NewMetaData() *Metadata {
	return &Metadata{
		Data: StoreData{},
	}
}

func (m Metadata) Save() error {
	layerOneBuffer := new(bytes.Buffer)
	if err := gob.NewEncoder(layerOneBuffer).Encode(&m.Data); err != nil {
		return err
	}
	SafeStoreData, err := pkg.Cipher([]byte(m.Password+m.Pin), layerOneBuffer.Bytes())
	if err != nil {
		return err
	}
	SafeToWrite, err := pkg.Cipher([]byte(m.Password), SafeStoreData)
	if err != nil {
		return err
	}
	if err := os.WriteFile(os.Getenv("LOCKER_PATH")+"/locker.txt", SafeToWrite, os.ModePerm); err != nil {
		return err
	}
	return nil
}
func (m *Metadata) SetFromJson(fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	var thing map[string]KeyValueStore
	if err := json.NewDecoder(f).Decode(&thing); err != nil {
		return err
	}
	for name, v := range thing {
		name = strings.ToLower(name)
		m.Data[name] = &KeyValueStore{Key: v.Key, Value: v.Value}
	}
	m.ChangesMade = true
	return nil
}
func (m *Metadata) NeedPin() {
	if time.Since(m.LastAction) > time.Minute*5 {
		pin := pkg.ScanLine("u have been AFK. now i need your pin: ")
		for pin != m.Pin {
			pin = pkg.ScanLine("Pin incorrect, insert again: ")
		}
		m.LastAction = time.Now()
	}
}
