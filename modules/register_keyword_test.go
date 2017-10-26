package modules

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	abs := Read()
	Load(abs)
	for _, v := range exprBuildinKeywords {
		fmt.Println(*v)
	}
}
