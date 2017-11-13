package modules

import (
	"encoding/json"
	"fmt"
	tree "go-study/expr2"

	"gopkg.in/olivere/elastic.v5"
)

func isInt(me *tree.Atomic) bool {
	fmt.Println("me.K.Type:", me.K.Type)
	return me.K.Type&1 != 0 //0非整数， 1整数
}

func isNested(me *tree.Atomic) bool {
	me.K.Type >>= 1
	return me.K.Type&1 != 0 //0非嵌套， 1嵌套
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
	fmt.Println("me.K", *me.K)
	var v elastic.Query

	switch me.Op {
	case tree.OpGeEq:
		v = elastic.NewRangeQuery(me.K.Name).Gte(me.V)
	case tree.OpGe:
		v = elastic.NewRangeQuery(me.K.Name).Gt(me.V)
	case tree.OpLeEq:
		v = elastic.NewRangeQuery(me.K.Name).Lte(me.V)
	case tree.OpLe:
		v = elastic.NewRangeQuery(me.K.Name).Lt(me.V)
	case tree.OpEq:
		if isInt(me) {
			v = elastic.NewRangeQuery(me.K.Name).Gte(me.V).Lte(me.V)
		} else {
			v = elastic.NewTermQuery(me.K.Name, me.V)
		}
	case tree.OpInclude:
		if me.K.Scope == tree.ScopeAll {
			v = elastic.NewFuzzyQuery("_all", me.V)
		} else {
			v = elastic.NewFuzzyQuery(me.K.Name, me.V)
		}
	case tree.OpNeq:
		if isInt(me) {
			fmt.Println("isInt")
			v = elastic.NewRangeQuery(me.K.Name).Gt(me.V).Lt(me.V)
		} else {
			fmt.Println("isNotInt")
			v = elastic.NewBoolQuery().MustNot(elastic.NewMatchQuery(me.K.Name, me.V))
		}
	default:
		panic(PANIC_UNKNOW_OPERATOR)
	}
	fmt.Println("me.K", *me.K)
	if isNested(me) {
		fmt.Println("isNested")
		v = elastic.NewNestedQuery("xdr", v)
	}

	return v
}
