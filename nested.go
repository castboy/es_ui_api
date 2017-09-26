package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v5"
)

func main() {
	nested := elastic.NewNestedQuery("xdr", elastic.NewMatchQuery("xdr.Http.ResponseLocation.Signature", "nana"))
	query_string := elastic.NewQueryStringQuery("192.168.1.188").DefaultField("my_xdr")
	query := elastic.NewBoolQuery().Must(nested, query_string)
	query = query.Should(query_string).MinimumShouldMatch("1")
	query = query.MustNot(nested)
	//query := elastic.NewBoolQuery()
	//query = query.Must(elastic.NewMatchQuery("user", "olivere"))
	//query = query.Filter(elastic.NewMatchQuery("account", 1))
	src, err := query.Source()
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	s := string(data)
	fmt.Println(s)
}
