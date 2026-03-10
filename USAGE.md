# Usage

<!-- docfresh begin
command:
  command: docfresh help-all
  hide_command: true
-->
```console
$ docfresh --help
NAME:
   docfresh - Make document maintainable, reusable, and testable. https://github.com/suzuki-shunsuke/docfresh

USAGE:
   docfresh [global options] [command [command options]]

VERSION:
   v3.0.0-local

COMMANDS:
   run         Update documents
   validate    Validate documents
   version     Show version
   help, h     Shows a list of commands or help for one command
   completion  Output shell completion script for bash, zsh, fish, or Powershell

GLOBAL OPTIONS:
   --log-level string  Log level (debug, info, warn, error) [$DOCFRESH_LOG_LEVEL]
   --help, -h          show help
   --version, -v       print the version
```

## docfresh run

```console
$ docfresh run --help
NAME:
   docfresh run - Update documents

USAGE:
   docfresh run [arguments...]

OPTIONS:
   --allow-unknown-field  Allow unknown fields in directive YAML
   --help, -h             show help
```

## docfresh validate

```console
$ docfresh validate --help
NAME:
   docfresh validate - Validate documents

USAGE:
   docfresh validate [arguments...]

OPTIONS:
   --allow-unknown-field  Downgrade unknown field errors to warnings
   --help, -h             show help
```

## docfresh version

```console
$ docfresh version --help
NAME:
   docfresh version - Show version

USAGE:
   docfresh version

OPTIONS:
   --json, -j  Output version in JSON format
   --help, -h  show help
```

## docfresh completion

```console
$ docfresh completion --help
NAME:
   docfresh completion - Output shell completion script for bash, zsh, fish, or Powershell

USAGE:
   docfresh completion

DESCRIPTION:
   Output shell completion script for bash, zsh, fish, or Powershell.
   Source the output to enable completion.

   # .bashrc
   source <(docfresh completion bash)

   # .zshrc
   source <(docfresh completion zsh)

   # fish
   docfresh completion fish > ~/.config/fish/completions/docfresh.fish

   # Powershell
   Output the script to path/to/autocomplete/docfresh.ps1 an run it.


OPTIONS:
   --help, -h  show help
```
<!-- docfresh end -->
