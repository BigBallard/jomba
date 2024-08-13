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

type TokenType string

const (
	TokenTypeField       TokenType = "FIELD"
	TokenTypeObjectStart TokenType = "OBJECT_START"
	TokenTypeArrayStart  TokenType = "ARRAY_START"
	TokenTypeObjectEnd   TokenType = "OBJECT_END"
	TokenTypeArrayEnd    TokenType = "ARRAY_END"
)

type Token struct {
	Type  TokenType `json:"token_type"`
	Depth int
	Name  string
	Count int
}

type TokenSet struct {
	tokens []*Token
}

func (t *TokenSet) AddToken(token *Token) {
	t.tokens = append(t.tokens, token)
}

func (t *TokenSet) FindToken(name string) (*Token, error) {
	i := len(t.tokens) - 1
	for i >= 0 {
		if t.tokens[i].Type == TokenTypeObjectStart {
			return nil, errors.New("field not in object scope")
		}
		if t.tokens[i].Name == name {
			return t.tokens[i], nil
		}
		i--
	}
	return nil, errors.New("field not found")
}

func printTokenSet(tokenSet *TokenSet) {
	for i, token := range tokenSet.tokens {
		if token.Type == TokenTypeObjectStart || token.Type == TokenTypeObjectEnd {
			printObject(token, i == len(tokenSet.tokens)-1)
		} else if token.Type == TokenTypeArrayStart || token.Type == TokenTypeArrayEnd {
			printArray(token)
		} else {
			printField(token)
		}
	}
}

func printField(token *Token) {
	pad := strings.Repeat(" ", token.Depth*2)
	fmt.Printf("%s%s: %d\n", pad, token.Name, token.Count)
}

func printObject(token *Token, end bool) {
	pad := strings.Repeat(" ", token.Depth*2)
	if end {
		fmt.Println("}")
	} else if token.Name == "" {
		if token.Depth == 0 || token.Name == "" {
			fmt.Printf("%s{\n", pad)
		} else {
			fmt.Printf("%s}\n", pad)
		}
	} else {
		fmt.Printf("%s%s: {\n", pad, token.Name)
	}
}

func printArray(token *Token) {
	pad := strings.Repeat(" ", token.Depth*2)
	if token.Name == "" {
		fmt.Printf("%s]\n", pad)
	} else {
		fmt.Printf("%s%s: [\n", pad, token.Name)
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

	tokenSet := &TokenSet{
		tokens: []*Token{},
	}

	processObject(0, "", jsonObj, tokenSet)

	printTokenSet(tokenSet)
}

func processObject(depth int, name string, object map[string]interface{}, tokenSet *TokenSet) {
	tokenSet.AddToken(&Token{
		Type:  TokenTypeObjectStart,
		Name:  name,
		Depth: depth,
	})
	for k, v := range object {
		if objValue, ok := v.(map[string]interface{}); ok {
			processObject(depth+1, k, objValue, tokenSet)
		} else if arrValue, ok := v.([]interface{}); ok {
			processArray(depth+1, k, arrValue, tokenSet)
		} else {
			token, findErr := tokenSet.FindToken(k)
			if findErr != nil {
				token = &Token{
					Type:  TokenTypeField,
					Name:  k,
					Count: 0,
					Depth: depth + 1,
				}
				tokenSet.AddToken(token)
			}
			token.Count++
		}
	}
	tokenSet.AddToken(&Token{
		Type:  TokenTypeObjectEnd,
		Depth: depth,
	})
}

func processArray(depth int, name string, array []interface{}, tokenSet *TokenSet) {
	tokenSet.AddToken(&Token{
		Type:  TokenTypeArrayStart,
		Depth: depth,
		Name:  name,
	})
	fields := make(map[string]int)
	for _, v := range array {
		if objValue, ok := v.(map[string]interface{}); ok {
			objSet := &TokenSet{tokens: make([]*Token, 0)}
			processObject(depth+1, "", objValue, objSet)

			for _, f := range objSet.tokens {
				if f.Type == TokenTypeField {
					if _, ok := fields[f.Name]; ok {
						fields[f.Name] += 1
					} else {
						fields[f.Name] = 1
					}
				}
			}
		} else if arrValue, ok := v.([]interface{}); ok {
			processArray(depth+1, "", arrValue, tokenSet)
			// TODO Fix this and correctly process array within array.
		}
	}
	tokenSet.AddToken(&Token{Type: TokenTypeObjectStart, Depth: depth + 1, Name: ""})
	for fk, fv := range fields {
		tokenSet.AddToken(&Token{
			Type:  TokenTypeField,
			Name:  fk,
			Count: fv,
			Depth: depth + 2,
		})
	}
	tokenSet.AddToken(&Token{Type: TokenTypeObjectEnd, Depth: depth + 1})
	tokenSet.AddToken(&Token{
		Type:  TokenTypeArrayEnd,
		Depth: depth,
	})
}
