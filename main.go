package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/apex/log"
)

func main() {
	app := kingpin.New("zal", "A docker and prometheus integration.")

	app.HelpFlag.Short('h')

	docker := app.Command("docker", "Run docker image.")
	dockerImage := docker.Arg("image", "Name of your docker image").Required().String()
	cpusAmm := docker.Arg("cpus", "The ammount of cpu to use.").Required().String()
	dockerPort := docker.Arg("port", "Docker port").Default("8000").String()
	dockerExposePort := docker.Arg("expose-port", "Port on which the server will be exposed.").Default("8000").String()

	if kingpin.MustParse(app.Parse(os.Args[1:])) == docker.FullCommand() {
		cmdName := "docker"
		cmdArgs := []string{"run", "-p", fmt.Sprintf("%v:%v", *dockerExposePort, *dockerPort), fmt.Sprintf("--cpus=%v", *cpusAmm), *dockerImage}
		cmd := exec.Command(cmdName, cmdArgs...)

		err := cmd.Start()
		if err != nil {
			log.Errorf("Error starting Cmd: %v", err)
		}

		checkServer(*dockerExposePort)

		err = cmd.Wait()
		if err != nil {
			log.Errorf("Error waiting for Cmd: %v", err)
		}

	}
}

func checkServer(port string) {
	begin := time.Now()
	for {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%v/", port))
		if err == nil && resp.StatusCode == http.StatusOK {
			end := time.Now()
			elapsed := end.Sub(begin)
			fmt.Printf("Server started: %vs", elapsed.Seconds())
			break
		}
	}
}
