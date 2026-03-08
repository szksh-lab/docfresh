# Embed Local Files

<!-- docfresh begin
file:
  path: file/foo.md
-->
This is read from foo.md
<!-- docfresh end -->

## Wrap the file content with Fenced Code Block

If `code_block` is true, the file content is wrapped with Fenced Code Block.

<!-- docfresh begin
file:
  path: file/create-index.sh
code_block: true
-->
```sh
#!/usr/bin/env bash

set -euo pipefail

dir=$(dirname $(dirname "$0"))

ls "${dir}"/*.md | sed "s|^./||" | grep -v README.md | sort -u | while read -r file; do
    title=$(head -n 1 "$file" | sed -E "s/^# //")
    echo "- [${title}](${file})"
done
```
<!-- docfresh end -->

## Range

Support extracting a specific range of lines from content fetched via `file`, `http`, or `github_content`.
Uses 0-based indexing with half-open interval [start, end).

<!-- docfresh begin
file:
  path: file/range.md
  range:
    start: 2
    end: 4
code_block: true
-->
```md
3
4
```
<!-- docfresh end -->

Negative values count from the end.

<!-- docfresh begin
file:
  path: file/range.md
  range:
    start: -2
code_block: true
-->
```md
9
10
```
<!-- docfresh end -->

## details_tag

[Please see command#details_tag.](10_command.md#details_tag)

<!-- docfresh begin
file:
  path: file/range.md
details_tag: {}
code_block: true
-->
<details>
<summary>file/range.md</summary>

```md
1
2
3
4
5
6
7
8
9
10
```

</details>
<!-- docfresh end -->
