// Package utils provides frequent used functions among project.
package utils

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
)

type config struct {
	Tg tgConfig
}

type tgConfig struct {
	Token        string
	ChatId       int64
	CommanderIds []int64
}

// Reads json file and returns struct of key.
func readConfig(obj interface{}, fieldName string) reflect.Value {
	s := reflect.ValueOf(obj).Elem()
	if s.Kind() != reflect.Struct {
		log.Fatalln("not a struct")
	}
	f := s.FieldByName(fieldName)
	if !f.IsValid() {
		log.Fatalln("not such struct with key")
	}
	return f
}

// Get config json file and returns interpreted interface.
func getConfig(key string) interface{} {
	path, _ := os.Getwd()
	file, _ := os.Open(path + "/config.json")
	defer file.Close()

	c := config{}
	err := json.NewDecoder(file).Decode(&c)
	if err != nil {
		log.Fatalln(err)
	}
	return readConfig(&c, key).Interface()
}

// Returns config of tg as tgConfig struct.
func TgConfig() tgConfig {
	return getConfig("Tg").(tgConfig)
}
