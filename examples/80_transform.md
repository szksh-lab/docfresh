# Transform command outputs and file contents before rendering

[#93](https://github.com/suzuki-shunsuke/docfresh/issues/93)

Preprocessing `Content`, `Stdout`, `Stderr`, and `CombinedOutput` before rendering them with a template.

Usecases:

- Mask sensitive information
- Replace dynamic values such as timestamps (which change on every execution) with fixed dummy values to reduce noise

<!-- docfresh begin
command:
  command: |
    echo "datetime: $(date "+%Y-%m-%d %H:%M:%S")"
transform:
  CombinedOutput: '{{regexReplaceAll "\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}" .CombinedOutput "2006-01-02 15:04:05" }}'
-->
```sh
echo "datetime: $(date "+%Y-%m-%d %H:%M:%S")"
```

Output:

```
datetime: 2006-01-02 15:04:05
```
<!-- docfresh end -->
