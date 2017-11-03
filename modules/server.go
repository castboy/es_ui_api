package modules

import (
	"encoding/base64"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const DefaultURL = "http://10.88.1.102:9200"

func Server(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	Includes()
	r.ParseForm()

	var query, esType string
	index := "apt"

	if val, ok := r.Form["type"]; ok {
		esType = val[0]
	}
	if val, ok := r.Form["query"]; ok {
		sDec, _ := base64.StdEncoding.DecodeString(val[0])
		query = string(sDec)
	}

	res := Res(index, esType, query)
	io.WriteString(w, string(*res))
}
