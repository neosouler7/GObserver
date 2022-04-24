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

func readConfig(obj interface{}, fieldName string) reflect.Value {
	s := reflect.ValueOf(obj).Elem()
	if s.Kind() != reflect.Struct {
		log.Fatalln("error not struct")
	}
	f := s.FieldByName(fieldName)
	if !f.IsValid() {
		log.Fatalln("error not such field")
	}
	return f
}

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

func TgConfig() tgConfig {
	return getConfig("Tg").(tgConfig)
}
