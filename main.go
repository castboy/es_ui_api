package main 

import (
    "gopkg.in/olivere/elastic.v5"
    "context"
    "fmt"
)

var index = "forbin_index_create_auto"

type Tweet struct {
    User string
    Message string
}

func main() {
    ctx := context.Background()

    client, err := elastic.NewClient(elastic.SetURL("http://10.88.1.102:9200"))
    if err != nil {
        fmt.Println("new client err")
    }

    exists, err := client.IndexExists(index).Do(ctx)
    if err != nil {
        fmt.Println("IndexExists() err")
    }
    if !exists {
        _, err = client.CreateIndex(index).Do(ctx)
        if err != nil {
            fmt.Println("create index failed")
        }
    }

    tweet := Tweet{User: "olivere", Message: "Take Five"}
    _, err = client.Index().
        Index("twitter").
        Type("tweet").
        //Id("1").
        BodyJson(tweet).
        Refresh("true").
        Do(ctx)
    if err != nil {
        panic(err)
    } else {
        fmt.Println("create document ok")
    }
}
