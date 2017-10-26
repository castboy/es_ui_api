package modules

import (
	"context"

	tree "go-study/expr2"

	"encoding/json"

	"gopkg.in/olivere/elastic.v5"
)

func Query(esIndex string, esType string, body elastic.Query) *[]byte {
	ctx := context.Background()

	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}

	res, err := client.Search().
		Index(esIndex).
		Type(esType).
		Query(body).
		From(0).Size(10).
		Pretty(true).
		Do(ctx)
	if err != nil {
		panic(err)
	}

	bytes, err := json.Marshal(res.Hits)
	if nil != err {
		panic(err)
	}

	return &bytes
}

func Res(esIndex string, esType string, expr string) *[]byte {
	body := Expr(tree.LineExpr(expr))
	return Query(esIndex, esType, body)
}
