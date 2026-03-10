# Test Command Results And Fetched File Contents

`test` is evaluated using [Expr](https://expr-lang.org/).
Clear, precise error messages with position indicators to help debug expressions quickly.

```
+ echo hello
hello
[ERROR] compile an expression
literal not terminated (1:23)
 | Stdout contains "hello
 | ......................^+ echo post
```

[About the language, please see the document of Expr.](https://expr-lang.org/docs/language-definition)

## Test Command

<!-- docfresh begin
pre_command:
  command: echo pre
  test: |
    Stdout contains "pre"
command:
  command: echo hello
  test: |
    Stdout contains "hello"
post_command:
  command: echo post
  test: |
    Stdout contains "post"
-->
```sh
echo hello
```

Output:

```
hello
```
<!-- docfresh end -->

## Test file

<!-- docfresh begin
file:
  path: file/foo.md
  test: |
    Content contains "foo.md"
-->
This is read from foo.md
<!-- docfresh end -->

## Available Variables In Expressions

Available Variables are same with available variables in templates.
