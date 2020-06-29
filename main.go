package main

import (
	"flag"
	"github.com/0x0f0f0f/gobba-golang/repl"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

func main() {
	opts := &repl.ReplOptions{}

	flag.BoolVar(&opts.ShowAST, "ast", false, "print the AST before evaluation")
	flag.BoolVar(&opts.ShowTok, "tok", false, "print lexed tokens before parsing")
	flag.BoolVar(&opts.DebugParser, "dparser", false, "enable parser debugging")

	flag.Parse()
	r := repl.New(opts)

	// Intercept sighup
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGQUIT)
	go func() {
		for {
			s := <-sigc
			switch s {
			case syscall.SIGQUIT:
				pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
				os.Exit(1)
			default:
				break
			}
		}
	}()

	r.Start()
}
