package main

import (
	"testing"

	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/require"
)

func TestFROM(t *testing.T) {
	dockerfile := &DOCKERFILE{}
	err := parser.ParseString("", `
FROM ruby:3.0.3-alpine
WORKDIR /app
`, dockerfile)
	require.NoError(t, err)
	repr.Println(dockerfile)
}

func TestRUN(t *testing.T) {
	dockerfile := &DOCKERFILE{}
	err := parser.ParseString("", `
FROM ruby:3.0.3-alpine

RUN apk add --update --virtual \
  runtime-deps \
  postgresql-client \
  build-base \
  libxml2-dev \
  libxslt-dev \
  nodejs \
  yarn \
  libffi-dev \
  readline \
  build-base \
  postgresql-dev \
  sqlite-dev \
  libc-dev \
  linux-headers \
  readline-dev \
  file \
  imagemagick \
  git \
  tzdata

WORKDIR /app

`, dockerfile)
	require.NoError(t, err)
	repr.Println(dockerfile)
}
