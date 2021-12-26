package main

import (
	"io/ioutil"

	"github.com/alecthomas/repr"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var dockerLexer = lexer.MustSimple([]lexer.Rule{
	{"BaseIdent", `^FROM`, nil},
	{"BaseValue", `[a-zA-z]*:\d+\.\d+\.\d+\-[a-zA-z]*`, nil},
	{"whitespace", `\s{1}`, nil},
})

type DOCKERFILE struct {
	From *From `@@`
}

type From struct {
	Key   string `@BaseIdent`
	Value string `@BaseValue`
}

var parser = participle.MustBuild(&DOCKERFILE{},
	participle.Lexer(dockerLexer),
)

func main() {
	dockerfile := &DOCKERFILE{}
	content, _ := ioutil.ReadFile("Dockerfile")

	err := parser.ParseString("", string(content), dockerfile)
	if err != nil {
		panic(err)
	}
	repr.Println(dockerfile, repr.Indent("  "), repr.OmitEmpty(true))

}
