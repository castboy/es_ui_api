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
)

var exprBuildinKeywords = map[string]*tree.Keyword{}
var abbreviations = make([]string, 0)

func Read(file string) []string {
	var fi *os.File
	var err error

	fi, err = os.Open(file)
	if err != nil {
		log.Fatal("can not find register_keyword file")
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		abbreviations = append(abbreviations, string(a))
	}

	return abbreviations
}

func parseType(str string) (t int) {
	t, err := strconv.Atoi(str)
	if nil == err && (-1 < t) && (t < 4) {
		return t
	}

	log.Fatal("mapping file err: wrong type, must: -1 < type < 4")
	return -1
}

func Load(ab []string) {
	var scope tree.Scope

	for _, v := range ab {
		s := strings.Fields(v)

		switch s[2] {
		case "all":
			scope = tree.ScopeAll
		case "zone":
			scope = tree.ScopeZone
		case "object":
			scope = tree.ScopeObject
		default:
			panic(PANIC_SEARCH_SCOPE)
		}

		keyword := tree.Keyword{
			Key:   s[0],
			Name:  s[1],
			Scope: scope,
			Type:  parseType(s[3]),
		}

		exprBuildinKeywords[keyword.Key] = &keyword
	}

	for _, v := range exprBuildinKeywords {
		fmt.Println(*v)
	}
}

func RegisterKeyword(file string) {
	abs := Read(file)
	Load(abs)

	for _, v := range exprBuildinKeywords {
		tree.RegisterKeyword(*v)
	}
}
