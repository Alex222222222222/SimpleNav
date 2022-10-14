package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Render      []string `long:"render" description:"render specific pages, usage: pageName:savePath:hidden"`
	RenderAll   []bool   `short:"r" long:"renderAll" description:"render all pages"`
	Server      []bool   `short:"s" long:"server" description:"server the web dir"`
	AddLinks    []string `short:"l" long:"addLinks" description:"add new links, use single url or specify other things using --name --description --img --fatherCategory --priority --tags flags"`
	AddCategory []string `short:"c" long:"addCategory" description:"add new categories, use single url or specify other things using --name --description --fatherCategory --priority --hidden flags"`
	// TODO "add tags command"

	SetName           []string `long:"name" description:"set name for category or links"`
	SetDescription    []string `long:"des" description:"set description for category or links"`
	SetImg            []string `long:"img" description:"set image for links"`
	SetFatherCategory []string `long:"fc" description:"set father category for category or links"`
	SetPriority       []int    `long:"prio" description:"set list priority for category or links"`
	SetTags           []string `long:"tags" description:"set tags for links, use tag1:tag2... or use this tag multiple times"`
	SetHidden         []bool   `long:"hidden" description:"set hidden for category"`
}

func InitFlags() (err error) {

	_, err = flags.ParseArgs(&opts, os.Args)

	return err

}
