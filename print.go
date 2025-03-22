package jomba

import (
	"fmt"
	"strings"
)

func printTokenSet(t *TokenSet) {
	printObject(t.root, 0)
}

func printObject(object *Token, depth int) {
	padding := strings.Repeat(" ", depth)
	if object.Name == "" {
		if depth > 0 {
			fmt.Println(fmt.Sprintf("%s{ %d", padding, object.Count))
		} else {
			fmt.Println(fmt.Sprintf("%s{", padding))
		}
	} else {
		fmt.Println(fmt.Sprintf("%s%s : { %d", padding, object.Name, object.Count))
	}
	for _, field := range object.fields {
		if field.Type == TokenTypeField {
			fmt.Println(fmt.Sprintf("%s%s : %d", strings.Repeat(" ", depth+1), field.Name, field.Count))
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
	if array.Name == "" { // Not an object field
		fmt.Println(fmt.Sprintf("%s[ %d", padding, array.Count))
	} else {
		fmt.Println(fmt.Sprintf("%s%s : [ %d", padding, array.Name, array.Count))
	}
	for _, field := range array.fields {
		if field.Type == TokenTypeField {
			fmt.Println(fmt.Sprintf("%s%s : %d", strings.Repeat(" ", depth+1), field.Name, field.Count))
		} else if field.Type == TokenTypeObject {
			printObject(field, depth+1)
		} else {
			printArray(field, depth+1)
		}
	}
	fmt.Println(fmt.Sprintf("%s]", padding))
}
