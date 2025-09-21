meson_go_fmt_program = find_program('meson_go_fmt.py')

sources = files(
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
  'cmd/root_test.go',
  'cmd/run.go',
  'cmd/utils.go',
  'pkg/nvidia/nvidia.go',
  'pkg/podman/container.go',
  'pkg/podman/containerInspect_test.go',
  'pkg/podman/errors.go',
  'pkg/podman/podman.go',
  'pkg/shell/shell.go',
  'pkg/shell/shell_test.go',
  'pkg/skopeo/skopeo.go',
  'pkg/term/term.go',
  'pkg/term/term_test.go',
  'pkg/utils/arch.go',
  'pkg/utils/errors.go',
  'pkg/utils/fedora.go',
  'pkg/utils/libsubid-wrappers.c',
  'pkg/utils/rhel.go',
  'pkg/utils/utils.go',
  'pkg/utils/utils_cgo.go',
  'pkg/utils/utils_test.go',
  'pkg/version/version.go',
  'toolbox.go',
)

go_tags = []
if migration_path_for_coreos_toolbox
  go_tags += 'migration_path_for_coreos_toolbox'
endif

toolbox_go = custom_target(
  'toolbox',
  command: [
    go,
    'build',
    '-C', meson.current_source_dir(),
    '-trimpath',
    '-buildmode=c-archive',
    '-tags', ','.join(go_tags),
    '-ldflags',
    '-X github.com/containers/toolbox/pkg/version.currentVersion='
    + meson.project_version(),
    '-o', meson.project_build_root() + '/@OUTPUT@',
  ],
  depend_files: sources,
  output: 'toolbox.a',
)

cc_exe = find_program(cc.cmd_array().get(-1))
readlink = find_program('readlink')
dirname = find_program('dirname')

libc_path = run_command(cc_exe, '--print-file-name=libc.so', capture: true, check: true).stdout()
libc_canon = run_command(readlink, '--canonicalize', libc_path, capture: true, check: true).stdout()
libc_dir = run_command(dirname, libc_canon, capture: true, check: true).stdout().strip()

message('Host machine libc:', libc_dir)

rpath = '/run/host' + libc_dir
message('Setting rpath:', rpath)

# objcopy only writes to files, but /proc/self/fd/1 is stdout so take that
objcopy = find_program('objcopy')
ldpath = run_command(
  objcopy,
  '-O', 'binary',
  '--only-section=.interp',
  cc_exe.full_path(),
  '/proc/self/fd/1',
  capture: true,
  check: true,
)

# .interp ends in a null char, so make sure to strip it
dynamic_linker = '/run/host' + ldpath.stdout().strip('\0')

message('Setting dynamic linker:', dynamic_linker)

toolbox = executable(
  'toolbox',
  sources: [],
  link_whole: toolbox_go,
  link_args: [
    '-Wl,-z,lazy',
    '-Wl,--export-dynamic',
    '-Wl,--unresolved-symbols=ignore-in-object-files',
    '-Wl,-dynamic-linker,' + dynamic_linker,
  ],
  build_rpath: rpath,
  install_rpath: rpath,
  install: true,
  install_dir: get_option('bindir'),
)

if bashcompletionsdir != ''
  custom_target(
    'bash-completion',
    capture: true,
    command: [
      toolbox,
      'completion',
      'bash',
    ],
    install: true,
    install_dir: bashcompletionsdir,
    output: 'toolbox.bash',
  )
endif

if fishcompletionsdir != ''
  custom_target(
    'fish-completion',
    capture: true,
    command: [
      toolbox,
      'completion',
      'fish',
    ],
    install: true,
    install_dir: fishcompletionsdir,
    output: 'toolbox.fish',
  )
endif

custom_target(
  'zsh-completion',
  capture: true,
  command: [
    toolbox,
    'completion',
    'zsh',
  ],
  install: true,
  install_dir: zshcompletionsdir,
  output: '_toolbox',
)

test('go fmt', meson_go_fmt_program, args: [meson.current_source_dir()])
test(
  'go vet',
  go,
  args: ['vet', '-c', '3', './...'],
  workdir: meson.current_source_dir(),
)
test(
  'go test',
  go,
  args: ['test', '-vet', 'off', './...'],
  workdir: meson.current_source_dir(),
)