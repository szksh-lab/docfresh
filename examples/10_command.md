# Embed Command Result

## Hello

<!-- docfresh begin
command:
  command: echo hello
-->
```sh
echo hello
```

```
hello
```
<!-- docfresh end -->

## command.dir

<!-- docfresh begin
command:
  command: cat foo.md
  dir: file
-->
```sh
cat foo.md
```

```
This is read from foo.md
```
<!-- docfresh end -->

## Set Environment Variables

<!-- docfresh begin
command:
  command: echo "$FOO"
  envs:
    FOO: foo
-->
```sh
echo "$FOO"
```

```
foo
```
<!-- docfresh end -->

## Change Shell

<!-- docfresh begin
command:
  command: console.log("hello")
  shell:
    - node
    - "-e"
template:
  content: |
    ```js
    {{.Command}}
    ```
    
    ```
    {{trimSuffix "\n" .CombinedOutput}}
    ```
-->
```js
console.log("hello")
```

```
hello
```
<!-- docfresh end -->

## command.script

Instead of `command.command`, you can specify an external script by `command.script`.

`shell` is automatically detected in case of some popular languages such as Go and Python.
If it can't be automatically detected, `shell` should be specified explicitly.

<!-- docfresh begin
command:
  script: file/hello.sh
-->
```sh
bash file/hello.sh
```

```
Hello
```
<!-- docfresh end -->

### Embed the content of command.script

If `command.embed_script` is true, the script content is embedded.
`script_language` and `shell` are automatically detected in case of some popular languages such as Go and Python.

<!-- docfresh begin
command:
  script: file/hello.sh
  embed_script: true
template:
  content: |
    ```sh
    {{trimSuffix "\n" .Content}}
    ```

    ```
    {{trimSuffix "\n" .CombinedOutput}}
    ```
-->
```sh
#!/usr/bin/env bash

echo Hello
```

```
Hello
```
<!-- docfresh end -->

### Automatic detection of script languages by file extensions

`script_language` and `shell` are automatically detected in case of some popular languages such as Go and Python.

[languages.yaml](../pkg/controller/run/languages.yaml)

<!-- docfresh begin
file:
  path: ../pkg/controller/run/languages.yaml
use_fenced_code_block_for_output: true
-->
```yaml
go:
  shell:
    - go
    - run
  extensions:
    - .go
hcl:
  extensions:
    - .hcl
js:
  shell:
    - node
  extensions:
    - .js
json:
  extensions:
    - .json
md:
  extensions:
    - .md
py:
  shell:
    - python3
  extensions:
    - .py
sh:
  shell:
    - bash
  extensions:
    - .sh
    - .bash
tf:
  extensions:
    - .tf
toml:
  extensions:
    - .toml
ts:
  extensions:
    - .ts
yaml:
  extensions:
    - .yaml
    - .yml
```
<!-- docfresh end -->

## ignore_fail

By default, `docfresh run` fails if any command fails.
If `.command.ignore_fail` is set to `true`, the command failure will be ignored.

<!-- docfresh begin
command:
  command: |
    echo "failed to install" >&2
    exit 1
  ignore_fail: true
-->
```sh
echo "failed to install" >&2
exit 1
```

```
failed to install
```
<!-- docfresh end -->

## pre_command, post_command

<!-- docfresh begin
pre_command:
  command: echo temporary > temporary.txt
command:
  command: cat temporary.txt
post_command:
  command: rm temporary.txt
-->
```sh
cat temporary.txt
```

```
temporary
```
<!-- docfresh end -->
