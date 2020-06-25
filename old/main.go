package main

import (
    "bufio"
    "os"
    parser "github.com/0x0f0f0f/gobba-golang/pkg/parser"
)

func check(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    p, err := parser.NewParser()
    check(err)
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        line := scanner.Text()
        expr, err := p.ParseString(line)
        check(err)
        parser.PrintExpression(expr)
    }
}

