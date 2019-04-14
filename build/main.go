package main

import (
	"os/exec"
)

func restartService(name string) {
	err := exec.Command("sc", "stop", name).Run()
	if err != nil {
		panic(err)
	}

	err = exec.Command("sc", "start", name).Run()
	if err != nil {
		panic(err)
	}
}

func buildServer() {
	err := exec.Command("go", "build", "./server").Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	buildServer()
	restartService("TeamShovBackend")
}
