package modules

import (
	//	"fmt"
	tree "go-study/expr2"
	"testing"
)

func Tree2() *tree.Expr {
	for _, v := range exprBuildinKeywords {
		tree.RegisterKeyword(*v)
	}

	line := `a1 a2||a3 le>1||lv<1`

	return tree.LineExpr(line)
}
func Test2(t *testing.T) {
	Read()
	tree := Tree2()
	expr := Expr(tree)
	BeJson(expr)
}
