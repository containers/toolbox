go_build_wrapper_file = files('go-build-wrapper')
go_build_wrapper_program = find_program('go-build-wrapper')

sources = files(
  'toolbox.go',
  'cmd/root.go',
  'pkg/podman/podman.go',
  'pkg/shell/shell.go',
  'pkg/utils/utils.go',
  'pkg/version/version.go',
)

if go.found()
  custom_target(
    'toolbox-go',
    build_by_default: true,
    command: [go_build_wrapper_program, meson.current_source_dir(), meson.current_build_dir()],
    input: sources,
    output: 'toolbox',
  )
endif

if shellcheck.found()
  test('shellcheck go-build-wrapper', shellcheck, args: [go_build_wrapper_file])
endif