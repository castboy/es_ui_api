package modules

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Server(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var query, esType string
	index := "apt"

	if val, ok := r.Form["type"]; ok {
		esType = val[0]
	}
	if val, ok := r.Form["query"]; ok {
		query = val[0]
	}

	res := Res(index, esType, query)
	io.WriteString(w, string(*res))
}
