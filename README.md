# The Monkey Programming Language

This repo contains two implementations of the monkey programming language. Both implementations share the same frontends (lexer + parser), however their backends differ. The first backend is a tree-walking interpreter, the second is a single pass bytecode compiler with a companion virtual machine.

Ultimately whether the interpreter or compiler+vm is running, the monkey code eventually gets executed in native Go.

## To Start the REPL
`go run main.go`

The REPL is configured to use the faster backend - the bytecode compiler and virtual machine. Monkey supports a wide variety of features which will be described below.

## Supported Types

**Booleans**

Booleans are backed by go's native bool type.


```
true
false
```

**Strings**

Strings are backed by go's native string type. Printing is supported via the built-in puts() function. String concatenation is supported with the `+` operator. Strings in Monkey take the form of characters delimited by a pair of double quotes.

Examples:


```
"Graciela"
"Daniela"
"Hugo"

puts("Monkey")
"Monkey " + "Bizness"
```

**Integers**

Integers are backed by go's native int type. Monkey supports basic arithmetic operations on integers.

Examples:


```
1
10000
9122873

1 / 2
1 + 2
1 - 2
1 * 2
```

**Arrays**

Arrays are backed by go's native slice type. Arrays are not scoped to a particular type in Monkey so you can mix and match to your hearts content. Monkey arrays take the form:

`[<expression>, <expression>, ...];`

You can index into an Array with an index expression. Array index expressions take the form:

 `<array>[<expression>];`

Examples:


```
["Ralph", "Abigail", "Bret", "Alejandro"]
[1,2,3,4]
[true, 1, false, "hello"]

["Ralph", "Abigail", "Bret", "Alejandro"][0] -> "Ralph"

let people = ["Ralph", "Abigail", "Bret", "Alejandro"]
people[3] -> "Alejandro"
```


**HashMaps/Dicts/Hashes**

Monkey's kv data type is the Hash and it is backed by a go map. Like Arrays, they are not typed. Hashes take the form:

`{<expression>:<expression, <expression>:<expression, ....};`

You can index into a Hash with an index expression. Hash index expressions takes the form:

 `<hash>[<expression>];`

Examples:

```javascript
{1:2, 3:4, 5:6}"
let animals = {"Rodrigo":"parrot", "William":"giraffe", "Matt":"octopus"}"

animals["Rodrigo"] -> "parrot"
animals["Rod" + "rigo"] -> "parrot"

```

**Functions**

Functions are first class in monkey. Additionally, closures are supported. If you don't have an explicit return in your monkey function, it will implicitly return the last expression.

Functions in Monkey take the form:

```
fn(<optional comma-delimited identifiers>) {
    <optional statements>
    <optional return statement>
}
```

Example closure w/ self referential fibonacci function:

```
let fibonacci = fn(x) {
  if (x == 0) {
    return 0;
  } else {
    if (x == 1) {
      return 1;
    } else {
      return fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};

fibonacci(10);
```

## Statements and Expressions

**Statements**

Programs in Monkey are a series of statements.

Statements don't produce values. There are three types of statements in Monkey.

1. let statements
    - Bind expressions to an identifier
2. return statements
    - return the value produced by an expression from a function
3. expression statements
    - wrap expressions, these values are not reused


**Expressions**

Expressions produce values. These values can be reused in other expressions and combined with the statements listed in the previous section in order to bind an expression to variable or return an expression, etc.

Monkey supports both infix and prefix expressions.

**Let Statements**

Let statements allow you to bind expressions to names in the environment. Let statements scope to where you define them. If you use a let statement in the global scope it will be available to all functions. If you use it within a function, it will be grouped to the lexical scope of the function.

`let <name> = <expression>;`

```
let result = 5 + 5 * 2 / 9
let concat = "fizz" + "buzz"
```


**If expressions**

Monkey supports conditional logic / flow control. This takes the form of:

`if (<expression>) { <statements> } else { <statements> };`

A nice feature of Monkey is it doesn't use if statements but rather if expressions. This allows you to assign to a variable based on conditional logic.

```
let comparison = 5 > 3;

let val = if (comparison) { "Greater than" } else { "less than or equal"};

val -> "Greater than"
```