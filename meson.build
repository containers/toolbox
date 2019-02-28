project(
  'toolbox',
  version: '0.0.6',
  license: 'ASL 2.0',
  meson_version: '>= 0.40.0',
)

go_md2man = find_program('go-md2man')

install_data(
  'toolbox',
  install_dir: get_option('bindir'),
)

subdir('doc')
