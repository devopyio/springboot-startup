package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"log"

	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

//The SIGTERM signal is a generic signal used to cause program termination. Unlike SIGKILL, this signal can be blocked, handled, and ignored. It is the normal way to politely ask a program to terminate.
const sigterm = "143"

func main() {
	app := kingpin.New("serc", "A docker images running time.")

	app.HelpFlag.Short('h')

	docker := app.Command("docker", "Run docker image.")
	dockerImageConfig := docker.Arg("config", "Path to your docker image config file.").Required().String()
	dockerServerAddr := docker.Arg("addr", "Server address(domain, ip) on which to run the docker image.").Default("localhost").String()
	dockerPort := docker.Arg("port", "Docker port.").Default("8080").String()
	dockerExposePort := docker.Arg("expose-port", "Port on which the server will be exposed.").Default("8080").String()
	cpusAmm := docker.Arg("cpus", "The ammount of cpu to use.").Default(".5", ".75", "1", "1.5", "2").Strings()

	if kingpin.MustParse(app.Parse(os.Args[1:])) == docker.FullCommand() {
		images, err := loadImagesFromConfig(*dockerImageConfig)
		if err != nil {
			log.Fatal("Couldn't load the images from config file: ", err)
		}
		for _, image := range images {
			fmt.Println(image)
			for _, cpu := range *cpusAmm {
				go checkServer(*dockerServerAddr, *dockerPort, image, cpu)
				if err := exec.Command("sudo", "docker", "run", "-p", fmt.Sprintf("%v:%v", *dockerExposePort, *dockerPort), fmt.Sprintf("--cpus=%v", cpu), image).Run(); err != nil && err.Error() != ("exit status "+sigterm) {
					log.Fatal(err)
				}
			}
		}
	}
}

func checkServer(addr string, port string, dockerImage string, cpu string) {
	begin := time.Now()
	for {
		resp, err := http.Get(fmt.Sprintf("http://%v:%v/", addr, port))
		if err == nil && resp.StatusCode == http.StatusOK {
			end := time.Now()
			elapsed := end.Sub(begin)
			fmt.Printf("(%v cpu): %vs \n", cpu, elapsed.Seconds())

			out, err := exec.Command("sudo", "docker", "ps", "-q", fmt.Sprintf("--filter=ancestor=%v", dockerImage)).Output()
			if err != nil {
				fmt.Printf("Couldn't get the running docker container ID: %v", err)
			}

			if err := exec.Command("docker", "stop", strings.TrimRight(string(out), "\n")).Run(); err != nil {
				fmt.Printf("Couldn't stop the running docker container: %v", err)
			}
			break
		}
	}
}

func loadImagesFromConfig(filename string) (map[string]string, error) {
	imagesFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "can't open the images file- %v", filename)
	}

	var images map[string]string
	err = yaml.Unmarshal(imagesFile, &images)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read the images file- %v", filename)
	}

	return images, nil
}
