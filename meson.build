project(
  'toolbox',
  'c',
  version: '0.0.99.1',
  license: 'ASL 2.0',
  meson_version: '>= 0.40.0',
)

cc = meson.get_compiler('c')
add_project_arguments('-pthread', language: 'c')
add_project_link_arguments('-pthread', language: 'c')

go = find_program('go')
go_md2man = find_program('go-md2man')
shellcheck = find_program('shellcheck', required: false)
skopeo = find_program('skopeo', required: false)

systemd_dep = dependency('systemd')
bash_completion = dependency('bash-completion', required: false)

profiledir = get_option('profile_dir')
tmpfilesdir = systemd_dep.get_pkgconfig_variable('tmpfilesdir')

if bash_completion.found()
  install_data(
    'completion/bash/toolbox',
    install_dir: bash_completion.get_pkgconfig_variable('completionsdir')
  )
endif

if not skopeo.found()
    message('Running system tests requires Skopeo for OCI image manipulation.')
endif

install_subdir(
  'test',
  install_dir: join_paths(get_option('datadir'), meson.project_name())
)

subdir('data')
subdir('doc')
subdir('libc-wrappers')
subdir('profile.d')

go_build_wrapper_file = files('go-build-wrapper')
go_build_wrapper_program = find_program('go-build-wrapper')

sources = files(
  'toolbox.go',
  'cmd/create.go',
  'cmd/enter.go',
  'cmd/help.go',
  'cmd/initContainer.go',
  'cmd/list.go',
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
