package main

import (
	"encoding/json"
	"log"
)

var difficultData = `{
  "items": {
    "DuW4cNvImHwF3cVx": {
      "key1": "val1",
      "key2": "val2",
      "key3": "val3"
    },
    "SkygbBwMUQp7zayq": {
      "key1": "val1",
      "key2": "val2",
      "key3": "val3"
    }
  }
}`

type difficultApi struct {
	Items map[string]record `json:"items"`
}

type record struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
	Key3 string `json:"key3"`
}

func main() {
	var difficult difficultApi
	err := json.Unmarshal([]byte(difficultData), &difficult)
	if err != nil {
		log.Fatal(err)
	}
	for i := range difficult.Items {
		log.Println("id", i)
		log.Println("key1", difficult.Items[i].Key1)
		log.Println("key2", difficult.Items[i].Key2)
		log.Println("key3", difficult.Items[i].Key3)
	}
}
