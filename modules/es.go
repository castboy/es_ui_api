package modules

import (
	"encoding/json"
	//	"fmt"

	"gopkg.in/olivere/elastic.v5"
)

func Demo() string {
	nested := elastic.NewNestedQuery("xdr", elastic.NewMatchQuery("xdr.Http.ResponseLocation.Signature", "nana"))
	query_string := elastic.NewQueryStringQuery("192.168.1.188").DefaultField("my_xdr")
	query := elastic.NewBoolQuery().Must(nested, query_string)
	query2 := elastic.NewBoolQuery().Must(query)
	src, err := query2.Source()
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	s := string(data)

	return s
}
