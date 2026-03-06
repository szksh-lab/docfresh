#!/usr/bin/env bash

set -euo pipefail

dir=$(dirname $(dirname "$0"))

ls "${dir}"/*.md | sed "s|^./||" | grep -v README.md | sort -u | while read -r file; do
    title=$(head -n 1 "$file" | sed -E "s/^# //")
    echo "- [${title}](${file})"
done
