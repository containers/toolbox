project(
  'toolbox',
  'c',
  version: '0.0.95',
  license: 'ASL 2.0',
  meson_version: '>= 0.40.0',
)

cc = meson.get_compiler('c')
add_project_arguments('-pthread', language: 'c')
add_project_link_arguments('-pthread', language: 'c')

go = find_program('go')
go_md2man = find_program('go-md2man')
shellcheck = find_program('shellcheck', required: false)

profiledir = get_option('profile_dir')

systemd_dep = dependency('systemd')
tmpfilesdir = systemd_dep.get_pkgconfig_variable('tmpfilesdir')

toolbox = files('toolbox')

if shellcheck.found()
  test('shellcheck', shellcheck, args: [toolbox])
endif

bash_completion = dependency('bash-completion', required: false)
if bash_completion.found()
  install_data(
    'completion/bash/toolbox',
    install_dir: bash_completion.get_pkgconfig_variable('completionsdir')
  )
endif

install_subdir(
  'test',
  install_dir: join_paths(get_option('datadir'), meson.project_name())
)

subdir('data')
subdir('doc')
subdir('profile.d')
subdir('src')
