# docfresh

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/suzuki-shunsuke/docfresh) | [Install](INSTALL.md)

docfresh is a CLI making document with command and code snippet maintainable, reusable, and testable.
It prevents document from being outdated.

## Features

- Execute external commands and embed their output into document
- Test external commands in document
- Unify a template file and generated file, which improves the maintainability of document
- Fetch documents from local and remote files and embed them into a document
  - You can share documents across multiple files and repositories
  - You can separate code snippets from the document and apply linters and formatters to them
- [Allow running commands in containers](examples/15_container.md)
  - Securely run commands in isolated environments
  - Consistent Command Results across different machines

Note that docfresh is intended to update markdown files.
Other markup language isn't supported.

## Getting Started

1. [Install docfresh](INSTALL.md)
2. [Run commands and embed the output into a document](#2-run-commands-and-embed-the-output-into-a-document)
3. [Embed Local Files](#3-embed-local-files)
4. [Embed Remote Files](#4-embed-remote-files)

### 1. Install docfresh

[Please see INSTALL.md.](INSTALL.md)

<!-- docfresh begin
command:
  command: docfresh -v
-->
```sh
docfresh -v
```

Output:

```
docfresh version v3.0.0-local
```
<!-- docfresh end -->

### 2. Run commands and embed the output into a document

2. Create a document `hello.md`

<!-- docfresh begin
file:
  path: examples/file/hello.md
code_block: true
-->
```md
# Embed the output of commands into a document

<!-- docfresh begin
command:
  command: echo hello
-->
<!-- docfresh end -->
```
<!-- docfresh end -->

The HTML comments `<!-- docfresh begin -->` and `<!-- docfresh end -->` are docfresh's directives to define the action and the output to be embedded into the document.
The action result is embedded between the begin and end directives.
They are HTML comments, so they are not visible in the rendered markdown.

3. Validate document by `docfresh run hello.md`.

<!-- docfresh begin
command:
  command: docfresh validate hello.md
  dir: examples/file
-->
```sh
docfresh validate hello.md
```

Output:

```
hello.md: valid
```
<!-- docfresh end -->

4. Update document by `docfresh run hello.md`.

```sh
docfresh run hello.md
```

Then `echo hello` is executed and the result is embedded into hello.md.
You can confirm commands in the document work properly and the document is updated.
By running `docfresh run` periodically, you can detect issues quickly and keep the document up-to-date.

````md
<!-- docfresh begin
command:
  command: echo hello
-->
```sh
echo hello
```

Output:

```
hello
```
<!-- docfresh end -->
````

<!-- docfresh begin
command:
  command: echo hello
-->
```sh
echo hello
```

Output:

```
hello
```
<!-- docfresh end -->

[You can also run commands in containers, which improves the security and consistency (portability).](examples/15_container.md)

### 3. Embed Local Files

You can embed local files into the document.
You can manage code snippets in documents separately, which improves the maintainability of code snippets.
You can run linter and formatter on code snippets, and edit code snippets by your editor.

For example, embedding [_typos.toml](_typos.toml).

```md
<!-- docfresh begin
file:
  path: _typos.toml
code_block: true
-->
<!-- docfresh end -->
```

<!-- docfresh begin
file:
  path: _typos.toml
code_block: true
-->
```toml
[default.extend-words]
ERRO = "ERRO"
intoto = "intoto"
typ = "typ"
```
<!-- docfresh end -->

### 4. Embed Remote Files

You can embed remote files into the document.
You can reuse documents across repositories.
For example, you can reuse CONTRIBUTING.md, AGENTS.md, CLAUDE.md, and so on.
You can reuse Go's template files, which can change contents dynamically.

```md
<!-- docfresh begin
code_block: true
http:
  url: https://jsonplaceholder.typicode.com/todos/1
-->
<!-- docfresh end -->
```

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

## Examples

Please see not only rendered markdowns but also raw source code because HTML comments are hidden.

> [!TIP]
> This list is created by docfresh.

<!-- docfresh begin
command:
  command: bash examples/file/create-index.sh
  hide_command: true
-->
- [Embed Command Result](examples/10_command.md)
- [Running Commands In Containers](examples/15_container.md)
- [Embed Local Files](examples/20_file.md)
- [Fetch Files Via HTTP](examples/30_http.md)
- [Fetch files by GitHub Contents API](examples/40_github_content.md)
- [Customize Template](examples/50_template.md)
- [Test Command Results And Fetched File Contents](examples/60_test.md)
- [Post (Cleanup)](examples/70_post.md)
- [Transform command outputs and file contents before rendering](examples/80_transform.md)
- [Validate Files](examples/90_validate.md)
<!-- docfresh end -->

## Motivation

Keeping documentation accurate is not easy.
Commands and code in documentation can quickly become outdated, which discourages readers.
When execution results are included in the documentation, it is tedious to manually rerun commands and update the results each time something changes.

With docfresh, commands in documentation can be executed automatically, and their results can be embedded directly into the documentation.
It also helps you quickly detect when commands start failing.

By running docfresh in CI, you can automate documentation updates and validation.

Another key feature of docfresh is that templates and generated files are unified, making documentation easier to maintain.

When templates and generated files are separate, it creates the question of where and how to manage the template files:

- Using a special extension like `.tpl`: syntax highlighting in editors may no longer work properly.
- Adding a suffix like `-tpl` to filenames: it becomes harder to navigate as the number of files increases, and static site generators may include template files unintentionally.
- Separating them into different directories: it becomes less obvious that templates exist and where they are located.

There are also practical issues.
Even if you delete a template file, the generated file may remain.
Since the editable file and the generated file are separate, it is also harder to edit while previewing the final output.

With docfresh, templates and generated files are unified, so these problems do not occur.

## Usecases

### Embed the common contribution guide

We maintain a lot of OSS projects and we maintain the common contribution guide at [suzuki-shunsuke/oss-contribution-guide](https://github.com/suzuki-shunsuke/oss-contribution-guide/blob/v0.1.0/docs/guide.md).
Previously, we noted the link to the document in CONTRUBTING.md of each project, but probably many people didn't follow it.
So now we embed the common guide to CONTRIBUTING.md of each project using docfresh.

- People and AI don't need to follow the link
  - They follow the guide more
- People and AI can search the content using tools like grep

https://github.com/suzuki-shunsuke/docfresh/blob/8df6bbf991279c8d8fafa99b5acc4b2ace3dc192/CONTRIBUTING.md?plain=1#L3-L6

## Install

[Please see INSTALL.md.](INSTALL.md)

## Security

docfresh may execute arbitrary external commands defined in templates.
Therefore, it is important to take appropriate security precautions.
Running docfresh on untrusted templates can be dangerous.
It is recommended to execute docfresh in an isolated environment such as a container.
Secrets should not be provided unless absolutely necessary.

## Running in CI

Running docfresh in CI helps prevent differences in generated results caused by variations in the execution environment.
It also enables quick detection of problems and prevents documentation from becoming inconsistent or outdated.
When running in CI, there are two possible approaches:

1. Fail CI if the documentation is updated.
2. Automatically update the documentation in CI.

Approach 2 generally reduces the maintenance burden.
For automatic updates:

- For public repositories, [autofix.ci](https://autofix.ci/) is convenient.
- For private repositories, [Securefix Action](https://github.com/csm-actions/securefix-action) is useful.

When running docfresh in CI, you may also want to restrict its execution to prevent arbitrary code execution from malicious PRs.

e.g.

- Restrict by PR author
- Restrict by `github.actor`
- Disable execution for PRs from forks

However, for `pull_request` events from forks, secrets cannot be accessed and `secrets.GITHUB_TOKEN` only has read permissions. Therefore, the risk is relatively limited, so excessive restrictions may not be necessary.
On the other hand, running docfresh with `pull_request_target` or `workflow_run` events in public repositories can be dangerous and should generally be avoided.

## Template Syntax

In docfresh, instructions are embedded into Markdown using HTML comments, as shown below:

```md
<!-- docfresh begin
command:
  command: npm test
-->

The result will be embedded here.

<!-- docfresh end -->
```

Instructions are written in YAML format inside `<!-- docfresh begin -->`, and the execution result is embedded between `<!-- docfresh begin -->` and `<!-- docfresh end -->`.
Since HTML comments are not rendered in Markdown, they don't affect the view of the documentation.
Because this mechanism relies on HTML comments, docfresh is designed specifically for Markdown and does not support other document formats.
Each directive must start with `<!-- docfresh begin` and must have a corresponding closing `<!-- docfresh end -->`.

## Command and File Processing Order

docfresh executes all file processing and commands sequentially.
Commands within the same file are executed from top to bottom. If a command fails, the file will not be updated.

## YAML Syntax In Begin Comment

```md
<!-- docfresh begin
command:
  command: npm test
-->
```

Please see [JSON Schema](json-schema/comment.json) and [Examples](#examples).
