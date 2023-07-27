test_system = files(
  '201-ipc.bats',
  '203-network.bats',
  '206-user.bats',
  '210-ulimit.bats',
  '211-dbus.bats',
)

if shellcheck.found()
  test('shellcheck test/system', shellcheck, args: [test_system])
endif
