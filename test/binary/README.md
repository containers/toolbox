# Binary tests

These tests test the Toolbx binary if it is properly set up (e.g., rpath, path
to dynamic linker)

## Dependencies

- `go`

## How to run the tests

The tests are best run using the build system that takes care of passing proper
values to them.

In case you want/need to run them manually run this command in this directory:

```shell
$ go test ./... -args <path-to-toolbx-binary> <path-to-file-with-build-overrides>
```

The tests read the tested values from a file where the build script wrote the
values it embeded into the binary. The values are space-separated on a single
line.
