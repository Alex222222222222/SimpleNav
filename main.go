package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

var Config map[string]string

func InitWebDir() (err error) {

	err = os.RemoveAll("./web")

	cmd := exec.Command("mkdir", "./web")
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("cp", "-r", "./static", "./web")
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil

}

func init() {

	err := InitWebDir()
	if err != nil {
		panic(err)
	}

	err = LoadConfig()
	if err != nil {
		panic(err)
	}

	err = DBInit()
	if err != nil {
		log.Panic(err)
	}

	err = InitTemplate()
	if err != nil {
		log.Panic(err)
	}

	err = Render("Index", true, "./web/index.html")
	if err != nil {
		panic(err)
	}

}

func LoadConfig() (err error) {
	config, err := os.ReadFile("./config.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(config, &Config)
	if err != nil {
		return err
	}

	return nil

}

func main() {

	Server()

}
