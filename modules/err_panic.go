package modules

import (
	"errors"
)

var ERR_DECODE_BASE64 = errors.New("Base64 Decode Err")

var PANIC_UNKNOW_ALERT = "Unknow Alert Type"
var PANIC_UNKNOW_OPERATOR = "Undefine Operator"
var PANIC_SEARCH_SCOPE = "Unknow search scope"
