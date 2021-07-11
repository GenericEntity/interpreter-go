# Interpreter in Golang
This is my 2021 summer personal project: to learn how to build an interpreter from scratch: lexer, parser and evaluator. The whole project is based on Thorsten Ball's [excellent book](https://interpreterbook.com/), so I do not take credit for its design or implementation.

The interpreted language is [Monkey](https://interpreterbook.com/#the-monkey-programming-language), designed by Thorsten Ball.

Note: it is not a fully-fledged Monkey interpreter yet (no support for objects and hashes yet)

## Dependencies
* [Golang](https://golang.org/) (I used `v1.13`, but earlier versions may work too)

## How to use
### Read from file
The `-f` flag specifies the path to a script file to interpret. For example, the following command will attempt to interpret the contents of `./source.monkey`
```bash
go run main.go -f ./source.monkey
```
### Interactive session (REPL)
If the `-f` flag is omitted, an interactive session (REPL) will be launched instead.
1. Launch the REPL
```bash
go run main.go
```

2. Then type [Monkey source code](https://interpreterbook.com/#the-monkey-programming-language)
  * E.g. the famous factorial function
```
let fact = fn(n) { if(n == 0) { return 1 } else { return n * fact(n - 1) } }
fact(5)
```

3. Enjoy! This step is mandatory.
