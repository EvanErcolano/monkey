# The Monkey Programming Language

This repo contains two implementations of the monkey programming language. Both implementations share the same frontends (lexer + parser), however their backends differ. The first backend is a tree-walking interpreter, the second is a single pass bytecode compiler with a companion virtual machine.

The whole shebang is done in go. Ultimately whether the interpreter or compiler+vm is running, the monkey code eventually gets executed in Go.

## To Start the REPL
`go run main.go`

The REPL is configured to use the faster backend - the bytecode compiler and virtual machine. Monkey supports a wide variety of features.

## Supported Types

**Booleans**

Booleans are backed by go's native bool type.

```
true
false
```

**Strings**

Strings are backed by go's native string type. Printing is supported via the built-in puts() function. String concatenation is supported with the `+` operator. Strings in Monkey take the form of characters delimited by a pair of double quotes.

```
"Graciela"
"Daniela"
"Hugo"

puts("Monkey")
"Monkey " + "Bizness"
```

**Integers**

Integers are backed by go's native int type. Monkey supports basic arithmetic operations on integers.

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

`[<expression>, <expression>, ...]`

You can index into an Array with an index expression. Array index expressions take the form:

 `<array>[<expression>]`

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

`{<expression>:<expression, <expression>:<expression, ....}`

You can index into a Hash with an index expression. Hash index expressions takes the form:

 `<hash>[<expression>]`

```
{1:2, 3:4, 5:6}"
let animals = {"Rodrigo":"parrot", "William":"giraffe", "Matt":"octopus"}"

animals["Rodrigo"] -> "parrot"
animals["Rod" + "rigo"] -> "parrot"

```

**Functions**

lorem ipsum...

## Statements

Programs in Monkey are a series of statements.

Statements don't produce values. There are three types of statements in Monkey.

1. let statements
    - Bind expressions to an identifier
2. return statements
    - return the value produced by an expression from a function
3. expression statements
    - wrap expressions, these values are not reused


## Expressions

Expressions produce values. These values can be reused in other expressions and combined with the statements listed in the previous section in order to bind an expression to variable or return an expression, etc.

There are two type of expressions in Monkey:

1. Infix Expressions
2. Prefix Expressions

explain...

### If expressions
