package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/cloyop/locker/pkg"
	"github.com/cloyop/locker/storage"
)

func NormalPath() {
	var valid bool
	f := pkg.GetLockerFIle()
	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, f); err != nil {
		log.Fatal(err)
	}
	f.Close()
	metadata := storage.NewMetaData()
	for !valid {
		password := pkg.ScanLine("insert your password:")
		if len(password) < 8 {
			fmt.Println("Invalid password length")
			continue
		}
		bytes, err := pkg.MustUnCipher([]byte(password), buffer.Bytes())
		if err != nil {
			fmt.Println("Invalid password")
			continue
		}
		if err := metadata.Decode(bytes); err != nil {
			log.Fatal(err)
		}
		metadata.SetPassword(password)
		valid = true
	}
	pkg.ClearTerminal()
	if metadata.PinPerm.Using {
		valid = false
		for !valid {
			pin := pkg.ScanLine("need your pin to do any read or write: ")
			if len(pin) < 4 || len(pin) > 16 {
				fmt.Printf("invalid pin length: %v ( 4-16 character)\n", len(pin))
				continue
			}
			metadata.PinPerm.SetPin(pin)
			err := metadata.DecryptInsideData()
			if err != nil {
				fmt.Println("invalid pin")
				continue
			}
			valid = true
			metadata.LastActionDone()
		}
	} else {
		err := metadata.DecryptInsideData()
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := metadata.Data.Decode(&metadata.SafeData); err != nil {
		log.Fatal(err)
	}
	CmdLoop(metadata)
}

func InitPath() {
	slog.Info("initializing...")
	metadata := storage.NewMetaData()
	fmt.Println("Set a first step password minimun 8 characters.")
	firstStepPassword := pkg.ScanLine("insert password: ")
	for len(firstStepPassword) < 8 {
		fmt.Println("Password must be minimun 8 characters")

		firstStepPassword = pkg.ScanLine("insert password: ")
	}
	metadata.SetPassword(firstStepPassword)
	fmt.Println("INFO: The password is used to decrypt the first part of your data. even if it leaks. your data would remain safe by your pin")
	fmt.Println("INFO: The pin is another step to keep your data safe")
	AcceptsOrDecline := pkg.ScanLine("want to set a pin? (Y/n): ")
	fmt.Printf("INFO: Make no changes in  %v minutes will mark u as AFK and will be asked your pin to read or write\n", metadata.PinPerm.RequestForIt)
	AcceptsOrDecline = strings.ToLower(AcceptsOrDecline)
	var enablePin = AcceptsOrDecline != "n"
	if enablePin {
		pin := pkg.ScanLine("insert your pin between 4-16 characters: ")
		for len(pin) < 4 && len(pin) > 16 {
			fmt.Println("invalid pin")
			pin = pkg.ScanLine("insert your pin between 4-16 characters: ")
		}
		metadata.PinPerm.Using = true
		metadata.PinPerm.SetPin(pin)
	}
	metadata.Save()
	metadata.LastActionDone()
	fmt.Println("Successfully Initialized")
	time.Sleep(time.Second * 2)
	CmdLoop(metadata)
}
