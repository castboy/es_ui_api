package modules

import (
	"io"
	"net/http"
)

func Server(w http.ResponseWriter, req *http.Request) {
	res := Demo()
	io.WriteString(w, res)
}
