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

If `use_fenced_code_block_for_output` is true, the file content is wrapped with Fenced Code Block.

<!-- docfresh begin
http:
  url: https://raw.githubusercontent.com/suzuki-shunsuke/docfresh/refs/heads/main/_typos.toml
use_fenced_code_block_for_output: true
-->
```toml
[default.extend-words]
ERRO = "ERRO"
intoto = "intoto"
typ = "typ"
```
<!-- docfresh end -->

## timeout, header

You can set the timeout and header.

<!-- docfresh begin
use_fenced_code_block_for_output: true
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
