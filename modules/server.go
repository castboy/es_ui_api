package modules

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/olivere/elastic.v5"
)

type Params struct {
	T     AlertType
	Query string
	From  int
	Size  int
}

type ResApi []interface{}

func (p *Params) ParseEsType(r *http.Request) *Params {
	if val, ok := r.Form["type"]; ok {
		switch val[0] {
		case EsType[Waf]:
			p.T = Waf
		case EsType[Vds]:
			p.T = Vds
		case EsType[Ids]:
			p.T = Ids
		default:
			panic(PANIC_UNKNOW_ALERT)
		}
	}

	return p
}

func (p *Params) ParseQuery(r *http.Request) *Params {
	if val, ok := r.Form["query"]; ok {
		s, err := base64.StdEncoding.DecodeString(val[0])
		if nil != err {

		}

		p.Query = string(s)
	}

	return p
}

func (p *Params) ParseFrom(r *http.Request) *Params {
	if val, ok := r.Form["from"]; ok {
		from, err := strconv.Atoi(val[0])
		if nil != err {
			fmt.Println("param from err")
		}

		p.From = from
	}

	return p
}

func (p *Params) ParseSize(r *http.Request) *Params {
	if val, ok := r.Form["size"]; ok {
		size, err := strconv.Atoi(val[0])
		if nil != err {
			fmt.Println("param size err")
		}

		p.Size = size
	}

	return p
}

func ParseParams(r *http.Request) Params {
	var p Params

	err := r.ParseForm()
	if nil != err {

	} else {
		p.ParseEsType(r).ParseQuery(r).ParseFrom(r).ParseSize(r)
	}

	return p
}

func ApiResWaf(hit *elastic.SearchHit) interface{} {
	var src WafSource
	err := json.Unmarshal(*hit.Source, &src)
	if nil != err {
		fmt.Println("Unmarshal WafSource err")
	}

	resWaf := ApiWafRes{
		WafAlert: WafAlert{
			Client:    src.Client,
			Rev:       src.Rev,
			Msg:       src.Msg,
			Attack:    src.Attack,
			Severity:  src.Severity,
			Maturity:  src.Maturity,
			Accuracy:  src.Accuracy,
			Hostname:  src.Hostname,
			Uri:       src.Uri,
			Unique_id: src.Unique_id,
			Ref:       src.Ref,
			Tags:      src.Tags,
			Rule:      src.Rule,
			Version:   src.Version,
		},
		Time:        src.Xdr[0].Time,
		Dest_ip:     src.Xdr[0].Conn.Dip,
		Dest_port:   src.Xdr[0].Conn.Dport,
		Src_ip:      src.Xdr[0].Conn.Sip,
		Src_ip_info: src.Xdr[0].Conn.SipInfo,
	}

	return resWaf
}

func ApiResVds(hit *elastic.SearchHit) interface{} {
	var src ApiVdsRes

	return src
}

func ApiResIds(hit *elastic.SearchHit) interface{} {
	var src ApiIdsRes

	return src
}

func (res *ResApi) ResSingle(hits *elastic.SearchHits) {
	for _, hit := range hits.Hits {
		switch hit.Type {
		case EsType[Waf]:
			*res = append(*res, ApiResWaf(hit))
		case EsType[Vds]:
			*res = append(*res, ApiResVds(hit))
		case EsType[Ids]:
			*res = append(*res, ApiResIds(hit))
		default:
			panic(PANIC_UNKNOW_ALERT)
		}
	}
}

func (res *ResApi) ResMulti(hits *elastic.SearchHits, t AlertType) {
	switch t {
	case Waf:
		res.MultiRes(hits, ApiResWaf)
	case Vds:
		res.MultiRes(hits, ApiResVds)
	case Ids:
		res.MultiRes(hits, ApiResIds)
	default:
		panic(PANIC_UNKNOW_ALERT)
	}
}

func (res *ResApi) MultiRes(hits *elastic.SearchHits, f func(hit *elastic.SearchHit) interface{}) {
	for _, hit := range hits.Hits {
		*res = append(*res, f(hit))
	}
}

func ApiRes(hits *elastic.SearchHits, t AlertType) string {
	var res ResApi

	if t == Multi {
		res.ResSingle(hits)
	} else {
		res.ResMulti(hits, t)
	}

	bytes, err := json.Marshal(res)
	if nil != err {
		fmt.Println("Marshal hits err")
	}

	return string(bytes)
}

func Server(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	p := ParseParams(r)
	hits := EsRes(p)
	s := ApiRes(hits, p.T)

	io.WriteString(w, s)
}
