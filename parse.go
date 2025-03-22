package jomba

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	jomba "jomba/internal"
)

var ParseErr = errors.New("failed to parse bytes")
type Parser struct {
	tokenSet *TokenSet
}

func NewParser() *Parser {
	return &Parser{
		tokenSet: &TokenSet{
			root: &Token{
				Id:     uuid.New().String(),
				Type:   TokenTypeObject,
				Depth:  0,
				Name:   "",
				Count:  1,
				fields: make([]*Token, 0),
			},
		},
	}
}

func (p *Parser) ParseBytes(b []byte) error {
	var jsonObj map[string]interface{}
	if uErr := json.Unmarshal(b, &jsonObj); uErr != nil {
		return errors.Join(ParseErr, uErr)
	}
	return p.parseObject(jsonObj, p.tokenSet.root)
}

// parseObject takes a JsonObject and iterates over the objet keys.
func (p *Parser) parseObject(jsonObject JsonObject, container *Token) error {
	for key, value := range jsonObject {
		valueToken, found := container.Find(key)
		if obj, isObj := value.(JsonObject); isObj { // Key is JsonObject
			if !found { // Not found, create one and add to container
				valueToken = &Token{
					Id: uuid.NewString(),
					Name: key,
					Count: 0,
					Depth: container.Depth + 1,
					Type: TokenTypeObject,
				}
				container.AddField(valueToken)
			}
			valueToken.Count++ // Increment instance counter
			if err := p.parseObject(obj, valueToken); err != nil { // Parse the next object in the in its own context
				return err
			}
		} else if arr, isArr :=value.(JsonArray); isArr {
			if !found {
				valueToken = &Token{
					Id: uuid.NewString(),
					Name: key,
					Count: 0,
					Depth: container.Depth + 1,
					Type: TokenTypeArray,
				}
				container.AddField(valueToken)
			}
			valueToken.Count++
			if err := p.parseArray(arr, valueToken); err != nil {
				return err
			}
		} else {
			if !found {
				valueToken = &Token{
					Id: uuid.NewString(),
					Name: key,
					Count: 0,
					Depth: container.Depth + 1,
					Type: TokenTypeField,
				}
				container.AddField(valueToken)
			}
			valueToken.Count++
		}
	}
	return nil
}

func (p *Parser) parseArray(jsonArray JsonArray, container *Token) error {
	// Iterate through each element in the array and determine its type. Objects will be an aggregate of the first level of keys.
	// This will make processing manageable and user presentation sensible since array values are generally either
	// arbitrary or context specific. Literals will be counted on a per value basis, so each unique literal value will be
	// an element entry.
	objToken := &Token{Type: TokenTypeObject, Depth: container.Depth+1}
	arrToken := &Token{Type: TokenTypeArray, Depth: container.Depth+1}
	fieldTokens := make(map[string]*Token)
	for _, value := range jsonArray {
		if obj, isObj := value.(JsonObject); isObj {
			objContainer := &Token{Type: TokenTypeObject}
			if err := p.parseObject(obj, objContainer); err != nil {
				return err
			}
			objToken.Count++
			if mErr := objToken.Merge(objContainer); mErr != nil {
				return mErr
			}
		} else if arr, isArr := value.(JsonArray); isArr {
			arrContainer := &Token{Type: TokenTypeArray}
			if err := p.parseArray(arr, arrContainer); err != nil {
				return err
			}
			arrToken.Count++
			if mErr := arrToken.Merge(arrContainer); mErr != nil {
				return mErr
			}
		} else {
			hash, literal, hErr := jomba.HashLiteral(value)
			if hErr != nil {
				return hErr
			}
			if fieldToken, found := fieldTokens[hash]; found {
				fieldToken.Count++
			} else {
				fieldTokens[hash] = &Token {
					Name: literal,
					Count: 1,
					Type: TokenTypeField,
					Depth: container.Depth + 1,
				}
			}
		}
	}
	for _, field := range fieldTokens {
		container.AddField(field)
	}
	if objToken.Count > 0 {
		container.AddField(objToken)
	}
	if arrToken.Count > 0 {
		container.AddField(arrToken)
	}
	return nil
}

func (p *Parser) Print() {
	printTokenSet(p.tokenSet)
}
