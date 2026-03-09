# validate command

<!-- docfresh begin
unknown_field: true
command:
  ignore_fail: true
  command: docfresh validate 90_validate.md
transform:
  CombinedOutput: '{{ regexReplaceAll "[A-Z][a-z]{2}\\s+\\d{1,2}\\s\\d{2}:\\d{2}:\\d{2}\\.\\d{3}" .CombinedOutput "Jan  2 15:04:05.999" }}'
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
   5 | 
Jan  2 15:04:05.999 ERR docfresh failed program=docfresh version=v3.0.0-local error="validate file 90_validate.md: parse file failed"
```
<!-- docfresh end -->
