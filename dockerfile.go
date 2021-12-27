package main

import (
	"io/ioutil"

	"github.com/alecthomas/repr"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var dockerLexer = lexer.MustSimple([]lexer.Rule{
	{"BaseIdent", `^FROM`, nil},
	{"BaseValue", `[a-zA-Z]*:\d+\.\d+\.\d+\-[a-zA-Z]*`, nil},
	{"RunDirective", `^RUN`, nil},
	{"Word", `\w+`, nil},
	{"Options", `--\w+`, nil},
	{"whitespace", `\s{1}`, nil},
	{`String`, `"(?:\\.|[^"])*"`, nil},
})

type DOCKERFILE struct {
	From    *From  `@@`
	Runners []*Run `@@*`
}

type From struct {
	Key   string `@BaseIdent`
	Value string `@BaseValue`
}

type Run struct {
	Key string `@RunDirective`
	Cmd *Value `@@`
}

type Value struct {
	Exe     string `@Word`
	Action  string `@Word`
	Options string `@Options @Options`
}

var parser = participle.MustBuild(&DOCKERFILE{},
	participle.Lexer(dockerLexer),
	participle.Unquote(),
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
