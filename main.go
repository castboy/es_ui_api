package main

import (
	"encoding/json"
	"fmt"

	"./modules"
	_ "./modules"

	//	"gopkg.in/olivere/elastic.v5"
)

func main() {
	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"], "Ha":{"Add":"ChanPing", "Tel":[123, 456]}}`)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if nil != err {
		fmt.Println("json.Unmarshal Err")
	}
	modules.Parse(f)
}
