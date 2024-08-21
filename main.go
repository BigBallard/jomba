package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jessevdk/go-flags"
	"os"
	"path/filepath"
	"strings"
)

type Cli struct {
	File   string `short:"f" long:"file" description:"Path to the JSON file. Can be relative or absolute." required:"true"`
	Output string `short:"o" long:"output" description:"Output file name."`
}

type TokenType string

const (
	TokenTypeField  TokenType = "field"
	TokenTypeObject           = "object"
	TokenTypeArray            = "array"
)

type Token struct {
	Id     string
	Type   TokenType
	Depth  int
	Name   string
	Count  int
	fields []*Token
}

func (c *Token) Find(name string) (*Token, bool) {
	for _, t := range c.fields {
		if t.Name == name {
			return t, true
		}
	}
	return nil, false
}

func (c *Token) AddField(token *Token) {
	c.fields = append(c.fields, token)
}

func (c *Token) Merge(other *Token) {
	combined := make([]*Token, 0)
	for _, ct := range c.fields {
		combined = append(combined, ct)
		if ot, contains := other.Find(ct.Name); contains {
			ct.Count++
			if ct.Type == TokenTypeObject {
				ct.Merge(ot)
			}
		}
	}
	for _, ot := range other.fields {
		if _, contains := c.Find(ot.Name); !contains {
			combined = append(combined, ot)
		}
	}
	c.fields = combined
}

type TokenSet struct {
	root *Token
}

func printTokenSet(t *TokenSet) {
	printObject(t.root, 0)
}

func printObject(object *Token, depth int) {
	padding := strings.Repeat(" ", depth)
	if object.Name == "" {
		fmt.Println(fmt.Sprintf("%s{ %d", padding, object.Count))
	} else {
		fmt.Println(fmt.Sprintf("%s%s: { %d", padding, object.Name, object.Count))
	}
	for _, field := range object.fields {
		if field.Type == TokenTypeField {
			fmt.Println(fmt.Sprintf("%s%s: %d", strings.Repeat(" ", depth+1), field.Name, field.Count))
		} else if field.Type == TokenTypeObject {
			printObject(field, depth+1)
		} else {
			printArray(field, depth+1)
		}
	}
	fmt.Println(fmt.Sprintf("%s}", padding))
}

func printArray(array *Token, depth int) {
	padding := strings.Repeat(" ", depth)
	if array.Name == "" {
		fmt.Println(fmt.Sprintf("%s[ %d", padding, array.Count))
	} else {
		fmt.Println(fmt.Sprintf("%s%s: [ %d", padding, array.Name, array.Count))
	}
	for _, field := range array.fields {
		if field.Type == TokenTypeField {
			fmt.Println(fmt.Sprintf("%s%s: %d", strings.Repeat(" ", depth+1), field.Name, field.Count))
		} else if field.Type == TokenTypeObject {
			printObject(field, depth+1)
		} else {
			printArray(field, depth+1)
		}
	}
	fmt.Println(fmt.Sprintf("%s]", padding))
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

	tokenSet := &TokenSet{
		root: &Token{
			Id:     uuid.New().String(),
			Type:   TokenTypeObject,
			Depth:  0,
			Name:   "",
			Count:  1,
			fields: make([]*Token, 0),
		},
	}

	processObject(jsonObj, tokenSet.root)

	printTokenSet(tokenSet)
}

func processObject(jsonObject map[string]interface{}, container *Token) {
	for key, value := range jsonObject {
		if childObject, ok := value.(map[string]interface{}); ok {
			childContainer, found := container.Find(key)
			if !found {
				childContainer = &Token{
					Id:     uuid.New().String(),
					Name:   key,
					Depth:  container.Depth + 1,
					Count:  0,
					Type:   TokenTypeObject,
					fields: make([]*Token, 0),
				}
			}
			childContainer.Count++
			container.AddField(childContainer)
			processObject(childObject, childContainer)
		} else if childArray, ok := value.([]interface{}); ok {
			arrayContainer, found := container.Find(key)
			if !found {
				arrayContainer = &Token{
					Id:     uuid.New().String(),
					Name:   key,
					Depth:  container.Depth + 1,
					Count:  0,
					Type:   TokenTypeArray,
					fields: make([]*Token, 0),
				}
			}
			arrayContainer.Count++
			container.AddField(arrayContainer)
			processArray(childArray, arrayContainer)
		} else {
			field, fieldFound := container.Find(key)
			if !fieldFound {
				field = &Token{
					Name:  key,
					Type:  TokenTypeField,
					Count: 0,
				}
				container.AddField(field)
			}
			field.Count++
		}
	}
}

func processArray(jsonArray []interface{}, container *Token) {
	for _, item := range jsonArray { // For each item in array
		if object, ok := item.(map[string]interface{}); ok { // and that item is an object
			objectContainer := &Token{
				Id:     uuid.New().String(),
				fields: make([]*Token, 0),
			}
			processObject(object, objectContainer)               // process the object
			for _, objectField := range objectContainer.fields { // iterate over the processed fields
				field, fieldFound := container.Find(objectField.Name) // determine if field exists in current container
				if !fieldFound {                                      // if not
					container.AddField(objectField) // add field
				} else { // otherwise
					field.Count++ // increment
					if field.Type == TokenTypeObject {
						field.Merge(objectField)
					}
				}
			}
		}
	}
}
