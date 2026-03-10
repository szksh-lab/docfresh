# Validate Files

<!-- docfresh begin
command:
  ignore_fail: true
  command: docfresh validate file/invalid.md
transform:
  CombinedOutput: '{{ regexReplaceAll "[A-Z][a-z]{2}\\s+\\d{1,2}\\s\\d{2}:\\d{2}:\\d{2}\\.\\d{3}" .CombinedOutput "Jan  2 15:04:05.999" }}'
-->
```sh
docfresh validate file/invalid.md
```

Output:

```
file/invalid.md:3
[1:1] unknown field "unknown_field"
>  1 | unknown_field: true
       ^
   2 | command:
   3 |   command: echo hello
Jan  2 15:04:05.999 ERR docfresh failed program=docfresh version=v3.0.0-local error="validate file file/invalid.md: parse file failed"
```
<!-- docfresh end -->
