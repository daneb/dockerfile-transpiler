package main

import (
	"testing"

	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/require"
)

func TestTranspile(t *testing.T) {
	dockerfile := &DOCKERFILE{}
	err := parser.ParseString("", `
FROM ruby:3.0.3-alpine

RUN apk add --update --virtual \
  nodejs \
  yarn \
  readline \
  libc-dev \
  file \
  imagemagick \
  git \
  tzdata

WORKDIR /app
COPY . /app/

ENV BUNDLE_PATH /gems
RUN yarn install
RUN gem install bundler
RUN bundle install

`, dockerfile)
	require.NoError(t, err)
	repr.Println(dockerfile)
}
