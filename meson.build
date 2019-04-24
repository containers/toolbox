project(
  'toolbox',
  version: '0.0.8',
  license: 'ASL 2.0',
  meson_version: '>= 0.40.0',
)

go_md2man = find_program('go-md2man')
shellcheck = find_program('shellcheck', required: false)

profiledir = get_option('profile_dir')

systemd_dep = dependency('systemd')
tmpfilesdir = systemd_dep.get_pkgconfig_variable('tmpfilesdir')

toolbox = files('toolbox')

if shellcheck.found()
  test('shellcheck', shellcheck, args: [toolbox])
endif

install_data(
  toolbox,
  install_dir: get_option('bindir'),
)

subdir('data')
subdir('doc')
subdir('profile.d')
