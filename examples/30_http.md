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

### Automatic detection of JSON language

If the URL path doesn't have the file extension, docfresh checks if the response body is JSON.
If the body is JSON, the language is `json` by default.

<!-- docfresh begin
code_block: true
http:
  url: https://jsonplaceholder.typicode.com/todos/1
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

### Specify the language explicitly

If `code_block` is true, the language is automatically detected by the file extension of the URL path.
You can specify the language explicitly by `language`.

<!-- docfresh begin
http:
  url: https://gist.githubusercontent.com/suzuki-shunsuke/7913edb4499fb83bfe86d99c6a2bd42d/raw/52ad913ccbc5ed1af00501175e78b8940768c8d1/docfresh-test-data
  language: yaml
code_block: true
-->
```yaml
name: foo
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
```json
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

```json
{
  "userId": 1,
  "id": 1,
  "title": "delectus aut autem",
  "completed": false
}
```

</details>
<!-- docfresh end -->
