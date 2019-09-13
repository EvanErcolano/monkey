# The Monkey Programming Language

This repo contains two implementations of the Monkey Programming language. I built this while following along with Thorsten Ball's incredible books on the topics of building [interpreters](https://interpreterbook.com/) and [compilers](https://compilerbook.com/).

Both implementations share the same frontends (lexer + parser), however their backends differ. The first backend is a tree-walking interpreter, the second is a single pass bytecode compiler with a companion virtual machine.

Ultimately whether the interpreter or Bytecode-Compiler+VM is running, the monkey code eventually is executed in native Go.

## To Start the REPL
`go run main.go`

The REPL is configured to use the faster backend - the Bytecode-Compiler and Virtual Machine. Monkey supports a wide variety of features...

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

```
{1:2, 3:4, 5:6}"
let animals = {"Rodrigo":"parrot", "William":"giraffe", "Matt":"octopus"}

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

Example self-referential recursive function:

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

Example Closure

```
let newClosure = fn(a,b) {
    let one = fn() {a;};
    let two = fn() {b;};
    return fn() {one() + two();};
};

let closure = newClosure(9,90);
closure(); -> 99
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

Expressions produce values. These values can be reused in other expressions and combined with the statements listed in the previous section in order to bind an expression to a variable or return an expression...

Monkey supports both infix and prefix expressions.

**Let Statements**

Let statements allow you to bind expressions to names in the environment. Let statements scope to where you define them. If you use a let statement in the global scope it will be available to all functions. If you use it within a function, it will be bound to the lexical scope of that function.

`let <name> = <expression>;`

```
let result = 5 + 5 * 2 / 9
let concat = "fizz" + "buzz"
```


**If Expressions**

Monkey supports conditional logic / flow control. This takes the form of:

`if (<expression>) { <statements> } else { <statements> };`

A nice feature of Monkey is it doesn't use if statements but rather if expressions. This allows you to assign to a variable based on conditional logic.

```
let comparison = 5 > 3;

let val = if (comparison) { "Greater than" } else { "less than or equal"};

val -> "Greater than"
```

## Nice to haves and things to improve

During this process I realized I take the python REPL for granted, it has so many neat features that are lacking here. For example the REPL:

- Does not support using your arrow keys to go backwards or forwards in the text you already wrote.
- Does not support multi-line REPL input. So you have to remove newlines from a long monkey function definition if you intend on inputting it via the REPL.

This could improved by integrating the REPL, lexer and parser together so that the REPL could know when a function is finished being defined, else it will continue to allow input...

Additionally, it was very noticable how important parser errors are. The parser is in a way the main user interface for your programming language and if you don't have good parser errors your users will not have any idea what they did wrong. For that reason I'd like to:

- Improve error handling in the parser
- Present nicer errors to the user and point to where exactly the error happened - reminenscent of error reporting I've seen in rust.
- Introduce stack traces in Monkey so when there is an error it's easy to trace execution and find the issue.

Extra language features that would be cool to add to Monkey:

- Pattern Matching
- for / while loops (we have recursion only)
- Ability to define libraries / modules and import them

The extra functionality I'd like to add to the Bytecode-Compiler mostly revolves around learning about optimization:

- Dead Code Elimination
- Constant Folding
- Strength Reduction
