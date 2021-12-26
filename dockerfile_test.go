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
`, dockerfile)
	require.NoError(t, err)
	repr.Println(dockerfile)
}
