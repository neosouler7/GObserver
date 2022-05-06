// Package utils provides frequent used functions among project.
package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

type config struct {
	Tg    tgConfig
	Pairs map[string]interface{}
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

// Returns pairs of certain exchange.
func GetPairs(exchange string) []string {
	var pairs []string
	for market, symbols := range getConfig("Pairs").(map[string]interface{})[exchange].(map[string]interface{}) {
		for _, symbolInfo := range symbols.([]interface{}) {
			pairs = append(pairs, fmt.Sprintf("%s:%s", market, symbolInfo))
		}
	}
	return pairs
}

// Returns pairs of certain exchange.
func GetVolumeMap(exchange string) map[string]string {
	m := make(map[string]string)
	for _, p := range getConfig("Pairs").(map[string]interface{})[exchange].([]interface{}) {
		s := strings.Split(p.(string), ":")
		m[fmt.Sprintf("%s:%s", s[0], s[1])] = s[2]
	}
	return m
}
