# gobba 
gobba is a purely functional, statically typed programming language,
heavily inspired from OCaml, Haskell, Go and Scheme languages. It is based
on Professors Gianluigi Ferrari and Francesca Levi's minicaml interpreter
example..  Development of the OCaml version has dropped, and from now on
the reference implementation of gobba will be written in the Go programming
language.


The goal for gobba is to be a practical language with
built in support for scientific computing, solving some of the problems
that exist in other commonly used interpreted languages like Python and
Javascript. A primary goal is also to offer a compromise between solidity,
ease of learning and the ability to express ideas quickly in the language.

## Goals/Roadmap
- [x] Lexer
- [x] Top Down Operator Precedence (Pratt) parser
- [ ] Formalization of type system
- [ ] Formalization of operational semantics
- [ ] Formalization of effect system 
- [ ] First class support for linear algebra data types and operations, 
- [ ] Go-like VCS based package/module system
- [ ] AST optimizations
- [ ] Aggressive static optimization of chains of linear algebraic operations
- [ ] Standard library
- [ ] Built-in data visualization in the interpreter
- [ ] Optimizing compiler

## Installation
```
go get -u -v github.com/0x0f0f0f/gobba-golang
```

## Changes from 0.4, or the last OCaml version
- Complex numbers literals are created during parsing instead of evaluation
- Introduced allow/deny for effects, including purity
- Comments are now in a C-like syntax
