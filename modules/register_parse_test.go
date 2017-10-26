package modules

import (
	//	"fmt"
	tree "go-study/expr2"
	"testing"
)

func Test2(t *testing.T) {
	RegisterKeyword()

	line := `a1 a2||a3 le==1 || lv>=2 || lp>=4`

	expr := Expr(tree.LineExpr(line))
	BeJson(expr)
}
