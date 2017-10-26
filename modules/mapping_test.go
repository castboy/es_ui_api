package modules

import (
	"bufio"
	"fmt"
	tree "go-study/expr2"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

var exprBuildinKeywords = map[string]*tree.Keyword{}

func Read() {
	fi, err := os.Open("mapping")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		Load(string(a))
	}
}

func parseType(str string) (t int) {
	t, err := strconv.Atoi(str)
	if nil == err && (-1 < t) && (t < 4) {
		return t
	}

	log.Fatal("mapping file err: wrong type, must: -1 < type < 4")
	return -1
}

func Load(line string) {
	var scope tree.Scope
	s := strings.Fields(line)

	switch s[2] {
	case "all":
		scope = tree.ScopeAll
	case "zone":
		scope = tree.ScopeZone
	case "object":
		scope = tree.ScopeObject
	default:
		panic("wrong scope type")
	}

	keyword := tree.Keyword{
		Key:   s[0],
		Name:  s[1],
		Scope: scope,
		Type:  parseType(s[3]),
	}

	exprBuildinKeywords[keyword.Key] = &keyword

}

func Test(t *testing.T) {
	Read()
	for _, v := range exprBuildinKeywords {
		fmt.Println(*v)
	}
}
