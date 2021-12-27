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
	{"WorkDirective", `^WORKDIR`, nil},
	{"Directory", `/\w+`, nil},
	{"Word", `-?\w+`, nil},
	{"Options", `--\w+`, nil},
	{`String`, `"(?:\\.|[^"])*"`, nil},
	{"whitespace", `\s`, nil},
	{"multiline", `\\\n`, nil},
	{"EOL", `[\n\r]+`, nil},
})

type DOCKERFILE struct {
	From    *From    `@@`
	Run     []*Run   `@@*`
	Workdir *WorkDir `@@`
}

type From struct {
	Key   string `@BaseIdent`
	Value string `@BaseValue`
}

type WorkDir struct {
	Key   string `@WorkDirective`
	Value string `@Directory`
}

type Run struct {
	Key   string `@RunDirective`
	Value *Value `@@*`
}

type Value struct {
	Exe      string     `@Word @Word @Options @Options`
	Packages []*Package `@@*`
}

type Package struct {
	Dependency string `@Word`
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
