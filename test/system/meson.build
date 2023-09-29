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
  'libs/helpers.bash',
)

if shellcheck.found()
  test('shellcheck test/system', shellcheck, args: [test_system])
endif
