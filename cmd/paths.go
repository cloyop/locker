package cmd

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/cloyop/locker/pkg"
	"github.com/cloyop/locker/storage"
)

func NormalPath() {
	var valid bool
	LockerBytes, err := os.ReadFile(os.Getenv("LOCKER_PATH") + "/locker.txt")
	if err != nil {
		log.Fatal(err)
	}
	m := storage.NewMetaData()
	var firstLayerBytes []byte
	for !valid {
		password := pkg.ScanLine("insert your password:")
		if len(password) < 8 {
			fmt.Println("Invalid password length")
			continue
		}
		firstLayerBytes, err = pkg.UnCipher([]byte(password), LockerBytes)
		if err != nil {
			fmt.Println("Invalid password")
			continue
		}
		m.Password = password
		valid = true
	}
	pkg.ClearTerminal()
	valid = false
	for !valid {
		pin := pkg.ScanLine("need your pin to do any read or write: ")
		if len(pin) < 4 || len(pin) > 16 {
			fmt.Printf("invalid pin length: %v ( 4-16 character)\n", len(pin))
			continue
		}
		storeDataBytes, err := pkg.UnCipher([]byte(m.Password+pin), firstLayerBytes)
		if err != nil {
			fmt.Println("invalid pin")
			continue
		}
		m.Pin = pin
		valid = true
		buffer := new(bytes.Buffer)
		_, err = buffer.Write(storeDataBytes)
		if err != nil {
			log.Fatal(err)
		}
		if err := gob.NewDecoder(buffer).Decode(&m.Data); err != nil {
			log.Fatal(err)
		}
		m.LastAction = time.Now()
	}
	CmdLoop(m)
}

func InitPath() {
	slog.Info("initializing...")
	m := storage.NewMetaData()
	fmt.Println("Set a first step password minimun 8 characters.")
	firstStepPassword := pkg.ScanLine("insert password: ")
	for len(firstStepPassword) < 8 {
		fmt.Println("Password must be minimun 8 characters")
		firstStepPassword = pkg.ScanLine("insert password: ")
	}
	m.Password = firstStepPassword
	fmt.Println("INFO: The password is used to decrypt the first part of your data. even if it leaks. your data would remain safe by your pin")
	fmt.Printf("INFO: Make no changes in 5 minutes will mark u as AFK and will be asked your pin to read or write\n")
	pin := pkg.ScanLine("insert your pin between 4-16 characters: ")
	for len(pin) < 4 && len(pin) > 16 {
		fmt.Println("invalid pin")
		pin = pkg.ScanLine("insert your pin between 4-16 characters: ")
	}
	m.Pin = pin
	m.Save()
	m.LastAction = time.Now()
	fmt.Println("Successfully Initialized")
	time.Sleep(time.Second * 1)
	CmdLoop(m)
}
