profile_toolbox_sh = files('toolbox.sh')

if shellcheck.found()
  test('shellcheck profile.d/toolbox.sh', shellcheck, args: ['--shell=sh', profile_toolbox_sh])
endif

install_data(
  profile_toolbox_sh,
  install_dir: profiledir,
)
