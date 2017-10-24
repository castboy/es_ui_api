package modules

import (
	"encoding/json"
	"fmt"
	tree "go-study/expr2"

	"testing"

	"gopkg.in/olivere/elastic.v5"
)

func IsNested(me *tree.Atomic) bool {
	return false
}

func IsTerm(me *tree.Atomic) bool {
	return false
}

func IsMatch(me *tree.Atomic) bool {
	return true
}

var boolQuery *elastic.BoolQuery

func Expr(expr *tree.Expr) elastic.Query {
	query := make([]elastic.Query, 0)

	if expr.IsAtomic() {
		exp := AtomicExpr(expr.Atomic)
		switch e := exp.(type) {
		case elastic.Query:
			return e
		}
	} else {
		for _, v := range expr.Children {
			query = append(query, Expr(v))
		}
	}

	return BoolExpr(expr, query)
}

func BeJson(query elastic.Query) {
	src, err := query.Source()
	if err != nil {
	}

	data, err := json.Marshal(src)
	if err != nil {

	}
	s := string(data)
	fmt.Println(s)
}

func BoolExpr(me *tree.Expr, query []elastic.Query) *elastic.BoolQuery {
	switch me.Logic {
	case tree.LogicAnd:
		return elastic.NewBoolQuery().Must(query...).MinimumShouldMatch("1")
	case tree.LogicOr:
		return elastic.NewBoolQuery().Should(query...).MinimumShouldMatch("1")
	case tree.LogicNot:
		return elastic.NewBoolQuery().MustNot(query...).MinimumShouldMatch("1")
	default:
		return elastic.NewBoolQuery().Must(query...).MinimumShouldMatch("1")
	}

	return nil
}

func AtomicExpr(me *tree.Atomic) interface{} {
	if IsNested(me) {
		if IsTerm(me) {
			return elastic.NewNestedQuery("xdr", elastic.NewTermQuery(me.K.Key, me.V))
		} else {
			return elastic.NewNestedQuery("xdr", elastic.NewMatchQuery(me.K.Key, me.V))
		}
	} else {
		if IsTerm(me) {
			return elastic.NewTermQuery(me.K.Key, me.V)
		} else {
			return elastic.NewMatchQuery(me.K.Key, me.V)
		}
	}

	return 0
}

func Tree() *tree.Expr {
	for i := 1; i < 10; i++ {
		tree.RegisterKeyword(tree.Keyword{
			Key:   fmt.Sprintf("zone%d", i),
			Scope: tree.ScopeZone,
			Type:  0,
		})
	}

	for i := 1; i < 10; i++ {
		tree.RegisterKeyword(tree.Keyword{
			Key:   fmt.Sprintf("obj%d", i),
			Scope: tree.ScopeObject,
			Type:  0,
		})
	}

	line := `a1 && a2 && a3`

	return tree.LineExpr(line)
}

func Test1(t *testing.T) {
	tree := Tree()
	expr := Expr(tree)
	BeJson(expr)
}
