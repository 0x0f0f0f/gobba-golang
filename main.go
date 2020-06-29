package main

import (
	"flag"
	"github.com/0x0f0f0f/gobba-golang/repl"
)

func main() {
	opts := &repl.ReplOptions{}

	flag.BoolVar(&opts.ShowAST, "ast", false, "print the AST before evaluation")
	flag.BoolVar(&opts.ShowTok, "tok", false, "print lexed tokens before parsing")
	flag.BoolVar(&opts.DebugParser, "dparser", false, "enable parser debugging")

	flag.Parse()
	r := repl.New(opts)
	r.Start()
}
