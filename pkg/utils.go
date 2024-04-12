package pkg

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func ClearTerminal() {
	rn := runtime.GOOS
	var clearCMD string
	if rn == "windows" {
		clearCMD = "cls"
	} else {
		clearCMD = "clear"
	}
	cmd := exec.Command(clearCMD)
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func ScanLine(show string) string {
	fmt.Print(show)
	br := bufio.NewScanner(os.Stdin)
	if !br.Scan() {
		return ""
	}
	password := br.Text()
	return password
}
func ShouldInit() bool {
	fileStats, err := os.Stat(os.Getenv("LOCKER_PATH") + "/locker.txt")
	if err != nil {
		slog.Info("couldnt find file called 'locker.txt' with the metadata. will init\n")
		return true
	}
	if fileStats.IsDir() {
		log.Fatal("locker file cant be a directory")
	}
	if fileStats.Size() == 0 {
		slog.Info("the file is empty, will init\n")
		return true
	}
	return false
}
func GetLockerFIle() *os.File {
	file, err := os.OpenFile(os.Getenv("LOCKER_PATH")+"/locker.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}
	return file
}
func Config() bool {
	var founded bool
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "LOCKER_PATH") {
			spl := strings.Split(v, "=")
			os.Setenv(spl[0], spl[1])
			founded = true
		}
	}
	return founded
}
