package es_ui_api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type OtherConf struct {
	Other string
}

type Broker struct {
	Host  string
	Port  int
	Url   string
	Other OtherConf
}

type Client interface {
	Listen()
}

func NewOtherConf(conf string) OtherConf {
	return OtherConf{
		Other: conf,
	}
}

func Dial(host string, port int, url string, conf OtherConf) (Broker, error) {
	broker := Broker{
		Host:  host,
		Port:  port,
		Url:   url,
		Other: conf,
	}
	return broker, nil
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func (broker Broker) Listen() {
	fmt.Println(broker.Other.Other)

	http.HandleFunc(broker.Url, HelloServer)
	err := http.ListenAndServe(":"+strconv.Itoa(broker.Port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
