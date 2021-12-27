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
	{"CopyDirective", `^COPY`, nil},
	{"EnvDirective", `^ENV`, nil},
	{"Directory", `/\w+`, nil},
	{"CopyDirectory", `\.\s/\w+/`, nil},
	{"Word", `-?\w+`, nil},
	{"Options", `--\w+`, nil},
	{`String`, `"(?:\\.|[^"])*"`, nil},
	{"whitespace", `\s`, nil},
	{"multiline", `\\\n`, nil},
	{"EOL", `[\n\r]+`, nil},
})

type DOCKERFILE struct {
	From       *From         `@@`
	ComplexRun []*ComplexRun `@@*`
	WorkDir    *WorkDir      `@@`
	Copy       *Copy         `@@`
	Env        *Env          `@@`
	SimpleRun  []*SimpleRun  `@@*`
}

type From struct {
	Key   string `@BaseIdent`
	Value string `@BaseValue`
}

type ComplexRun struct {
	Key   string `@RunDirective`
	Value *Value `@@*`
}

type SimpleRun struct {
	Key   string `@RunDirective`
	Value Words  `@@*`
}

type Words struct {
	Words string `@Word`
}

type WorkDir struct {
	Key   string `@WorkDirective`
	Value string `@Directory`
}

type Copy struct {
	Key   string `@CopyDirective`
	Value string `@CopyDirectory`
}

type Env struct {
	Key   string `@EnvDirective`
	Value string `@Word @Directory`
}

type Value struct {
	Exe      string     `@Word @Word @Options @Options`
	Packages []*Package `@@*`
	Command  string     `| @Word`
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
