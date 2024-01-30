test_system_dir = get_option('datadir') / meson.project_name() / 'test'

awk = find_program('awk', required: false)
bats = find_program('bats', required: false)
htpasswd = find_program('htpasswd', required: false)
openssl = find_program('openssl', required: false)
podman = find_program('podman', required: false)
skopeo = find_program('skopeo', required: false)

test_system = files(
  '001-version.bats',
  '002-help.bats',
  '101-create.bats',
  '102-list.bats',
  '103-container.bats',
  '104-run.bats',
  '105-enter.bats',
  '106-rm.bats',
  '107-rmi.bats',
  '108-completion.bats',
  '201-ipc.bats',
  '203-network.bats',
  '206-user.bats',
  '210-ulimit.bats',
  '211-dbus.bats',
  '220-environment-variables.bats',
  'setup_suite.bash',
)

if shellcheck.found()
  test('shellcheck test/system', shellcheck, args: [test_system])
endif

if fs.exists('/run/.containerenv')
  message('System tests can not be run in a container.')
elif not awk.found() or not bats.found() or not htpasswd.found() or not openssl.found() or not podman.found() or not skopeo.found()
  message('System tests require: awk, bats, htpasswd, openssl, podman and skopeo being installed.')
else
  system_tests_env = environment()
  system_tests_env.set('TOOLBOX', toolbox_go.full_path())

  test(
    'system tests',
    bats,
    args: [
      '--formatter', 'tap',
      '--timing',
      './'],
    depends: toolbox_go,
    env: system_tests_env,
    is_parallel: false,
    protocol: 'tap',
    suite: 'system',
    timeout: 3600,
    workdir: meson.current_source_dir(),
  )
endif

install_data(
  install_dir: test_system_dir,
  sources: [test_system, 'README.md'],
)

subdir('libs')
