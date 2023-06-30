test_system = files(
  '201-ipc.bats',
  '203-network.bats',
  '210-ulimit.bats',
)

if shellcheck.found()
  test('shellcheck test/system', shellcheck, args: [test_system])
endif
