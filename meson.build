project(
  'toolbox',
  'c',
  version: '0.3',
  license: 'ASL 2.0',
  default_options: 'c_std=c99',
  meson_version: '>= 0.58.0',
)

fs = import('fs')

cc = meson.get_compiler('c')
if not cc.has_argument('-print-file-name=libc.so')
  error('C compiler does not support the -print-file-name argument.')
endif

subid_dep = cc.find_library('subid', has_headers: ['shadow/subid.h'])

go = find_program('go')
go_md2man = find_program('go-md2man')

bats = find_program('bats', required: false)
codespell = find_program('codespell', required: false)
dbus_daemon = find_program('dbus-daemon', required: false)
htpasswd = find_program('htpasswd', required: false)
openssl = find_program('openssl', required: false)

p11kit_server_works = false
p11kit = find_program('p11-kit', required: false)
if p11kit.found()
  res = run_command(p11kit, 'server', check: false)
  if res.returncode() == 0
    error('Command \'p11-kit server\' was supposed to fail')
  endif

  res_stdout = res.stdout()
  if res_stdout.contains('--name') and res_stdout.contains('--provider')
    p11kit_server_works = true
  else
    warning('Command \'p11-kit server\' doesn\'t work')
  endif
endif

if not p11kit_server_works
  warning('Containers won\'t have access to the CA certificates from the host')
endif

podman = find_program('podman', required: false)
shellcheck = find_program('shellcheck', required: false)
skopeo = find_program('skopeo', required: false)

bashcompletionsdir = get_option('bash_completions_dir')
if bashcompletionsdir == ''
  bash_completion_dep = dependency('bash-completion', required: get_option('bash_completions'))
  if bash_completion_dep.found()
    bashcompletionsdir = bash_completion_dep.get_variable(pkgconfig: 'completionsdir')
  endif
endif

fishcompletionsdir = get_option('fish_completions_dir')
if fishcompletionsdir == ''
  fish_completion_dep = dependency('fish', required: get_option('fish_completions'))
  if fish_completion_dep.found()
    fishcompletionsdir = fish_completion_dep.get_variable(pkgconfig: 'completionsdir')
  endif
endif

zshcompletionsdir = get_option('zsh_completions_dir')
if zshcompletionsdir == ''
  zshcompletionsdir = get_option('datadir') / 'zsh' / 'site-functions'
endif

migration_path_for_coreos_toolbox = get_option('migration_path_for_coreos_toolbox')
profiledir = get_option('profile_dir')

tmpfilesdir = get_option('tmpfiles_dir')
if tmpfilesdir == '' or not fs.exists('/run/.containerenv')
  systemd_dep = dependency('systemd')

  if tmpfilesdir == ''
    tmpfilesdir = systemd_dep.get_variable(pkgconfig: 'tmpfilesdir')
  endif
endif

toolbox_sh = files('toolbox')

if codespell.found()
  test(
    'codespell',
    codespell,
    args: [
      '--check-filenames',
      '--check-hidden',
      '--context', '3',
      '--exclude-file', meson.project_source_root() / '.codespellexcludefile',
      '--skip', meson.project_build_root(),
      '--skip', meson.project_source_root() / '.git',
      '--skip', meson.project_source_root() / 'src/go.sum',
      '--skip', meson.project_source_root() / 'test/system/libs/bats-assert',
      '--skip', meson.project_source_root() / 'test/system/libs/bats-support',
      meson.project_source_root(),
    ],
  )
endif

if shellcheck.found()
  test('shellcheck toolbox (deprecated)', shellcheck, args: [toolbox_sh])
endif

install_subdir(
  'test',
  install_dir: get_option('datadir') / meson.project_name(),
  exclude_files: [
    'meson.build',
    'system/libs/bats-assert/.git',
    'system/libs/bats-assert/.gitignore',
    'system/libs/bats-assert/.travis.yml',
    'system/libs/bats-assert/package.json',
    'system/libs/bats-support/.git',
    'system/libs/bats-support/.gitignore',
    'system/libs/bats-support/.travis.yml',
    'system/libs/bats-support/package.json',
    'system/meson.build',
    ],
  exclude_directories: [
    'system/libs/bats-assert/script',
    'system/libs/bats-assert/test',
    'system/libs/bats-support/script',
    'system/libs/bats-support/test'
  ]
)

subdir('data')
subdir('doc')
subdir('profile.d')
subdir('src')
subdir('test')

meson.add_install_script('meson_post_install.py')
