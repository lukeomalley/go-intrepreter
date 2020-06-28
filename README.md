# 🐒 Monkey Lang

Monkey Lang is a "toy" programming language built to help understand how an intrepreted language works under the hood. The intrepreter currently supports functions, higher-order functions, closures and integers and arithmetic.

This was written with the help of [Writing an Intrepreter in Go](https://interpreterbook.com/) by Thorsten Ball.

## 🚀 Getting Started

Currently the project only supports input from the command line repl.

1. Clone the repository: `git clone https://github.com/lukeomalley/go-intrepreter.git`

2. Change into the root directory of the project: `cd go-intrepreter`

3. Start the interactive REPL: `go run main.go`

## ✍️ Sample Mokney Code

Declare a Variable:

```
let x = 5;
```

Define and Apply a Function:

```js
let add = fn(x, y) {
  return x + y;
}

add(5, 5); // => 10
```

Closures:

```js
let newAdder = fn(x) {
  fn(y) { x + y };
};

let addTwo = newAdder(2);

addTwo(2); // => 4
```

Nth Fibonacci Number:

```js
let fib = fn(n) {
  if (n < 1) {
    return n;
  }

  return fib(n - 1) + fib (n - 2);
};
```

## 🛠 How it Works

### Lexing

### Parsing

### Evaluating

### Repl

## 🛰 Status

4.2 : p 154 Strings
