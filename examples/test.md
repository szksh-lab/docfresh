# Test

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

```
hello
```
<!-- docfresh end -->
