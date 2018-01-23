package modules

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
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

type CurlBody struct {
	From  int         `json:"from"`
	Size  int         `json:"size"`
	Query interface{} `json:"query"`
	Sort  `json:"sort"`
}

type Sort struct {
	Time Order `json:"TimeAppend"`
}

type Order struct {
	Order string `json:"order"`
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

func (p *Params) ParseEsType() *Params {

	return p
}

func (p *Params) ParseQuery() *Params {
	s, err := base64.StdEncoding.DecodeString(p.Query)
	if nil != err {
		p.Err = ERR_DECODE_BASE64
		Log("ERR", "parseQuery decode_base64 err %s", p.Err)

		return p
	}
	p.Query = strings.ToLower(string(s))
	p.Query = strings.Replace(p.Query, "dprt==未知", "dprt==0", -1)
	p.Query = strings.Replace(p.Query, "dprt=未知", "dprt=0", -1)

	Log("INF", "query is: %s", p.Query)

	return p
}

func (p *Params) ParseFrom() *Params {
	Log("INF", "from is: %d", p.From)

	return p
}

func (p *Params) ParseSize() *Params {
	Log("INF", "size is: %d", p.Size)

	return p
}

func ParseReq(r *http.Request) (*Params, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	var params Params
	json.Unmarshal(body, &params)

	p := params.ParseEsType().ParseQuery().ParseFrom().ParseSize()

	if nil != err || nil != p.Err {
		return nil, ERR_HTTP_REQ
	}

	return p, nil
}

func ApiResWaf(b []byte) interface{} {
	var src SrcWaf
	err := json.Unmarshal(b, &src)
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

		HttpReq: src.Xdr[0].Http.RequestLocation.ReqCont,
		HttpRes: src.Xdr[0].Http.ResponseLocation.ResCont,
	}

	return resWaf
}

func ApiResVds(b []byte) interface{} {
	var src SrcVds
	err := json.Unmarshal(b, &src)
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

func ApiResIds(b []byte) interface{} {
	var src SrcIds
	err := json.Unmarshal(b, &src)
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

func ResLast(result []byte) Res {
	var resHitsType ResHitsType
	var res []interface{}
	var curlRes CurlRes

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

	if curlRes.Hits.Total > MaxResNum {
		curlRes.Hits.Total = 10000
	}

	return ResStruct(curlRes.Hits.Total, res, 0)
}

func Server(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var res Res

	p, err := ParseReq(r)
	if 0 != p.Size {
		if nil != err {
			res = ResStruct(0, nil, WRONG)
		} else {
			e, err := RecoverLineExpr(p)
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
				bytes, err := EsRes(p, e)
				if nil != err {
					Log("ERR", "%s", "CurlEs res.Body")
				} else {
					res = ResLast(bytes)
				}
			}
		}

		s := ApiRes(res)

		io.WriteString(w, *s)

		Log("INF", "query res: %s", *s)
	} else {
		e, _ := RecoverLineExpr(p)
		if nil == e {
			io.WriteString(w, `{"code":400, "data":""}`)
			return
		}

		body := Expr(e)
		i, _ := body.Source()

		b, _ := json.Marshal(i)

		s, err2 := Stat(string(b))
		if nil != err2 {
			io.WriteString(w, `{"code": 400, "data": ""`)
		} else {
			io.WriteString(w, s)
		}

	}

}
