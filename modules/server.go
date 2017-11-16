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
	Err   error
}

type Res struct {
	Total int64       `json:"total"`
	Hits  interface{} `json:"hits"`
	Code  ResCode     `json:"code"`
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
		case EsType[Multi]:
			p.T = Multi
		default:
			panic(PANIC_UNKNOW_ALERT)
		}
	} else {
		p.T = Multi
	}

	return p
}

func (p *Params) ParseQuery(r *http.Request) *Params {
	if val, ok := r.Form["query"]; ok {
		s, err := base64.StdEncoding.DecodeString(val[0])
		if nil != err {
			p.Err = ERR_DECODE_BASE64

			return p
		}
		p.Query = string(s)

		return p
	}
	p.Err = ERR_HTTP_REQ

	return p
}

func (p *Params) ParseFrom(r *http.Request) *Params {
	if val, ok := r.Form["from"]; ok {
		from, err := strconv.Atoi(val[0])
		if nil != err {
			p.Err = ERR_HTTP_REQ

			return p
		}
		p.From = from

		return p
	}
	p.Err = ERR_HTTP_REQ

	return p
}

func (p *Params) ParseSize(r *http.Request) *Params {
	if val, ok := r.Form["size"]; ok {
		size, err := strconv.Atoi(val[0])
		if nil != err {
			p.Err = ERR_HTTP_REQ

			return p
		}
		p.Size = size

		return p
	}
	p.Err = ERR_HTTP_REQ

	return p
}

func ParseReq(r *http.Request) (*Params, error) {
	err := r.ParseForm()
	p := ParseParams(r)
	if nil != err || nil != p.Err {
		return nil, ERR_HTTP_REQ
	}

	return p, nil
}

func ParseParams(r *http.Request) *Params {
	var p Params

	return p.ParseEsType(r).ParseQuery(r).ParseFrom(r).ParseSize(r)
}

func ApiResWaf(hit *elastic.SearchHit) interface{} {
	var src SrcWaf
	err := json.Unmarshal(*hit.Source, &src)
	if nil != err {
		fmt.Println("Unmarshal WafSource err")
	}

	resWaf := ResWaf{
		Client:    src.Client,
		Rev:       src.Rev,
		Msg:       src.Msg,
		Attack:    src.Attack,
		Severity:  src.SeverityAppend,
		Maturity:  src.Maturity,
		Accuracy:  src.Accuracy,
		Hostname:  src.Hostname,
		Uri:       src.Uri,
		Unique_id: src.Unique_id,
		Ref:       src.Ref,
		Tags:      src.Tags,
		Rule:      src.Rule,
		Version:   src.Version,
		Time:      src.Xdr[0].TimeAppend,
		ConnInfo: ConnInfo{
			Proto:        src.Xdr[0].Conn.ProtoAppend,
			Dest_ip:      src.Xdr[0].Conn.Dip,
			Dest_port:    src.Xdr[0].Conn.Dport,
			Dest_ip_info: src.Xdr[0].Conn.DipInfo,
			Src_ip:       src.Xdr[0].Conn.Sip,
			Src_port:     src.Xdr[0].Conn.Sport,
			Src_ip_info:  src.Xdr[0].Conn.SipInfo,
		},
		Operators: "",
		Type:      "waf",
	}

	return resWaf
}

func ApiResVds(hit *elastic.SearchHit) interface{} {
	var src SrcVds
	err := json.Unmarshal(*hit.Source, &src)
	if nil != err {
		fmt.Println("Unmarshal WafSource err")
	}

	resVds := ResVds{
		Subfile:          src.Subfile,
		Threatname:       src.Threatname,
		Local_threatname: src.Local_threatname,
		Attack:           src.Attack,
		Local_platfrom:   src.Local_platfrom,
		Local_vname:      src.Local_vname,
		Severity:         src.SeverityAppend,
		Local_enginetype: src.Local_enginetype,
		Local_logtype:    src.Local_logtype,
		Local_engineip:   src.Local_engineip,
		Time:             src.Xdr[0].TimeAppend,
		HttpUrl:          src.Xdr[0].Http.Url,
		Filepath:         src.Xdr[0].App.File,
		ConnInfo: ConnInfo{
			Proto:        src.Xdr[0].Conn.ProtoAppend,
			Dest_ip:      src.Xdr[0].Conn.Dip,
			Dest_port:    src.Xdr[0].Conn.Dport,
			Dest_ip_info: src.Xdr[0].Conn.DipInfo,
			Src_ip:       src.Xdr[0].Conn.Sip,
			Src_port:     src.Xdr[0].Conn.Sport,
			Src_ip_info:  src.Xdr[0].Conn.SipInfo,
		},
		Type: "vds",
	}

	return resVds
}

