toolbox_sh = files('toolbox.sh')

if shellcheck.found()
  test('shellcheck profile.d', shellcheck, args: ['--shell=sh', toolbox_sh])
endif

install_data(
  toolbox_sh,
  install_dir: profiledir,
)
