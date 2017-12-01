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
	Total int64     `json:"total"`
	Hits  HitsInner `json:"hits"`
}

type HitsInner []OneResComplete

type OneResComplete struct {
	Source OneResSource `json:"_source"`
}

type OneResSource interface{}

type ResHits []OneResSource

type ResHitsType struct {
	Type string
}

var NodesSlice []string

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

func Nodes(nodes []string) {
	NodesSlice = nodes
}

func CurlEs(body string) ([]byte, error) {
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

func Query(body string) Res {
	var resHitsType ResHitsType
	var res []interface{}
	var curlRes CurlRes

	result, err := CurlEs(body)
	if nil != err {
		Log("ERR", "%s", "CurlEs res.Body")
	} else {
		json.Unmarshal(result, &curlRes)

		for _, v := range curlRes.Hits.Hits {
			bytes, _ := json.Marshal(v.Source)
			json.Unmarshal(bytes, &resHitsType)
			switch resHitsType.Type {
			case "waf":
				res = append(res, ApiResWaf(bytes))
			case "vds":
				res = append(res, ApiResVds(bytes))
			case "ids":
				res = append(res, ApiResIds(bytes))
			}
		}
	}

	return ResStruct(curlRes.Hits.Total, res, 0)
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

func EsRes(p *Params, e *tree.Expr) Res {

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
