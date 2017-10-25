package modules

import (
	"bufio"
	"fmt"
	tree "go-study/expr2"
	"io"
	"os"
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

func Load(line string) {
	s := strings.Fields(line)
	var scope tree.Scope

	switch s[2] {
	case "all":
		scope = tree.ScopeAll
	case "zone":
		scope = tree.ScopeAll
	case "object":
		scope = tree.ScopeObject
	default:
		panic("wrong scope type")
	}

	keyword := tree.Keyword{
		Key:   s[0],
		Name:  s[1],
		Scope: scope,
		Type:  1,
	}

	exprBuildinKeywords[keyword.Key] = &keyword

}

func Test(t *testing.T) {
	Read()
	for _, v := range exprBuildinKeywords {
		fmt.Println(*v)
	}
}
