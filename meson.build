project(
  'toolbox',
  'c',
  version: '0.0.99.3',
  license: 'ASL 2.0',
  meson_version: '>= 0.58.0',
)

fs = import('fs')

cc = meson.get_compiler('c')
if not cc.has_argument('-print-file-name=libc.so')
  error('C compiler does not support the -print-file-name argument.')
endif

go = find_program('go')
go_md2man = find_program('go-md2man')
shellcheck = find_program('shellcheck', required: false)
skopeo = find_program('skopeo', required: false)

bash_completion = dependency('bash-completion', required: false)

migration_path_for_coreos_toolbox = get_option('migration_path_for_coreos_toolbox')
profiledir = get_option('profile_dir')

tmpfilesdir = get_option('tmpfiles_dir')
if tmpfilesdir == '' or not fs.exists('/run/.containerenv')
  systemd_dep = dependency('systemd')

  if tmpfilesdir == ''
    tmpfilesdir = systemd_dep.get_variable(pkgconfig: 'tmpfilesdir')
  endif
endif

if bash_completion.found()
  install_data(
    'completion/bash/toolbox',
    install_dir: bash_completion.get_variable(pkgconfig: 'completionsdir')
  )
endif

if not skopeo.found()
    message('Running system tests requires Skopeo for OCI image manipulation.')
endif

install_subdir(
  'test',
  install_dir: join_paths(get_option('datadir'), meson.project_name()),
  exclude_files: [
    'system/libs/bats-assert/.git',
    'system/libs/bats-assert/.gitignore',
    'system/libs/bats-assert/.travis.yml',
    'system/libs/bats-assert/package.json',
    'system/libs/bats-support/.git',
    'system/libs/bats-support/.gitignore',
    'system/libs/bats-support/.travis.yml',
    'system/libs/bats-support/package.json'
    ],
  exclude_directories: [
    'system/libs/bats-assert/.git',
    'system/libs/bats-assert/script',
    'system/libs/bats-assert/test',
    'system/libs/bats-support/.git',
    'system/libs/bats-support/script',
    'system/libs/bats-support/test'
  ]
)

subdir('data')
subdir('doc')
subdir('profile.d')
subdir('src')

meson.add_install_script('meson_post_install.py')
