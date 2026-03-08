# Fetch Files Via HTTP

<!-- docfresh begin
http:
  url: https://raw.githubusercontent.com/suzuki-shunsuke/docfresh/refs/heads/main/_typos.toml
-->
[default.extend-words]
ERRO = "ERRO"
intoto = "intoto"
typ = "typ"
<!-- docfresh end -->

## Wrap the file content with Fenced Code Block

If `code_block` is true, the file content is wrapped with Fenced Code Block.

<!-- docfresh begin
http:
  url: https://raw.githubusercontent.com/suzuki-shunsuke/docfresh/refs/heads/main/_typos.toml
code_block: true
-->
```toml
[default.extend-words]
ERRO = "ERRO"
intoto = "intoto"
typ = "typ"
```
<!-- docfresh end -->

## Specify the language explictly

If `code_block` is true, the language is automatically deteced by the file extension of the URL path.
You can specify the language explictly by `language`.

<!-- docfresh begin
http:
  url: https://jsonplaceholder.typicode.com/todos/1
  language: json
code_block: true
-->
```json
{
  "userId": 1,
  "id": 1,
  "title": "delectus aut autem",
  "completed": false
}
```
<!-- docfresh end -->

## timeout, header

You can set the timeout and header.

<!-- docfresh begin
code_block: true
http:
  url: https://jsonplaceholder.typicode.com/todos/1
  timeout: -1 # Disable timeout. The default timeout is 5 seconds.
  header:
    Content-Type:
      - application/json
-->
```
{
  "userId": 1,
  "id": 1,
  "title": "delectus aut autem",
  "completed": false
}
```
<!-- docfresh end -->

## Range

Please see [file#Range](20_file.md#range).

## details_tag

[Please see command#details_tag.](10_command.md#details_tag)

<!-- docfresh begin
http:
  url: https://jsonplaceholder.typicode.com/todos/1
details_tag: {}
code_block: true
-->
<details>
<summary>https://jsonplaceholder.typicode.com/todos/1</summary>

```
{
  "userId": 1,
  "id": 1,
  "title": "delectus aut autem",
  "completed": false
}
```

</details>
<!-- docfresh end -->
