go_build_wrapper_file = files('go-build-wrapper')
go_build_wrapper_program = find_program('go-build-wrapper')

meson_go_fmt_program = find_program('meson_go_fmt.py')

sources = files(
  'toolbox.go',
  'cmd/completion.go',
  'cmd/create.go',
  'cmd/enter.go',
  'cmd/help.go',
  'cmd/initContainer.go',
  'cmd/list.go',
  'cmd/rm.go',
  'cmd/rmi.go',
  'cmd/root.go',
  'cmd/rootDefault.go',
  'cmd/rootMigrationPath.go',
  'cmd/run.go',
  'cmd/utils.go',
  'pkg/podman/podman.go',
  'pkg/shell/shell.go',
  'pkg/utils/errors.go',
  'pkg/utils/utils.go',
  'pkg/version/version.go',
)

cpu_family = host_machine.cpu_family()
endian = host_machine.endian()

dynamic_linker = ''
if cpu_family == 'aarch64' and endian == 'little'
  dynamic_linker = '/lib/ld-linux-aarch64.so.1'
elif cpu_family == 'arm' and endian == 'little'
  dynamic_linker = '/lib/ld-linux-armhf.so.3'
elif cpu_family == 'ppc64' and endian == 'little'
  dynamic_linker = '/lib64/ld64.so.2'
elif cpu_family == 's390x' and endian == 'big'
  dynamic_linker = '/lib/ld64.so.1'
elif cpu_family == 'x86' and endian == 'little'
  dynamic_linker = '/lib/ld-linux.so.2'
elif cpu_family == 'x86_64' and endian == 'little'
  dynamic_linker = '/lib64/ld-linux-x86-64.so.2'
else
  host_machine_description = cpu_family + ' (' + endian + ' endian)'
  error('Please specify dynamic linker for:', host_machine_description)
endif

message('Host machine dynamic linker:', dynamic_linker)

toolbox_go = custom_target(
  'toolbox',
  command: [
    go_build_wrapper_program,
    meson.current_source_dir(),
    meson.project_build_root(),
    '@OUTPUT@',
    meson.project_version(),
    cc.cmd_array().get(-1),
    dynamic_linker,
    migration_path_for_coreos_toolbox.to_string(),
  ],
  input: sources,
  install: true,
  install_dir: get_option('bindir'),
  output: 'toolbox',
)

if shellcheck.found()
  test('shellcheck src/go-build-wrapper', shellcheck, args: [go_build_wrapper_file])
endif

test('go fmt', meson_go_fmt_program, args: [meson.current_source_dir()])
test('toolbox go unit tests', go, args: ['test', '-v', './...'], workdir: meson.current_source_dir())
