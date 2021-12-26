package main

import (
	"io/ioutil"

	"github.com/alecthomas/repr"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var dockerLexer = lexer.MustSimple([]lexer.Rule{
	{`Ident`, `[a-zA-Z][a-zA-Z_\d]*`, nil},
	{`String`, `"(?:\\.|[^"])*"`, nil},
	{`Float`, `\d+(?:\.\d+)?`, nil},
	{`Punct`, `[][=]`, nil},
	{"comment", `[#;][^\n]*`, nil},
	{"whitespace", `\s+`, nil},
})

type DOCKERFILE struct {
	Start []*From `@@*`
}

type From struct {
	Key   string `@Ident`
	Value *Value `@@`
}

type Value struct {
	Pos    lexer.Position
	String *string ` @String`
}

var parser = participle.MustBuild(&DOCKERFILE{},
	participle.Lexer(dockerLexer),
	participle.Unquote("String"),
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
