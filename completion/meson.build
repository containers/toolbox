generate_completions_program = find_program('generate_completions.py')

if bash_completion_dep.found()
  bashcompletionsdir = bash_completion_dep.get_variable(pkgconfig: 'completionsdir')
else
  bashcompletionsdir = get_option('datadir') / 'bash-completion' / 'completions'
  message('bash-completion not found: using', get_option('prefix') / bashcompletionsdir, 'as a falback install directory')
endif

if fish_dep.found()
  fishcompletionsdir = fish_dep.get_variable(pkgconfig: 'completionsdir')
else
  fishcompletionsdir = get_option('datadir') / 'fish' / 'completions'
  message('fish not found: using', get_option('prefix') / fishcompletionsdir, 'as a fallback install directory')
endif

custom_target(
  'bash-completion',
  capture: true,
  command: [generate_completions_program, meson.global_source_root() / 'src', 'bash'],
  depends: [toolbox_go],
  install: true,
  install_dir: bashcompletionsdir,
  output: 'toolbox'
)

custom_target(
  'zsh-completion',
  capture: true,
  command: [generate_completions_program, meson.global_source_root() / 'src', 'zsh'],
  depends: [toolbox_go],
  install: true,
  install_dir: get_option('datadir') / 'zsh' / 'site-functions',
  output: '_toolbox'
)

custom_target(
  'fish-completion',
  capture: true,
  command: [generate_completions_program, meson.global_source_root() / 'src', 'fish'],
  depends: [toolbox_go],
  install: true,
  install_dir: fishcompletionsdir,
  output: 'toolbox.fish'
)
