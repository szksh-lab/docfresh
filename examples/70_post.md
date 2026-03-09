# Post (Cleanup)

If you want to ensure that some processing is always executed at the end of a file’s processing, you can use `<!-- docfresh post ... -->`.  
Commands defined in this comment are executed at the end of the file processing regardless of whether the file processing succeeded or failed.

e.g.

```md
<!-- docfresh post
command: |
  test ! -f aqua.yaml || rm aqua.yaml
-->
````

The `post_command` in `<!-- docfresh begin ... -->` is executed at the end of processing for that specific comment.
However, if you want to carry results across multiple comments, cleanup may not be reliably handled with `begin`’s `post_command`.

For example, suppose one comment **A** generates a file `A`, the next comment **B** references that file, and finally you want to delete file `A`.
In this case:

- You cannot delete file `A` in comment **A**’s `post_command`, because comment **B** still needs it.
- If comment **A**’s command fails, file `A` may remain without being cleaned up.

If multiple `<!-- docfresh post ... -->` comments are defined in a file, all of them are executed in reverse order of their definition (i.e., the most recently defined runs first).

If a post comment fails, processing does not stop immediately; all post comments will still be executed.
However, if any post comment fails, `docfresh run` will ultimately fail.

<!-- docfresh post
command: |
  echo "post 1" >&2
  test ! -f foo.txt || rm foo.txt
-->

<!-- docfresh post
command: |
  echo "post 2" >&2
  test ! -f foo.txt || rm foo.txt
-->

<!-- docfresh begin
command:
  command: echo foo > foo.txt
  hide_output: true
-->
```sh
echo foo > foo.txt
```

<!-- docfresh end -->

<!-- docfresh begin
command:
  command: cat foo.txt
-->
```sh
cat foo.txt
```

Output:

```
foo
```
<!-- docfresh end -->
