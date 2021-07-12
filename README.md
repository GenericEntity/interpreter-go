# Interpreter in Golang
This is my 2021 summer personal project: to learn how to build an interpreter from scratch: lexer, parser and evaluator. The whole project is based on Thorsten Ball's [excellent book](https://interpreterbook.com/), so I do not take credit for its design or implementation (except a few minor add-ons).

The interpreted language is [Monkey](https://interpreterbook.com/#the-monkey-programming-language), designed by Thorsten Ball.

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

## License
Note: A lot of the code in this repository follows the code presented in the book very closely. The main differences are a slightly nicer testing framework, a flag to interpret from a file, and support for escape characters in strings.

Thorsten Ball licensed his code under the MIT license. His license file has been reproduced here as the file `LICENSES/LICENSE_THORSTEN_BALL`.
