package es_ui_api

import (
	"io"
	"net/http"

	"github.com/castboy/es_ui_api/modules"
	"github.com/julienschmidt/httprouter"
)

func Server(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	res := modules.Demo()
	io.WriteString(w, res)
}
