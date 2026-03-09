# validate command

<!-- docfresh begin
unknown_field: true
command:
  ignore_fail: true
  command: docfresh validate 90_validate.md
-->
```sh
docfresh validate 90_validate.md
```

Output:

```
[1:1] unknown field "unknown_field"
>  1 | unknown_field: true
       ^
   2 | command:
   3 |   ignore_fail: true
   4 |   command: docfresh validate 90_validate.md
Mar  9 23:24:51.309 ERR docfresh failed program=docfresh version=v3.0.0-local error="validate file 90_validate.md: parse file failed"
```
<!-- docfresh end -->
