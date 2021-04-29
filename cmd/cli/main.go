package main

import (
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/cli"
)

func init() {
	internal.InitEnv()
}

func main() {
	cli.Execute()
}
