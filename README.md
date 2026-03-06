# docfresh

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
- Fetch document from local and remote files and embed them into document

Note that docfresh is intended to update markdown files.
Other markup language isn't supported.

## Getting Started

1. [Install docfresh](#install)

```sh
: Check version
docfresh -v
```

2. Checkout the repository

```sh
git clone https://github.com/suzuki-shunsuke/docfresh
cd examples
```

Please see [date.md](examples/date.md).
In this document, the result of `date` command is embedded.

```sh
cat date.md
```

Please run `docfresh run date.md` to update date.md.

```sh
docfresh run date.md
```

Then the datetime is updated.

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
Running docfresh on untrusted templates can be dangerous. It is recommended to execute docfresh in an isolated environment such as a container. Secrets should not be provided unless absolutely necessary.
Support for executing commands inside containers is also being considered for future releases.

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

Support for parallel processing across multiple files may be added in the future.

## YAML Syntax In Begin Comment

```md
<!-- docfresh begin
command:
  command: npm test
  shell:
    - bash
    - "-c"
-->
```

- command.command: External Command
- command.dir: The relative path from the current file to the directory where the command is executed. By default, the directory of the file.
- command.shell: The list of shell command executing command. By default, `["bash", "-c"]`
- command.ignore_fail: Ignore command failure. By default, `false`
- command.envs: Environment variables.
- pre_command: External Command executed before `command`. If it fails, docfresh fails and `command` isn't run. The command and output are outputted to the console but the result isn't affected to the document. This is used for setup and checking the requirement. The format is same with `command`
- post_command: External Command executed after `command`. If it fails, docfresh fails. The command and output are outputted to the console but the result isn't affected to the document. This is used for testing the command result and cleaning up. The format is same with `command`. `post_command` is run even if `pre_command` and `command` fail.
- file.path: The relative path from the current file to the loaded file
- http.url: The URL to fetch the content from
- template.content: The content to be rendered by the template engine
- github_content.owner: GitHub repository owner
- github_content.repo: GitHub repository name
- github_content.path: GitHub repository content path
- github_content.ref: GitHub repository content ref. This is optional.

### Examples

[Please see examples.](examples)

### Template Engine

docfresh uses Go's [text/template](https://pkg.go.dev/text/template) and [sprig](https://masterminds.github.io/sprig/).
Note that the following sprig functions aren't available due to security concerns:

- env
- expandenv
- getHostByName

#### Available Variables In Templates

command:

- Command
- Stdout
- Stderr
- CombinedOutput
- ExitCode

http, file:

- Content

### Run Command

```md
<!-- docfresh begin
command:
  command: npm test
-->
```

### Change Shell

```md
<!-- docfresh begin
command:
  command: echo hello
  shell:
    - zsh
    - "-c"
-->
```

### Ignore Command Failure

By default, `docfresh run` fails if any command fails.
If `.command.ignore_fail` is set to `true`, the command failure will be ignored.

```md
<!-- docfresh begin
command:
  command: npm t
  ignore_fail: true
-->
```

### Pre-Command, Post-Command 

```md
<!-- docfresh begin
pre_command:
  command: npm ci
command:
  command: npm t
  ignore_fail: true
post_command:
  command: rm -rf node_modules
-->
```

### Read File

```md
<!-- docfresh begin
file:
  path: foo.md
-->
```

### Fetch File via HTTP

```md
<!-- docfresh begin
http:
  url: https://raw.githubusercontent.com/suzuki-shunsuke/docfresh/refs/heads/main/_typos.toml
-->
```

### Change Template

```md
<!-- docfresh begin
command:
  command: echo hello
template:
  content: |
    Command:
    {{.Command}}
    
    Stdout:
  
    {{.Stdout}}
    
    Stderr:
    
    {{.Stderr}}
-->
```

### Fetch File by GitHub Content API

When ref is not set, the content is fetched from the default branch.

> [!WARNING]
> GitHub caches the content.
> So when a branch is specified, even if the branch is updated the old content may be fetched.
> This is the problem of GitHub, not docfresh.
> You can avoid the issue by specifying a tag or commit SHA and updating it continuously.

```md
<!-- docfresh begin
github_content:
  owner: suzuki-shunsuke
  repo: docfresh
  path: README.md
  ref: main # ref is optional
-->
```

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

### Read Files As Templates

```
Hello, {{.Vars.name}}
```

```md
<!-- docfresh begin
file:
  path: file/template.md
  template:
    vars:
      name: foo
-->
```

### Test

```md
<!-- docfresh begin
command:
  command: echo test
  test: |
    Stdout contains "test"
-->
```

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

#### Available Variables In Expressions

Available Variables are same with available variables in templates.
