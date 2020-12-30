package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/go-vgo/robotgo"
)

type Config struct {
	XmrigLocation    string        `json:"xmrigLocation"`
	Xmfilename       string        `json:"xmFilename"`
	TimeoutStatusOn  time.Duration `json:"timeoutStatusOn"`
	TimeoutStatusOff time.Duration `json:"timeoutStatusOff"`
}

func LoadConfiguration(filename string) Config {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("can't open config file: ", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config

	json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatal("can't decode config JSON: ", err)
	}
	return config
}

func main() {
	previousX := 0
	previousY := 0
	launchStatus := false
	cmd := &exec.Cmd{}
	config := LoadConfiguration("config.json")
	fmt.Println(config)
	for {
		x, y := robotgo.GetMousePos()

		if (x == previousX && y == previousY) && launchStatus == false {
			cmd = exec.Command(config.XmrigLocation + "./" + config.Xmfilename)
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			fmt.Println("Starting Miner...")
			if err := cmd.Start(); err != nil {
				log.Fatal(err)
			}
			launchStatus = true
			fmt.Println("Launch Status: ", launchStatus, "x: ", x, "y: ", y)
		}

		if (x != previousX || y != previousY) && launchStatus {
			fmt.Println("Killing Miner...")

			pgid, err := syscall.Getpgid(cmd.Process.Pid)
			if err == nil {
				syscall.Kill(-pgid, 15)
			}

			cmd.Wait()
			launchStatus = false

			fmt.Println("Launch Status: ", launchStatus, "x: ", x, "y: ", y)

		}

		previousX = x
		previousY = y

		if launchStatus {
			time.Sleep(config.TimeoutStatusOn * time.Second)
		} else {
			time.Sleep(config.TimeoutStatusOff * time.Second)

		}

	}

}
