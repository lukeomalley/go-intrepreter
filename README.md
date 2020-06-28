# Monkey Lang

An intrepreter for the monkey programing language written in Go.

## 🚀 Getting Started

Currently the project only supports input from the command line repl.

1. Clone the repository: `git clone https://github.com/lukeomalley/go-intrepreter.git`

2. Change into the root directory of the project: `cd go-intrepreter`

3. Start the interactive REPL: `go run main.go`

## ✍️ Sample Mokney Code

Nth Fibonacci Number:

```
let fib = fn(n) {
  if (n < 1) {
    return n;
  }

  return fib(n - 1) + fib (n - 2);
};
```

Closures:

```
let newAdder = fn(x) {
  fn(y) { x + y };
};

let addTwo = newAdder(2);

addTwo(2); // => 4
```

## 🛠 How it Works

### Lexing

### Parsing

### Evaluating

### Repl

## 🛰 Status

4.2 : p 154 Strings
