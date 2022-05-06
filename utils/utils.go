// Package utils provides frequent used functions among project.
package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
)

// Returns if a is contained in s.
func Contains[T comparable](s []T, a T) bool {
	for _, b := range s {
		if b == a {
			return true
		}
	}
	return false
}

// Returns smaller value.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Converts any type as bytes.
func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	err := encoder.Encode(i)
	if err != nil {
		log.Fatalln(err)
	}
	return aBuffer.Bytes()
}

// Converts bytes to json.
func Bytes2Json(data []byte, i interface{}) {
	r := bytes.NewReader(data)
	err := json.NewDecoder(r).Decode(i)
	if err != nil {
		log.Fatalln(err)
	}
}
