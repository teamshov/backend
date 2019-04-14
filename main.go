package main

import (
	"os"
	"os/exec"
	"syscall"
)

func restartService(name string) {
	binary, lookErr := exec.LookPath("sc")
	if lookErr != nil {
		panic(lookErr)
	}

	stop := []string{"sc", "stop", name}
	start := []string{"sc", "start", name}

	env := os.Environ()

	err := syscall.Exec(binary, stop, env)
	if err != nil {
		panic(err)
	}

	err = syscall.Exec(binary, start, env)
	if err != nil {
		panic(err)
	}
}

func buildServer() {
	binary, lookErr := exec.LookPath("go")
	if lookErr != nil {
		panic(lookErr)
	}

	stop := []string{"sc", "build", "./server"}

	env := os.Environ()

	err := syscall.Exec(binary, stop, env)
	if err != nil {
		panic(err)
	}

}

func main() {
	buildServer()
	restartService("TeamShovBackend")
}
