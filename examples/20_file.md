# Embed Local Files

<!-- docfresh begin
file:
  path: file/foo.md
-->
This is read from foo.md
<!-- docfresh end -->

## Wrap the file content with Fenced Code Block

If `use_fenced_code_block_for_output` is true, the file content is wrapped with Fenced Code Block.

<!-- docfresh begin
file:
  path: file/create-index.sh
use_fenced_code_block_for_output: true
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
use_fenced_code_block_for_output: true
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
use_fenced_code_block_for_output: true
-->
```md
9
10
```
<!-- docfresh end -->
