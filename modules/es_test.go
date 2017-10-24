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

var query = make([]elastic.Query, 0)
var boolQuery *elastic.BoolQuery

func Expr(expr *tree.Expr) *elastic.BoolQuery {
	if expr.IsAtomic() {
	} else {
		for _, v := range expr.Children {
			expr := AtomicExpr(v.Atomic)
			switch e := expr.(type) {
			case elastic.Query:
				query = append(query, e)
			}
		}
		fmt.Println(expr.Logic)
		boolQuery = boolExpr(expr, query)
		src, err := boolQuery.Source()
		if err != nil {
		}

		data, err := json.Marshal(src)
		if err != nil {

		}
		s := string(data)
		fmt.Println(s)
	}

	return nil
}

func boolExpr(me *tree.Expr, query []elastic.Query) *elastic.BoolQuery {
	switch me.Logic {
	case tree.LogicAnd:
		return elastic.NewBoolQuery().Must(query...).MinimumShouldMatch("1")
	case tree.LogicOr:
		return elastic.NewBoolQuery().Should(query...).MinimumShouldMatch("1")
	case tree.LogicNot:
		return elastic.NewBoolQuery().MustNot(query...).MinimumShouldMatch("1")
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
	Expr(Tree())
}
