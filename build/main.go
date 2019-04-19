package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

var serviceName string

func stopService() {
	err := exec.Command("sc", "stop", serviceName).Run()
	if err != nil {
		fmt.Println(err)
	}
}
func startService() {
	err := exec.Command("sc", "start", serviceName).Run()
	if err != nil {
		fmt.Println(err)
	}
}

func buildServer() {
	cmd := exec.Command("go", "build", "./server")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ":\n" + stderr.String())
		return
	}
	fmt.Println("Build Successful")
}

func main() {
	serviceName = "TeamShovBackend"

	stopService()
	buildServer()
	startService()
}
