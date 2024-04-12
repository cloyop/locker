package storage

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cloyop/locker/pkg"
)

type Metadata struct {
	changesMade bool
	password    string
	PinPerm     pinPerm
	lastAction  time.Time
	SafeData    []byte
	Data        StoreData
}

func NewMetaData() *Metadata {
	return &Metadata{
		Data:    make(StoreData),
		PinPerm: pinPerm{RequestForIt: 10},
	}
}

// encrypt
func (m *Metadata) DecryptInsideData() error {
	byts, err := pkg.MustUnCipher([]byte(m.password+m.PinPerm.pin), m.SafeData)
	if err != nil {
		return err
	}
	m.SafeData = byts
	return nil
}
func (m *Metadata) encryptInsideData() {
	buffer := new(bytes.Buffer)
	if err := m.Data.Encode(buffer); err != nil {
		log.Fatal(err)
	}
	m.SafeData = pkg.MustCipherText([]byte(m.password+m.PinPerm.pin), buffer.Bytes())
	m.Data = StoreData{}
}

// Exit
func (m Metadata) OnlySave() {
	if m.Save() {
		fmt.Println("Changes Save Succesfully")
		m.ChangesMade(false)
	}
}
func (m *Metadata) Save() bool {
	m.encryptInsideData()
	buffer := new(bytes.Buffer)
	if err := m.Encode(buffer); err != nil {
		log.Fatal(err)
	}
	safeBytes := pkg.MustCipherText([]byte(m.password), buffer.Bytes())
	f := pkg.GetLockerFIle()
	defer f.Close()
	if w, err := f.Write(safeBytes); w == 0 || err != nil {
		return false
	}
	return true
}
func (m *Metadata) Exit() {
	if m.changesMade {
		fmt.Println("Saving changes...")
		if m.Save() {
			fmt.Println("Changes Saved Succesfully")
		}
	}
}

// Encoding
func (m *Metadata) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(m)
}
func (m *Metadata) Decode(data []byte) error {
	b := new(bytes.Buffer)
	if _, err := b.Write(data); err != nil {
		return err
	}
	return gob.NewDecoder(b).Decode(m)
}

// Misc
func (m *Metadata) SetFromJson(f *os.File) error {
	var thing map[string]KeyValueStore
	if err := json.NewDecoder(f).Decode(&thing); err != nil {
		return err
	}
	for name, v := range thing {
		name = strings.ToLower(name)
		m.Data[name] = KeyValueStore{Key: v.Key, Value: v.Value}
	}
	return nil
}
func (m *Metadata) NeedPin() {
	if m.PinPerm.Using {
		if time.Since(m.lastAction) > time.Minute*time.Duration(m.PinPerm.RequestForIt) {
			pin := pkg.ScanLine("u have been AFK. now i need your pin: ")
			for !m.PinPerm.PinIsValid(pin) {
				pin = pkg.ScanLine("Pin incorrect, insert again: ")
			}
			m.LastActionDone()
		}
	}
}
func (m *Metadata) SetPassword(pwd string) {
	m.password = pwd
}
func (m *Metadata) LastActionDone() {
	m.lastAction = time.Now()
}
func (m *Metadata) ChangesMade(b bool) {
	m.changesMade = b
}
