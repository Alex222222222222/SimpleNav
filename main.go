package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

	err := InitFlags()
	if err != nil {
		log.Fatal(err)
	}
}

func InitAll() (err error) {

	err = DBInit()
	if err != nil {
		return err
	}

	err = LoadConfig()
	if err != nil {
		return err
	}

	err = InitTemplate()
	if err != nil {
		return err
	}

	return nil
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

	// TODO database management flags
	// TODO init database here
	if len(opts.AddCategory) > 0 && len(opts.AddCategory[0]) > 0 {

		// init database
		err := DBInit()
		if err != nil {
			log.Fatal(err)
		}

		if len(opts.SetName) == 0 {
			log.Fatal(errors.New("Add new category required a name"))
		}
		name := opts.SetName[0]
		d := ""
		if len(opts.SetDescription) != 0 {
			d = opts.SetDescription[0]
		}
		hidden := false
		if len(opts.SetHidden) != 0 {
			hidden = opts.SetHidden[0]
		}
		p := 0
		if len(opts.SetPriority) != 0 {
			p = opts.SetPriority[0]
		}

		err = AddCategory(name, opts.SetFatherCategory, d, hidden, p)
		if err != nil {
			log.Fatal(err)
		}

	}

	if len(opts.AddLinks) > 0 && len(opts.AddLinks[0]) > 0 {

		// init database
		err := DBInit()
		if err != nil {
			log.Fatal(err)
		}

	}

	// renderAll flag
	if len(opts.RenderAll) > 0 && opts.RenderAll[0] {
		err := InitAll()
		if err != nil {
			log.Fatal(err)
		}
		err = RenderAll()
		if err != nil {
			log.Fatal(err)
		}
	}

	// render flag
	for i := 0; i < len(opts.Render); i += 1 {
		err := InitAll()
		if err != nil {
			log.Fatal(err)
		}

		a := strings.Split(opts.Render[i], ":")
		if len(a) != 3 {
			log.Fatal(
				errors.New("Not enough argument provided in the render flags"),
			)
		}
		hidden, err := strconv.ParseBool(a[2])
		if err != nil {
			log.Fatal(err)
		}

		err = Render(a[0], hidden, a[1])
		if err != nil {
			log.Fatal(err)
		}

	}

	// Server flag
	if len(opts.Server) > 0 && opts.Server[0] {
		err := Server()
		if err != nil {
			log.Fatal(err)
		}
	}

}
