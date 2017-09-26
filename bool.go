package main 

import (
    "gopkg.in/olivere/elastic.v5"
    //"context"
    "encoding/json"
    "fmt"
)

func main() {
    query := elastic.NewBoolQuery()
    query = query.Must(elastic.NewMatchQuery("user", "olivere"))
    query = query.Filter(elastic.NewMatchQuery("account", 1))
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
