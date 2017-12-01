package modules

import (
	//	"context"
	"bytes"
	"encoding/json"
	"fmt"
	tree "go-study/expr2"
	"io/ioutil"
	//	"log"
	"net/http"

	"gopkg.in/olivere/elastic.v5"
)

var EsClient *elastic.Client
var NodesSlice []string

func Nodes(nodes []string) {
	NodesSlice = nodes
}
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

func Query(body string) ([]byte, error) {
	var err error
	var res *http.Response
	var result []byte

	b := bytes.NewBuffer([]byte(body))
	res, err = http.Post("http://"+NodesSlice[0]+":9200/apt/_search", "application/json;charset=utf-8", b)
	if err != nil {
		Log("ERR", "node can not run: %", NodesSlice[0])
		res, err = http.Post("http://"+NodesSlice[1]+":9200/apt/_search", "application/json;charset=utf-8", b)
		if err != nil {
			Log("ERR", "node can not run: %", NodesSlice[1])
			res, err = http.Post("http://"+NodesSlice[2]+":9200/apt/_search", "application/json;charset=utf-8", b)
			if nil != err {
				Log("CRT", "%s", "all es node can not run")
			}
		}
	}

	result, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	return result, err
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

func EsRes(p *Params, e *tree.Expr) ([]byte, error) {

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
