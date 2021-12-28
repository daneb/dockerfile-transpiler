# Dockerfile Transpiler

A GoLang transpiler for Dockerfiles.
It leverages of [Participle](https://github.com/alecthomas/participle).

The transpiler relies on distinct lexing and parsing phases. The lexer takes raw bytes and produces tokens which the parser consumes. The parser transforms these tokens into Go values.

A conversion of a Dockerfile into Go values (abstract syntax tree).

## Why?

Consider trying to build security automation around Dockerfiles. The system would review multiple repositories and their associated Dockerfiles, trying to look for mis-configurations or company policy violations. However pattern matching (regex) would only provide a certain level of aid, and the developer would have to understand a vast array of permutations in configuration.

Using the process of a transpiler, we can leverage off tokenization (lexer) and parsing (parsing) to build an AST that is easy to work with in your language of choice and at the same time give context.

## Usage

Supported Dockerfile directives are: [From](https://docs.docker.com/engine/reference/builder/#from), [Run](https://docs.docker.com/engine/reference/builder/#run), [WorkDir](https://docs.docker.com/engine/reference/builder/#workdir), [Env](https://docs.docker.com/engine/reference/builder/#env), [EntryPoint](https://docs.docker.com/engine/reference/builder/#entrypoint), [Cmd](https://docs.docker.com/engine/reference/builder/#cmd) and [Expose](https://docs.docker.com/engine/reference/builder/#expose).

### Grammar

```go
type DOCKERFILE struct {
	From       *From         `@@`
	ComplexRun []*ComplexRun `@@*`
	WorkDir    *WorkDir      `@@`
	Copy       *Copy         `@@`
	Env        *Env          `@@`
	SimpleRun  []*SimpleRun  `@@*`
	EntryPoint *EntryPoint   `@@`
	Cmd        *Cmd          `@@`
	Expose     *Expose       `@@`
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

type Expose struct {
	Key   string `@ExposeDirective`
	Value string `@Port`
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
```

### Lexer Rules

```go
var dockerLexer = lexer.MustSimple([]lexer.Rule{
	{"BaseIdent", `^FROM`, nil},
	{"BaseValue", `[a-zA-Z]*:\d+\.\d+\.\d+\-[a-zA-Z]*`, nil},
	{"RunDirective", `^RUN`, nil},
	{"WorkDirective", `^WORKDIR`, nil},
	{"CopyDirective", `^COPY`, nil},
	{"EnvDirective", `^ENV`, nil},
	{"EntryPointDirective", `^ENTRYPOINT`, nil},
	{"CmdDirective", `^CMD`, nil},
	{"ExposeDirective", `^EXPOSE`, nil},
	{"Directory", `/\w+`, nil},
	{"CopyDirectory", `\.\s/\w+/`, nil},
	{"Word", `-?\w+`, nil},
	{"AppToRun", `\["\w+/\w+"\]`, nil},
	{"Options", `--\w+`, nil},
	{"String", `"(?:\\.|[^"])*"`, nil},
	{"StringArgs", `\[[\s?"-?\w?",?]+\s?]`, nil},
	{"Port", `\s[0-9]+`, nil},
	{"whitespace", `\s`, nil},
	{"multiline", `\\\n`, nil},
	{"EOL", `[\n\r]+`, nil},
})
```

## Test

```sh
# go test
&main.DOCKERFILE{
  From: &main.From{
    Key: "FROM",
    Value: "ruby:3.0.3-alpine",
  },
  ComplexRun: []*main.ComplexRun{
    {
      Key: "RUN",
      Value: &main.Value{
        Exe: "apkadd--update--virtual",
        Packages: []*main.Package{
          {
            Dependency: "nodejs",
          },
          {
            Dependency: "yarn",
          },
          {
            Dependency: "readline",
          },
          {
            Dependency: "libc",
          },
          {
            Dependency: "-dev",
          },
          {
            Dependency: "file",
          },
          {
            Dependency: "imagemagick",
          },
          {
            Dependency: "git",
          },
          {
            Dependency: "tzdata",
          },
        },
      },
    },
  },
  WorkDir: &main.WorkDir{
    Key: "WORKDIR",
    Value: "/app",
  },
  Copy: &main.Copy{
    Key: "COPY",
    Value: ". /app/",
  },
  Env: &main.Env{
    Key: "ENV",
    Value: "BUNDLE_PATH/gems",
  },
  SimpleRun: []*main.SimpleRun{
    {
      Key: "RUN",
      Value: main.Words{
        Words: "install",
      },
    },
    {
      Key: "RUN",
      Value: main.Words{
        Words: "bundler",
      },
    },
    {
      Key: "RUN",
      Value: main.Words{
        Words: "install",
      },
    },
  },
  EntryPoint: &main.EntryPoint{
    Key: "ENTRYPOINT",
    Value: "[\"bin/rails\"]",
  },
  Cmd: &main.Cmd{
    Key: "CMD",
    Arguments: "[ \"s\",\"-b\",\"0.0.0.0\" ]",
  },
  Expose: &main.Expose{
    Key: "EXPOSE",
    Value: " 3000",
  },
}
PASS
ok  	dockerfiletranspiler	0.140s

```

## References

1. [Dockerfile specification](https://docs.docker.com/engine/reference/builder/)
