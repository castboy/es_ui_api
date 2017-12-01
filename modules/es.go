package modules

import (
	//	"context"
	"bytes"
	"encoding/json"
	"fmt"
	tree "go-study/expr2"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/olivere/elastic.v5"
)

var EsClient *elastic.Client

type CurlBody struct {
	From  int         `json:"from"`
	Size  int         `json:"size"`
	Query interface{} `json:"query"`
}

type CurlRes struct {
	Took     int  `json:"took"`
	Time_out bool `time_out`
	_shards  interface{}
	Hits     HitsOuter `json:"hits"`
}

type HitsOuter struct {
	Total int       `json:"total"`
	Hits  HitsInner `json:"hits"`
}

type HitsInner []OneResComplete

type OneResComplete struct {
	Source OneResSource `json:"_source"`
}

type OneResSource interface{}

func Cli(nodes []string, port string) {
	var err error
	var nodePort []string

	for _, v := range nodes {
		nodePort = append(nodePort, "http://"+v+":"+port)
	}

	EsClient, err = elastic.NewClient(elastic.SetURL(nodePort...))
	if err != nil {
		Log("CRT", "new es client failed, nodes: %s", nodes)
	}
}

func Query(body string) string {
	b := bytes.NewBuffer([]byte(body))
	res, err := http.Post("http://10.88.1.102:9200/apt/_search", "application/json;charset=utf-8", b)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return ""
	}

	var curlRes CurlRes

	json.Unmarshal(result, &curlRes)

	fmt.Println(curlRes.Hits.Total)

	//	var one OneResComplete
	for _, v := range curlRes.Hits.Hits {
		fmt.Println(v)
		//		json.Unmarshal([]byte(v), &one)
		//		fmt.Println(one.Source)
	}

	return string(result)
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

func EsRes(p *Params, e *tree.Expr) string {

	body := Expr(e)
	i, _ := body.Source()

	curlBody := CurlBody{
		From:  p.From,
		Size:  p.Size,
		Query: i,
	}

	bytes, _ := json.Marshal(curlBody)
	Log("INF", "es query exe: %s", string(bytes))

	return Query(string(bytes))
}
