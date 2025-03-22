package jomba

import (
	"errors"
)

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

func (t *Token) Clone() *Token {
	clone := &Token {
		Id: t.Id,
		Name: t.Name,
		Depth: t.Depth,
		Count: t.Count,
		Type: t.Type,
		fields: make([]*Token, 0),
	}
	for _, f := range t.fields {
		clone.fields = append(clone.fields, f.Clone())
	}
	return clone
}

// Find searches for a child Token field with the given name.
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

// Merge takes another token and appends the new fields to the caller or increments the token count if it exists
// in the caller token.
func (c *Token) Merge(other *Token) error {
	if c.Type != other.Type {
		return errors.New("cannot merge differently typed tokens")
	}

	if c.Type == TokenTypeObject {
		return mergeObjects(c, other)
	} else if c.Type == TokenTypeArray {
		return mergeArrays(c, other)
	} else {
		return errors.New("cannot merge literal tokens")
	}
}

func mergeObjects(caller *Token, other *Token) error {
	combined := make([]*Token, 0)
	for _, ct := range caller.fields {
		combined = append(combined, ct)
		if ot, contains := other.Find(ct.Name); contains {
			ct.Count++
			if ct.Type == TokenTypeObject {
				if mErr := ct.Merge(ot); mErr != nil {
					return mErr
				}
			}
		}
	}
	for _, ot := range other.fields {
		if _, contains := caller.Find(ot.Name); !contains {
			combined = append(combined, ot)
		}
	}
	caller.fields = combined
	return nil
}

func mergeArrays(caller *Token, other *Token) error {
	addedTokens := make([]*Token, 0)
	for _, otherToken := range other.fields {
		elem, found := SliceElementWhere[Token](caller.fields, func(t Token) bool {
			return t.Type == otherToken.Type
		})
		if found { // Element exists in the caller array
			typedElem := elem.(*Token)
			if typedElem.Type == TokenTypeField {
				if typedElem.Name == otherToken.Name {
					typedElem.Count++
				} else {
					addedTokens = append(addedTokens, otherToken)
				}
			} else {
				if err := typedElem.Merge(otherToken); err != nil {
					return err
				}
			}
		} else {
			elem, found = SliceElementWhere[Token](addedTokens, func(t Token) bool {
				return t.Name == otherToken.Name
			})
			if found {
				typedElem := elem.(*Token)
				typedElem.Count++
			} else {
				addedTokens = append(addedTokens, otherToken.Clone())
			}
		}
	}
	for _, add := range addedTokens {
		caller.fields = append(caller.fields, add)
	}
	return nil
}

type TokenSet struct {
	root *Token
}

type JsonObject = map[string]interface{}
type JsonArray = []interface{}

// ArrayElement represents an array element of an array. This represents literal values, objects instances (based on structure matching),
// and array elemtent structure matching. The struct holds a hash of the entire element to make matching simpler.
type ArrayElement struct {
	Value interface{}
	Count int
}

// ObjectElement represents an object element of an array. This represents literal values, objects instances (based on structure matching),
// and array elemtent structure matching. The struct holds a hash of the entire element to make matching simpler.
type ObjectElement struct {
	Value interface{}
	Count int
}