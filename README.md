# Dockerfile Transpiler

A GoLang transpiler for Dockerfiles.
It leverages of [Participle](https://github.com/alecthomas/participle).

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
