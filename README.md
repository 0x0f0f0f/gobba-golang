# gobba 
gobba is a purely functional interpreted programming
language, heavily inspired from the OCaml, Haskell and Scheme languages. It
is based on Professors Gianluigi Ferrari and Francesca Levi's minicaml
interpreter example. At first, gobba was written in OCaml.
Development of the OCaml version has dropped, and from now on 
the reference implementation of gobba will be written in the Go programming language.


The goal for gobba is to be a practical language with
built in support for scientific computing, solving some of the problems
that exist in other commonly used interpreted languages like Python and
Javascript. A primary goal is also to offer a compromise between solidity,
ease of learning and the ability to express ideas quickly in the language.

## Goals/Roadmap
- [] Painless Go-like, VCS based package and module system
- [] Formalized type and effect system and operational semantics
- [] Type and Memory safety
- [] First class support for linear algebra data types and operations, 
  abstracting away the physical representation
- [] Built-in support for probability distributions
- [] Aggressive static optimization of chains of linear algebraic operations
- [] Built-in data visualization for interpreted programs
- [] Optimizing compiler

## Installation
```
go get -u -v github.com/0x0f0f0f/gobba-golang
```

## Changes from 0.4, or the last OCaml version
- Function application syntax: `fun(a1, a2, a3)` instead of `fun a1 a2 a3`.
- Hand written Top Down Operator Precedence (Pratt) parser
- Complex numbers literals are created during parsing instead of evaluation
- Introduced allow/deny for effects, including purity
- Comments are now in a C-like syntax
