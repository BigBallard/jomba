package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"path/filepath"
	"strings"
)

type Cli struct {
	File   string `short:"f" long:"file" description:"Path to the JSON file. Can be relative or absolute." required:"true"`
	Output string `short:"o" long:"output" description:"Output file name."`
}

type PropMap struct {
	Level    int
	Entries  map[string]int
	Children map[string]PropMap
}

func (p *PropMap) Print() {
	padding := strings.Repeat(" ", p.Level)
	for key, value := range p.Entries {
		fmt.Printf("%s%s: %d\n", padding, key, value)
		if child, ok := p.Children[key]; ok {
			child.Print()
		}
	}
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

	var jsonObj map[string]interface{}
	if err = json.Unmarshal(fileBytes, &jsonObj); err != nil {
		panic(err)
	}
	propMap := PropMap{
		Level:    0,
		Entries:  map[string]int{},
		Children: map[string]PropMap{},
	}
	processObject(jsonObj, &propMap)

	propMap.Print()
}

func processObject(object map[string]interface{}, propMap *PropMap) {
	for key, value := range object {
		if _, ok := propMap.Entries[key]; !ok {
			propMap.Entries[key] = 0
		}
		propMap.Entries[key] = propMap.Entries[key] + 1
		if _, ok := value.(map[string]interface{}); ok {
			nested := PropMap{
				Level:    propMap.Level + 1,
				Entries:  map[string]int{},
				Children: map[string]PropMap{},
			}
			propMap.Children[key] = nested
			processObject(value.(map[string]interface{}), &nested)
		} else if array, ok := value.([]interface{}); ok {
			for _, entry := range array {
				if _, ok := entry.(map[string]interface{}); ok {
					nested := PropMap{
						Level:    propMap.Level + 1,
						Entries:  map[string]int{},
						Children: map[string]PropMap{},
					}
					propMap.Children[key] = nested
					processObject(entry.(map[string]interface{}), &nested)
				}
			}
		}
	}
}
