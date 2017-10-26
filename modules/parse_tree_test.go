package modules

import (
	"fmt"
	tree "go-study/expr2"
	"testing"
)

func Test1(t *testing.T) {
	for i := 1; i < 10; i++ {
		tree.RegisterKeyword(tree.Keyword{
			Key:   fmt.Sprintf("zone%d", i),
			Name:  fmt.Sprintf("zone%d", i),
			Scope: tree.ScopeZone,
			Type:  0,
		})
	}

	for i := 1; i < 10; i++ {
		tree.RegisterKeyword(tree.Keyword{
			Key:   fmt.Sprintf("obj%d", i),
			Name:  fmt.Sprintf("obj%d", i),
			Scope: tree.ScopeObject,
			Type:  0,
		})
	}

	line := `a1 a2||a3 zone1="v1 \'v2\'"&&(obj1>=0&&(obj2>1||obj3<1)||!(obj4==0&&obj5==0))||(obj6<0&&obj7>0)`

	expr := Expr(tree.LineExpr(line))
	BeJson(expr)
}
