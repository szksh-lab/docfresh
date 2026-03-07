# docfresh

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/suzuki-shunsuke/docfresh)

docfresh is a CLI making document with command and code snippet maintainable, reusable, and testable.
It prevents document from being outdated.

## Status

This is still alpha.
Many features are not implemented yet, and the API is subject to change.
Please don't use this in production.

## Features

- Execute external commands and embed their output into document
- Test external commands in document
- Unify a template file and generated file, which improves the maintainability of document
- Fetch documents from local and remote files and embed them into a document
  - You can share documents across multiple files and repositories
  - You can separate code snippets from the document and apply linters and formatters to them

Note that docfresh is intended to update markdown files.
Other markup language isn't supported.

## Getting Started

1. [Install docfresh](#install)

```sh
: Check version
docfresh -v
```

2. Create a document `date.md`

```md
# Embed the output of date command into a document

<!-- docfresh begin
command:
  command: date "+%Y-%m-%d %H:%M:%S"
-->
<!-- docfresh end -->
```

3. Run `docfresh run date.md` to update date.md.

```sh
docfresh run date.md
```

Then the date command and the result is embedded into the document.

## Examples

Please see not only rendered markdowns but also raw source code because HTML comments are hidden.

> [!TIP]
> This list is created by docfresh.

<!-- docfresh begin
command:
  command: bash examples/file/create-index.sh
template:
  content: |
    {{trimSuffix "\n" .Stdout}}
-->
- [Embed Command Result](examples/10_command.md)
- [Embed Local Files](examples/20_file.md)
- [Fetch Files Via HTTP](examples/30_http.md)
- [Fetch files by GitHub Contents API](examples/40_github_content.md)
- [Customize Template](examples/50-template.md)
- [Test Command Results And Fetched File Contents](examples/60-test.md)
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

```sh
go install github.com/suzuki-shunsuke/docfresh/cmd/docfresh@latest
```

## Security

docfresh may execute arbitrary external commands defined in templates. Therefore, it is important to take appropriate security precautions.
Running docfresh on untrusted templates can be dangerous.
It is recommended to execute docfresh in an isolated environment such as a container.
Secrets should not be provided unless absolutely necessary.

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

## Template Engine

docfresh uses Go's [text/template](https://pkg.go.dev/text/template) and [sprig](https://masterminds.github.io/sprig/).
Note that the following sprig functions aren't available due to security concerns:

- env
- expandenv
- getHostByName

### Available Variables In Templates

command:

- Command
- Stdout
- Stderr
- CombinedOutput
- ExitCode

http, file:

- Content

## Test

`test` is evaluated using [Expr](https://expr-lang.org/).
Clear, precise error messages with position indicators to help debug expressions quickly.

```
+ echo hello
hello
[ERROR] compile an expression
literal not terminated (1:23)
 | Stdout contains "hello
 | ......................^+ echo post
```

[About the language, please see the document of Expr.](https://expr-lang.org/docs/language-definition)

### Available Variables In Expressions

Available Variables are same with available variables in templates.
