package modules

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type AggsRes struct {
	Aggregations `json:"aggregations"`
}

type Aggregations struct {
	Type  Base            `json:"type"`
	Time  TimeBase        `json:"time"`
	Os    Base            `json:"os"`
	Prot  NestedBase      `json"prot"`
	SrcCY NestedBase      `json:"srcCY"`
	Dport DPortNestedBase `json:"dPort"`
}

type DPortNestedBase struct {
	DPortBase `json:"innerBase"`
}
type DPortBase struct {
	Other   int64         `json:"sum_other_doc_count"`
	Buckets []DPortBucket `json:"buckets"`
}
type DPortBucket struct {
	Key   int   `json:"key"`
	Count int64 `json:"doc_count"`
}

type NestedBase struct {
	Base `json:"innerBase"`
}
type Base struct {
	Other   int64    `json:"sum_other_doc_count"`
	Buckets []Bucket `json:"buckets"`
}
type Bucket struct {
	Key   string `json:"key"`
	Count int64  `json:"doc_count"`
}

type TimeBase struct {
	Other   int64        `json:"sum_other_doc_count"`
	Buckets []TimeBucket `json:"buckets"`
}
type TimeBucket struct {
	Key   string `json:"key_as_string"`
	Count int64  `json:"doc_count"`
}

type UiAggs struct {
	Code int      `json:"code"`
	Data AggsData `json:"data"`
}
type AggsData struct {
	Type  map[string]int64 `json:"type"`
	Time  map[string]int64 `json:"time"`
	SrcCY map[string]int64 `json:"srcCY"`
	DPort map[string]int64 `json:"dPort"`
	Prot  map[string]int64 `json:"prot"`
	Os    map[string]int64 `json:"os"`
}

func AggsBody() string {
	return `{"size":0,"aggs":{"type":{"terms":{"field":"Type"}},"time":{"date_histogram":{"field":"TimeAppend","interval":"month","format":"yyyy-MM"}},"srcCY":{"nested":{"path":"Xdr"},"aggs":{"innerBase":{"terms":{"field":"Xdr.Conn.SipInfo.Country"}}}},"dPort":{"nested":{"path":"Xdr"},"aggs":{"innerBase":{"terms":{"field":"Xdr.Conn.Dport"}}}},"prot":{"nested":{"path":"Xdr"},"aggs":{"innerBase":{"terms":{"field":"Xdr.Conn.ProtoAppend"}}}},"os":{"terms":{"field":"Local_platfrom"}}}}`
}

func BaseElmt(base Base) map[string]int64 {
	var m = make(map[string]int64)

	if 0 != base.Other {
		m["other"] = base.Other
	}

	for _, v := range base.Buckets {
		if "" == v.Key {
			m["other"] += v.Count
		} else {
			m[v.Key] = v.Count
		}
	}

	return m
}

func TimeBaseElmt(base TimeBase) map[string]int64 {
	var m = make(map[string]int64)
	for _, v := range base.Buckets {
		m[v.Key] = v.Count
	}

	if 0 != base.Other {
		m["other"] = base.Other
	}

	return m
}

func DportBaseElmt(base DPortBase) map[string]int64 {
	var m = make(map[string]int64)
	for _, v := range base.Buckets {
		m[strconv.Itoa(v.Key)] = v.Count
	}

	if 0 != base.Other {
		m["other"] = base.Other
	}

	return m
}

func UiStat(b []byte) string {
	var res AggsRes

	json.Unmarshal(b, &res)

	uiAggs := UiAggs{
		Code: 200,
		Data: AggsData{
			Type:  BaseElmt(res.Type),
			Time:  TimeBaseElmt(res.Time),
			Os:    BaseElmt(res.Os),
			Prot:  BaseElmt(res.Prot.Base),
			SrcCY: BaseElmt(res.SrcCY.Base),
			DPort: DportBaseElmt(res.Dport.DPortBase),
		},
	}

	b, err := json.Marshal(uiAggs)
	if nil != err {
		fmt.Println("marshal err")
	}

	return string(b)
}

func Stat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b, err := Query(AggsBody())
	if nil != err {
		io.WriteString(w, `{"code": 400, "data": null}`)
		Log("ERR", "statistic res: %s", `{"code": 400, "data": null}`)
	} else {
		io.WriteString(w, UiStat(b))
	}

}
