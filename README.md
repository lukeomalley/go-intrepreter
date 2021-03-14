# 🐒 Monkey Lang

Monkey Lang is a programming language built to help understand how an programming languages work under the hood. The compiler currently supports functions, higher-order functions, closures, strings, integers, arrays, hashes, and integer arithmetic.

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

  return fib(n - 1) + fib(n - 2);
};
```

Array Map Function

```js
let map = fn(arr, f) {
  let iter = fn(arr, accumulated) {
    if (len(arr) == 0) {
      accumulated
    } else {
      iter(rest(arr), push(accumulated, f(first(arr))));
    }
  };

  return iter(arr, []);
};
```

Array Reduce Function

```js
let reduce = fn(arr, initial, f) {
  let iter = fn(arr, result) {
    if (len(arr) == 0) {
      result
    } else {
      iter(rest(arr), f(result, first(arr)));
    }
  };

  iter(arr, initial);
};
```

## 🛠 How it Works

### Lexing

### Parsing

### Evaluating

### Repl

## 🛰 Status

Chapter 8: Built-In Functions (pg. 216)
