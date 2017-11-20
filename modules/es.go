package modules

import (
	"context"
	"encoding/json"
	"fmt"
	tree "go-study/expr2"
	"log"

	"gopkg.in/olivere/elastic.v5"
)

var EsClient *elastic.Client

func Cli(nodes []string, port string) {
	var err error
	var nodePort []string

	for _, v := range nodes {
		nodePort = append(nodePort, v+port)
	}

	EsClient, err = elastic.NewClient(elastic.SetURL(nodePort...))
	if err != nil {
		log.Fatal("init es client err")
	}
}

func Query(p *Params, body elastic.Query) *elastic.SearchHits {
	ctx := context.Background()
	fetchSrcCtx := FetchSrcCtx(p)

	res, err := EsClient.Search().
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

func IncludesItems(p *Params) []string {
	switch p.T {
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

func FetchSrcCtx(p *Params) *elastic.FetchSourceContext {
	include := IncludesItems(p)
	ctx := elastic.NewFetchSourceContext(true).Include(include...)

	return ctx
}

func RecoverLineExpr(p *Params) (expr *tree.Expr, err ExprErr) {
	defer func() {
		err = ExprErr(fmt.Sprint(recover()))
	}()

	expr = tree.LineExpr(p.Query)

	return expr, ""
}

func EsRes(p *Params, e *tree.Expr) *elastic.SearchHits {

	body := Expr(e)
	i, _ := body.Source()

	bytes, _ := json.Marshal(i)
	fmt.Println(string(bytes))

	return Query(p, body)
}
