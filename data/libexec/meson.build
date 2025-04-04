host_runner = files('host-runner')
sudo_wrapper = files('sudo-wrapper')

if shellcheck.found()
  test('shellcheck', shellcheck, args: [host_runner, sudo_wrapper])
endif

install_data(
  host_runner,
  sudo_wrapper,
  install_dir: get_option('libexecdir'),
)
