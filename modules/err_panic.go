package modules

import (
	"errors"
)

var ERR_DECODE_BASE64 = errors.New("Base64 Decode Err")
var ERR_HTTP_REQ = errors.New("Http Request Err")
var ERR_UNMARSHAL_ES_RES = errors.New("Unmarshal Es-Response  Err")

var PANIC_UNKNOW_ALERT = "Unknow Alert Type"
var PANIC_UNKNOW_OPERATOR = "Undefine Operator"
var PANIC_SEARCH_SCOPE = "Unknow search scope"
