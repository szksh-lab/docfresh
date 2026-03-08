# Fetch files by GitHub Contents API

When ref is not set, the content is fetched from the default branch.

> [!WARNING]
> GitHub caches the content.
> So when a branch is specified, even if the branch is updated the old content may be fetched.
> This is the problem of GitHub, not docfresh.
> You can avoid the issue by specifying a tag or commit SHA and updating it continuously.

<!-- docfresh begin
github_content:
  owner: suzuki-shunsuke
  repo: docfresh
  path: .gitignore
  ref: main # ref is optional
-->
dist
.coverage
third_party_licenses
.serena
ja
*-ja.md
<!-- docfresh end -->

You can pass a GitHub access token via environment variables `DOCFRESH_GITHUB_TOKEN` or `GITHUB_TOKEN`.

```sh
export DOCFRESH_GITHUB_TOKEN=xxx
```

```sh
export GITHUB_TOKEN=xxx
```

If you use [ghtkn](https://github.com/suzuki-shunsuke/ghtkn), you can pass an access token by ghtkn integration.

```sh
export DOCFRESH_GHTKN_ENABLED=true
```

## Range

Please see [file#Range](20_file.md#range).

## details_tag

[Please see command#details_tag.](10_command.md#details_tag)

<!-- docfresh begin
github_content:
  owner: suzuki-shunsuke
  repo: docfresh
  path: .gitignore
details_tag: {}
use_fenced_code_block_for_output: true
-->
<details>
<summary>suzuki-shunsuke/docfresh/.gitignore</summary>

```
dist
.coverage
third_party_licenses
.serena
ja
*-ja.md
```

</details>
<!-- docfresh end -->
