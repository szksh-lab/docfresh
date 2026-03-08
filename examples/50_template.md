# Customize Template

<!-- docfresh begin
command:
  command: echo Hello
template:
  content: |
    ```console
    $ {{.Command}}
    ```
    
    ```
    {{trimSuffix "\n" .CombinedOutput}}
    ```
-->
```console
$ echo Hello
```

```
Hello
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

## Change Delimiters

The default delimiters is `{{` and `}}`.

<!-- docfresh begin
command:
  command: echo Hello
template:
  delims:
    left: "<["
    right: "]>"
  content: |
    ```sh
    <[.Command]>
    ```
-->
```sh
echo Hello
```
<!-- docfresh end -->

`*.template.delims` work as well.

<!-- docfresh begin
file:
  path: file/template4.md
  template:
    vars:
      Name: foo
    delims:
      left: "<["
      right: "]>"
-->
Name: foo
<!-- docfresh end -->

## Render file as templates

<!-- docfresh begin
file:
  path: file/template.md
  template:
    vars:
      name: foo
-->
Hello, foo
<!-- docfresh end -->

## Template Variables

<!-- docfresh begin
command:
  command: echo Hello
template:
  path: file/template3.md
  vars:
    project: foo
-->
project: foo

```console
$ echo Hello
```

```
Hello
```
<!-- docfresh end -->
