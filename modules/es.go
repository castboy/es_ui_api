package modules

import (
	"context"
	"fmt"
	tree "go-study/expr2"
	"log"

	"gopkg.in/olivere/elastic.v5"
)

func Cli() *elastic.Client {
	client, err := elastic.NewClient()
	if err != nil {
		log.Fatal("please conf es-cluster-api-host in: GOPATH/src/gopkg.inolivere/elastic.v5/client.go  --line 30")
	}

	return client
}

func Query(p Params, body elastic.Query) *elastic.SearchHits {
	ctx := context.Background()
	client := Cli()
	fetchSrcCtx := FetchSrcCtx(p.T)

	res, err := client.Search().
		Index(ES_INDEX_ALERT).
		Type(EsType[p.T]).
		Query(body).
		FetchSourceContext(fetchSrcCtx).
		From(p.From).Size(p.Size).
		Pretty(true).
		Do(ctx)

	if nil != err {
		fmt.Println("Query Exe Err")
	}

	return res.Hits
}

func IncludesItems(t AlertType) []string {
	switch t {
	case Waf:
		return WafItems
	case Vds:
		return VdsItems
	case Ids:
		return IdsItems
	case Multi:
		return append(append(WafItems, VdsItems...), IdsItems...)
	default:
		panic(PANIC_UNKNOW_ALERT)
	}

	return []string{}
}

func FetchSrcCtx(t AlertType) *elastic.FetchSourceContext {
	include := IncludesItems(t)
	ctx := elastic.NewFetchSourceContext(true).Include(include...)

	return ctx
}

func EsRes(p Params) *elastic.SearchHits {
	body := Expr(tree.LineExpr(p.Query))

	return Query(p, body)
}
