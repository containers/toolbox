% toolbox-help 1

## NAME
toolbox\-help - Display help information about Toolbx

## SYNOPSIS
**toolbox help** [*COMMAND*]

## DESCRIPTION

When no COMMAND is specified, the `toolbox(1)` manual is shown. If a COMMAND
is specified, a manual page for that command is brought up.

Note that `toolbox --help ...` is identical to `toolbox help ...` because the
former is internally converted to the latter.

This page can be displayed with `toolbox help help` or `toolbox help --help`.

## EXAMPLES

### Show the toolbox manual

```
$ toolbox help
```

### Show the manual for the create command

```
$ toolbox help create
```

## SEE ALSO

`toolbox(1)`
