package es_ui_api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type OtherConf struct {
	ClientID string
}

type HttpConf struct {
	Host string
	Port int
	Url  string
}

type Broker struct {
	Basic HttpConf
	Other OtherConf
}

type Client interface {
	Listen()
}

func NewOtherConf(conf string) OtherConf {
	return OtherConf{
		ClientID: conf,
	}
}

func Dial(http HttpConf, other OtherConf) (Broker, error) {
	broker := Broker{
		Basic: http,
		Other: other,
	}
	return broker, nil
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func (broker Broker) Listen() {
	fmt.Println(broker.Other.ClientID)

	http.HandleFunc(broker.Basic.Url, HelloServer)
	err := http.ListenAndServe(":"+strconv.Itoa(broker.Basic.Port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
