package main

import (
    "io/ioutil"
    "fmt"
    "encoding/json"
    "net/http"
)

type Params struct {
    And []string `json:"and"`
    Or  []string `json:"or"`
    Not []string `json:"not"`
}

type Match struct {
    All []string `json:"_all"`    
}

type Must struct {
    Match2 Match `json:"match"`   
}

type Should struct {
    Matches []Match `json:"should"`
}

type Bool struct {
    Must2 Must `json:"must"`
    MustNot Must `json:"must_not"` 
}

func hello(rw http.ResponseWriter, req *http.Request) {
    var obj Params 
    body, _ := ioutil.ReadAll(req.Body)
    json.Unmarshal(body, &obj)

    for _, val := range obj.Or {
        match := Match{All: }     
    }
   
    h := Bool{Must2: Must{Match2: Match{All: obj.And}}, MustNot: Must{Match2: Match{All: obj.Not}}}

    b, _ := json.Marshal(h)
    fmt.Println(string(b))
}


func main() {
    http.HandleFunc("/", hello) 
    http.ListenAndServe(":8080", nil)
}
