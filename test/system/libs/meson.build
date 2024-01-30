test_system_libs_dir = test_system_dir / 'libs'

test_system_libs_files = files(
  'helpers.bash'
)

install_data(
  install_dir: test_system_libs_dir,
  sources: test_system_libs_files
)

install_subdir(
  'bats-support',
  install_dir: test_system_libs_dir,
  exclude_files: [
    '.git',
    '.gitignore',
    '.travis.yml',
    'package.json'
  ],
  exclude_directories: [
    '.git',
    'script',
    'test'
  ]
)

install_subdir(
  'bats-assert',
  install_dir: test_system_libs_dir,
  exclude_files: [
    '.git',
    '.gitignore',
    '.travis.yml',
    'package.json'
  ],
  exclude_directories: [
    '.git',
    'script',
    'test'
  ]
)
