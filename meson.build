project(
  'toolbox',
  version: '0.0.16',
  license: 'ASL 2.0',
  meson_version: '>= 0.40.0',
)

go_md2man = find_program('go-md2man')
shellcheck = find_program('shellcheck', required: false)

profiledir = get_option('profile_dir')

systemd_dep = dependency('systemd')
tmpfilesdir = systemd_dep.get_pkgconfig_variable('tmpfilesdir')

toolbox = files('toolbox')

if shellcheck.found()
  test('shellcheck', shellcheck, args: [toolbox])
endif

install_data(
  toolbox,
  install_dir: get_option('bindir'),
)

bash_completion = dependency('bash-completion', required: false)
if bash_completion.found()
  install_data(
    'completion/bash/toolbox',
    install_dir: bash_completion.get_pkgconfig_variable('completionsdir')
  )
endif

subdir('data')
subdir('doc')
subdir('profile.d')

#mytarget = custom_target('toolbox-rs',
#  output : 'toolbox-rs/target/release/toolbox-rs',
#  input : 'toolbox-rs/src/coretoolbox.rs',
#  command : ['/bin/sh', '-c', 'cd toolbox-rs && cargo build --release && ln target/release'],
#  install : true,
#  install_dir : 'libexec/toolbox-rs')

