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
	{"EntryPointDirective", `^ENTRYPOINT`, nil},
	{"CmdDirective", `^CMD`, nil},
	{"Directory", `/\w+`, nil},
	{"CopyDirectory", `\.\s/\w+/`, nil},
	{"Word", `-?\w+`, nil},
	{"AppToRun", `\["\w+/\w+"\]`, nil},
	{"Options", `--\w+`, nil},
	{`String`, `"(?:\\.|[^"])*"`, nil},
	{`StringArgs`, `\[[\s?"-?\w?",?]+\s?]`, nil},
	{`Ip`, ` (\b25[0-5]|\b2[0-4][0-9]|\b[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`, nil},
	{`Char`, `\w`, nil},
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
	EntryPoint *EntryPoint   `@@`
	Cmd        *Cmd          `@@`
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

type EntryPoint struct {
	Key   string `@EntryPointDirective`
	Value string `@AppToRun`
}

type Cmd struct {
	Key       string `@CmdDirective`
	Arguments string `@StringArgs`
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

type Argument struct {
	Arg  *string     `	@String`
	Ip   *string     `| @Ip`
	List []*Argument `| "[" (@@)? "]"`
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
