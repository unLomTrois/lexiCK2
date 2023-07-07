# What is this?
Golang implementation of linting PDXScript files for Crusader Kings 3

## Structure
It is composed of three parts:
- lexer (tokenizer)
- parser (LL(1) parser)
- linter

Lexer creates stream of tokens that is consumed by parser, catches lexical errors, like unknown tokens (e.g. you can't write `!=`)

Parser makes AST (Abstract Syntax Tree), catches syntax errors (e.g. not closed curly brace)

Linter takes AST and rewrites the file by set linting rules