func ApiResIds(hit *elastic.SearchHit) interface{} {
	var src SrcIds
	err := json.Unmarshal(*hit.Source, &src)
	if nil != err {
		fmt.Println("Unmarshal WafSource err")
	}

	resIds := ResIds{
		Time:     src.Xdr[0].TimeAppend,
		Attack:   src.Attack,
		Details:  "",
		Severity: src.SeverityAppend,
		Engine:   src.Engine,
		ConnInfo: ConnInfo{
			Src_ip:       src.Xdr[0].Conn.Sip,
			Src_port:     src.Xdr[0].Conn.Sport,
			Src_ip_info:  src.Xdr[0].Conn.SipInfo,
			Dest_ip:      src.Xdr[0].Conn.Dip,
			Dest_port:    src.Xdr[0].Conn.Dport,
			Dest_ip_info: src.Xdr[0].Conn.DipInfo,
			Proto:        src.Xdr[0].Conn.ProtoAppend,
		},
		Type: "ids",
	}

	return resIds
}

func (res *ResApi) ResMulti(hits *elastic.SearchHits) {
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

func (res *ResApi) ResSingle(hits *elastic.SearchHits, p Params) {
	switch p.T {
	case Waf:
		res.SingleRes(hits, ApiResWaf)
	case Vds:
		res.SingleRes(hits, ApiResVds)
	case Ids:
		res.SingleRes(hits, ApiResIds)
	default:
		panic(PANIC_UNKNOW_ALERT)
	}
}

func (res *ResApi) SingleRes(hits *elastic.SearchHits, f func(hit *elastic.SearchHit) interface{}) {
	for _, hit := range hits.Hits {
		*res = append(*res, f(hit))
	}
}

func Hits(hits *elastic.SearchHits, p Params) ResApi {
	var res ResApi

	if p.T == Multi {
		res.ResMulti(hits)
	} else {
		res.ResSingle(hits, p)
	}

	return res
}

func ResStruct(total int64, hits interface{}, code ResCode) Res {
	return Res{Total: total, Hits: hits, Code: code}
}

func ApiRes(i interface{}) *string {
	bytes, err := json.Marshal(i)
	if nil != err {
		fmt.Println("Marshal hits err")
	}

	res := string(bytes)
	return &res
}

func Server(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var res Res

	p, err := ParseReq(r)
	if nil != err {
		res = ResStruct(0, nil, WRONG)
	} else {
		e, err := RecoverLineExpr(p)
		fmt.Println("e:", e)
		fmt.Println("err:", err)
		if nil == e {
			switch err {
			case ErrExprType[NOT_CLOSE_QUOTES_SINGLE]:
				res = ResStruct(0, nil, NOT_CLOSE_QUOTES_SINGLE)
			case ErrExprType[NOT_CLOSE_QUOTES_DOUBLE]:
				res = ResStruct(0, nil, NOT_CLOSE_QUOTES_DOUBLE)
			case ErrExprType[NOT_FOUND_QUOTES_SINGLE_NEXT]:
				res = ResStruct(0, nil, NOT_FOUND_QUOTES_SINGLE_NEXT)
			case ErrExprType[TOKEN_IS_NOT_EXPRESS]:
				res = ResStruct(0, nil, TOKEN_IS_NOT_EXPRESS)
			case ErrExprType[NOT_CLOSE_PARENTHESIS]:
				res = ResStruct(0, nil, NOT_CLOSE_PARENTHESIS)
			default:
				res = ResStruct(0, nil, ERR_EXPRESS)
			}
		} else {
			hits := EsRes(p, e)
			res = ResStruct(hits.TotalHits, Hits(hits, *p), SUCCESS)
		}
	}

	s := ApiRes(res)

	io.WriteString(w, *s)
}
