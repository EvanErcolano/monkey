### TODO

- Extend the Repl to allow multi-line input
    - will need to connect repl into parser to do this
- When an error occurs return stacktrace, filename, line number and so on for better debugabillity
    - will need to look at how we do macros b/c when we replace nodes / create nodes dynamically this will create issues with this error reporting
- Add for / while loops into monkey
- Implement a global scope
- Implement modulo operator