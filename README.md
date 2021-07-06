# Interpreter in Golang
This is my 2021 summer personal project: to learn how to build an interpreter from scratch: lexer, parser and evaluator. The whole project is based on Thorsten Ball's [excellent book](https://interpreterbook.com/), so I do not take credit for its design or implementation.

## Dependencies
[Golang](https://golang.org/)

## How to use
1. Launch the interpreter
```bash
go run main.go
```

2. Then type [Monkey source code](https://interpreterbook.com/#the-monkey-programming-language)
  * E.g. the famous factorial function
```
let fact = fn(n) { if(n == 0) { return 1 } else { return n * fact(n - 1) } }
fact(5)
```
  * Note: it is not a fully-fledged Monkey interpreter yet (no support for strings, arrays, objects, hashes yet)

3. Enjoy!
