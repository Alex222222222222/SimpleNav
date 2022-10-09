package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Render    []string `long:"render" description:"render specific pages, usage: pageName:savePath:hidden"`
	RenderAll []bool   `short:"r" long:"renderAll" description:"render all pages"`
	Server    []bool   `short:"s" long:"server" description:"server the web dir"`
}

func InitFlags() (err error) {

	_, err = flags.ParseArgs(&opts, os.Args)

	return err

}
