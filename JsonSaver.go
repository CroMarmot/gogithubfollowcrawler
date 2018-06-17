package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	jslog = "log/jsonSaver.log"
)

type JsonSaver struct {
	dataInMemory string
	logger       *log.Logger
}

func NewJsonSaver() *JsonSaver {
	ret := new(JsonSaver)
	f, err := os.OpenFile(jslog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	ret.logger = log.New(f, "JsonSaver:\t", log.Lshortfile)
	return ret
}

func (js *JsonSaver) Save(filename string, item interface{}) {
	b, err := json.Marshal(item)
	if err != nil {
		js.logger.Println("error:", err)
		return
	}
	js.SaveFile(filename, string(b))
}

func (js *JsonSaver) Load(filename string, v interface{}) {
	dataFromfile := js.LoadFile(filename)
	err := json.Unmarshal([]byte(dataFromfile), v)
	if err != nil {
		js.logger.Println("[JsonSaver Load ERROR]:", err)
	}
}

func (js *JsonSaver) SaveMem(item interface{}) string {
	b, err := json.Marshal(item)
	if err != nil {
		js.logger.Println("error:", err)
		js.dataInMemory = ""
		return ""
	}
	js.dataInMemory = string(b)
	return js.dataInMemory

}

func (js *JsonSaver) LoadMem(v interface{}) {
	js.logger.Printf(" Load Op on :%v\n", js.dataInMemory)
	err := json.Unmarshal([]byte(js.dataInMemory), v)
	if err != nil {
		js.logger.Println("[JsonSaver Load ERROR]:", err)
	}
}

func (js *JsonSaver) SaveFile(filename string, data string) {
	err := ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (js *JsonSaver) LoadFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
