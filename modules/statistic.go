package modules

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type AggsRes struct {
	Aggregations `json:"aggregations"`
}

type Aggregations struct {
	Type  Base            `json:"type"`
	Time  TimeBase        `json:"time"`
	Os    Base            `json:"os"`
	Prot  NestedBase      `json:"prot"`
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
	Type  []map[string]int64 `json:"type"`
	Time  []map[string]int64 `json:"time"`
	SrcCY []map[string]int64 `json:"srcCY"`
	DPort []map[string]int64 `json:"dPort"`
	Prot  []map[string]int64 `json:"prot"`
	Os    []map[string]int64 `json:"os"`
}

func AggsBody() string {
	return `"aggs":{"type":{"terms":{"field":"Type"}},"time":{"date_histogram":{"field":"TimeAppend","interval":"month","format":"yyyy-MM"}},"srcCY":{"nested":{"path":"Xdr"},"aggs":{"innerBase":{"terms":{"field":"Xdr.Conn.SipInfo.Country", "size": 250}}}},"dPort":{"nested":{"path":"Xdr"},"aggs":{"innerBase":{"terms":{"field":"Xdr.Conn.Dport", "size":65536}}}},"prot":{"nested":{"path":"Xdr"},"aggs":{"innerBase":{"terms":{"field":"Xdr.Conn.ProtoAppend", "size": 500}}}},"os":{"terms":{"field":"Local_platfrom"}}}`
}

func BaseElmt(base Base) []map[string]int64 {
	var s = make([]map[string]int64, 0)

	for _, v := range base.Buckets {
		var m = make(map[string]int64)
		if "" == v.Key || "亚太互联网络信息中心(未知)" == v.Key {
			m["未知"] += v.Count
			s = append(s, m)
		} else {
			m[v.Key] = v.Count
			s = append(s, m)
		}
	}

	return s
}

func TimeBaseElmt(base TimeBase) []map[string]int64 {
	var s = make([]map[string]int64, 0)

	for _, v := range base.Buckets {
		var m = make(map[string]int64)
		m[v.Key] = v.Count
		s = append(s, m)
	}

	return s
}

func DportBaseElmt(base DPortBase) []map[string]int64 {
	var s = make([]map[string]int64, 0)

	for _, v := range base.Buckets {
		var m = make(map[string]int64)
		if 0 == v.Key {
			m["未知"] += v.Count
			s = append(s, m)
		} else {
			m[strconv.Itoa(v.Key)] = v.Count
			s = append(s, m)
		}
	}

	return s
}

func UiStat(b []byte) string {
	var res AggsRes
	var s string

	json.Unmarshal(b, &res)
	fmt.Println(res)

	if 0 == len(res.Type.Buckets) {
		s = `{"code": 200, "data": ""}`
	} else {
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
		byte, _ := json.Marshal(uiAggs)
		s = string(byte)
	}

	return s
}

func Stat(query string) (string, error) {
	aggsBody := fmt.Sprintf(`{"size": 0, "query": %s, %s}`, query, AggsBody())

	b, err := Query(aggsBody)
	if nil != err {
		Log("ERR", "statistic res: %s", `{"code": 400, "data": {"type":[], "time":[], "srcCY":[], "dPort":[], "prot":[], "os":[]}}`)
		return `{"code": 400, "data": null}`, err
	}

	return UiStat(b), nil

}
