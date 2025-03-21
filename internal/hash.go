package jomba

import (
	"crypto/sha256"
    "encoding/hex"
    "fmt"
    "reflect"
    "strconv"
)


var (
	trueHash string
	falseHash string
)

func init() {
	trueHash = string(sha256.New().Sum([]byte{'1'}))
	falseHash = string(sha256.New().Sum([]byte{'0'}))
}

func HashLiteral(value interface{}) (string, string, error) {
	var hash string
	var literal string
	if strValue, isStr := value.(string); isStr {
		hash, literal = stringHash(strValue)
	} else if floatValue, isFloat := value.(float64); isFloat {
		hash, literal = floatHash(floatValue)
	} else if boolValue, isBool := value.(bool); isBool {
		hash, literal = boolHash(boolValue)
	} else {
		realType := reflect.TypeOf(value)
		return "","", fmt.Errorf("invalid literal type %s", realType.Name())
	}
	return hash, literal, nil
}

func stringHash(value string) (string, string) {
	return hash([]byte(value)), value
}

func intHash(value int) (string, string) {
	strValue := strconv.Itoa(value)
	return hash([]byte(strValue)), strValue
}

func floatHash(value float64) (string, string) {
	strValue := strconv.FormatFloat(value, 'g', -1,  64)
	return hash([]byte(strValue)), strValue
}

func boolHash(value bool) (string, string) {
	if value {
		return trueHash, "true"
	} else {
		return falseHash, "false"
	}
}

func hash(value []byte) string {
	return hex.EncodeToString(sha256.New().Sum(value))
}