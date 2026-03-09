# Running Commands In Containers

Creating containers and running commands inside them.

<!-- docfresh container
id: foo
engine: docker-cli
image: ubuntu:24.04
-->

<!-- docfresh begin
command:
  command: whoami
  container:
    id: foo
-->
```sh
whoami
```

Output:

```
root
```
<!-- docfresh end -->

`container` directive itself doesn't update documents.
This directive performs the following steps:

1. Create a container
1. (Optional) Copy specified host files into the container
1. (Optional) Execute the specified command (e.g., to install required tools or set up the environment)

A new container is created each time.

`id` in `container` directive is not the actual container ID but an identifier used by docfresh to track containers.
`id` must be unique within the same document, but there is no risk of conflicting with existing containers.

After the container is created, subsequent `begin` and `post` directives can specify the container ID to run commands inside the container.

Unlike `begin` directive, `container` directive doesn't support `file`, `http`, `github_content`, `pre_command`, or `post_command`.
Only `command` is supported.

The `command` configuration is mostly the same as in `begin` directive, with the following differences:

- `env` is applied to the container.
- `dir` refers to a path inside the container.

## Container Removal

By default, containers are removed after file processing finishes.
If a command fails, the container used to run that command is not removed and remains for debugging.

If `keep: true` is specified, the container will not be removed.
This is mainly intended for temporary troubleshooting.
Be aware that if `keep: true` remains enabled and docfresh is run repeatedly, containers will accumulate.

### Labels

Containers are automatically assigned labels so that they can be identified as containers created by docfresh.

- `docfresh.file_path`: the file path passed as an argument to `docfresh run`. This may be either a relative or an absolute path.
- `docfresh.absolute_file_path`: the absolute file path.
- `docfresh.id`: the container ID used internally by docfresh.

You can search for containers created by docfresh using these labels.

<!-- docfresh begin
command:
  command: |
    docker ps --filter "label=docfresh.id"
  quiet: true
-->
```sh
docker ps --filter "label=docfresh.id"
```

<!-- docfresh end -->

Remove all containers created by docfresh:

```sh
docker ps --filter "label=docfresh.id" -q | xargs docker rm -f
```

## Engine

How containers are managed depends on the `container` directive's `engine` field.
Currently, the only supported engine is `docker-cli`, which executes commands such as `docker run`, `docker exec`, and `docker rm`.

### docker-cli engine

This engine executes commands such as `docker run`, `docker exec`, and `docker rm`.

`Stdout`, `Stderr`, `CombinedOutput`, and `ExitCode` correspond to the output of `docker exec`.
`Command` refers to the command executed inside `docker exec` (excluding the `docker exec` wrapper itself).
