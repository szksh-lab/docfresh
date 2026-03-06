# Fetch Files Via HTTP

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

```md
<!-- docfresh begin
http:
  url: https://jsonplaceholder.typicode.com/todos/1
  timeout: -1 # Disable timeout. The default timeout is 5 seconds.
  header:
    Content-Type:
      - application/json
-->
```
