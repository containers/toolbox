subdir('libc-wrappers')

go_build_wrapper_file = files('go-build-wrapper')
go_build_wrapper_program = find_program('go-build-wrapper')

sources = files(
  'toolbox.go',
  'cmd/create.go',
  'cmd/enter.go',
  'cmd/help.go',
  'cmd/initContainer.go',
  'cmd/list.go',
  'cmd/rename.go',
  'cmd/rm.go',
  'cmd/rmi.go',
  'cmd/root.go',
  'cmd/run.go',
  'pkg/podman/podman.go',
  'pkg/shell/shell.go',
  'pkg/utils/utils.go',
  'pkg/version/version.go',
)

custom_target(
  'toolbox',
  build_by_default: true,
  command: [
    go_build_wrapper_program,
    meson.current_source_dir(),
    meson.current_build_dir(),
    meson.project_version(),
    libc_wrappers.full_path(),
  ],
  depends: libc_wrappers,
  input: sources,
  install: true,
  install_dir: get_option('bindir'),
  output: 'toolbox',
)

if shellcheck.found()
  test('shellcheck go-build-wrapper', shellcheck, args: [go_build_wrapper_file])
endif

test('toolbox go unit tests', go, args: ['test', '-v', './...'], workdir: meson.current_source_dir())
