package main

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
    "jomba"
    "os"
	"path/filepath"
)

type Cli struct {
	File   string `short:"i" long:"input" description:"Path to the JSON file. Can be relative or absolute." required:"true"`
	Output string `short:"o" long:"output" description:"Output file name."`
	Verbos bool `short:"v" long:"verbose" description:"Prints the traversal to the console."`
}

func main() {
	var cli Cli
	_, err := flags.Parse(&cli)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsing file %s\n", cli.File)
	if filepath.Ext(cli.File) != ".json" {
		fmt.Printf("WARN: file extension of %s is not standard\n", filepath.Ext(cli.File))
	}
	info, statErr := os.Stat(cli.File)
	if statErr != nil {
		panic(statErr)
	}
	if info.Size() == 0 {
		panic(errors.New("file is empty\n"))
	}

	fileBytes, readErr := os.ReadFile(cli.File)
	if readErr != nil {
		panic(fileBytes)
	}

	parser := jomba.NewParser()
	if pErr := parser.ParseBytes(fileBytes); pErr != nil {
		panic(pErr)
	}

	parser.Print()
}
