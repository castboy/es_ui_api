package modules

import (
	"context"

	tree "go-study/expr2"

	"encoding/json"
	"fmt"
	"log"

	"gopkg.in/olivere/elastic.v5"
)

func Query(esIndex string, esType string, body elastic.Query, include interface{}) *[]byte {
	ctx := context.Background()

	client, err := elastic.NewClient()
	if err != nil {
		log.Fatal("please conf es-cluster-api-host in: GOPATH/src/gopkg.inolivere/elastic.v5/client.go  --line 30")
	}

	res, err := client.Search().
		Index(esIndex).
		Type(esType).
		Query(body).
		Source(include).
		From(0).Size(10).
		Pretty(true).
		Do(ctx)
	if err != nil {
		Log("Err", "search response err")
	}

	bytes, err := json.Marshal(res.Hits)
	if nil != err {
		fmt.Println("search response err")
	}

	return &bytes
}

func Includes() interface{} {
	waf := []string{"Client", "Rev", "Msg", "Attack", "Severity", "Maturity", "Accuracy",
		"Hostname", "Uri", "Unique_id", "Ref", "Tags", "Rule", "Version"}

	include, err := elastic.NewFetchSourceContext(true).Include(waf...).Source()
	if nil != err {
		fmt.Println("NewFetchSourceContext.Source err")
	}

	return include
}

func Res(esIndex string, esType string, expr string) *[]byte {
	body := Expr(tree.LineExpr(expr))
	BeJson(body)
	include := Includes()
	return Query(esIndex, esType, body, include)
}
