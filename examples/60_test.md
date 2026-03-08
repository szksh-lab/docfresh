# Test Command Results And Fetched File Contents

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
