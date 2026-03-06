# command.shell & template

## command.shell & template.content

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

## template.path

<!-- docfresh begin
command:
  command: echo "read template file"
template:
  path: file/template2.md
-->
```console
$ echo "read template file"
read template file
```
<!-- docfresh end -->
