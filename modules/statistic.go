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
	Code int      `code`
	Data AggsData `code`
}
type AggsData struct {
	Type  map[string]int64 `time`
	Time  map[string]int64 `time`
	SrcCY map[string]int64 `srcVY`
	DPort map[string]int64 `dPort`
	Prot  map[string]int64 `prot`
	Os    map[string]int64 `prot`
}

func AggsBody() string {
	return `{"took":1170,"timed_out":false,"_shards":{"total":5,"successful":5,"failed":0},"hits":{"total":4978861,"max_score":0.0,"hits":[]},"aggregations":{"dPort":{"doc_count":4978861,"innerBase":{"doc_count_error_upper_bound":15,"sum_other_doc_count":84581,"buckets":[{"key":80,"doc_count":4891799},{"key":8001,"doc_count":1074},{"key":8080,"doc_count":1054},{"key":8081,"doc_count":195},{"key":81,"doc_count":40},{"key":8443,"doc_count":33},{"key":4800,"doc_count":32},{"key":16464,"doc_count":21},{"key":25,"doc_count":17},{"key":53,"doc_count":15}]}},"srcCY":{"doc_count":4978861,"innerBase":{"doc_count_error_upper_bound":10399,"sum_other_doc_count":1431547,"buckets":[{"key":"中国","doc_count":2521573},{"key":"美国","doc_count":137065},{"key":"俄罗斯","doc_count":124567},{"key":"瑞士","doc_count":113512},{"key":"巴西","doc_count":112793},{"key":"加拿大","doc_count":108463},{"key":"马达加斯加","doc_count":108013},{"key":"捷克","doc_count":107531},{"key":"波兰","doc_count":107197},{"key":"黑山","doc_count":106600}]}},"prot":{"doc_count":4978861,"innerBase":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"tcp","doc_count":4895369},{"key":"udp","doc_count":83492}]}},"os":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"windows","doc_count":602938},{"key":"win32","doc_count":987}]},"time":{"buckets":[{"key_as_string":"2017-08","key":1501545600000,"doc_count":2160},{"key_as_string":"2017-09","key":1504224000000,"doc_count":0},{"key_as_string":"2017-10","key":1506816000000,"doc_count":864},{"key_as_string":"2017-11","key":1509494400000,"doc_count":0},{"key_as_string":"2017-12","key":1512086400000,"doc_count":4923563},{"key_as_string":"2018-01","key":1514764800000,"doc_count":52274}]},"type":{"doc_count_error_upper_bound":0,"sum_other_doc_count":0,"buckets":[{"key":"waf","doc_count":4282581},{"key":"vds","doc_count":603925},{"key":"ids","doc_count":92355}]}}}`
}

func BaseElmt(base Base) map[string]int64 {
	var m = make(map[string]int64)
	for _, v := range base.Buckets {
		m[v.Key] = v.Count
	}

	if 0 != base.Other {
		m["other"] = base.Other
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
