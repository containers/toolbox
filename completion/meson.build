generate_completions_program = find_program('generate_completions.py')

if bashcompletionsdir != ''
  custom_target(
    'bash-completion',
    capture: true,
    command: [generate_completions_program, meson.global_source_root() / 'src', 'bash'],
    depends: [toolbox_go],
    install: true,
    install_dir: bashcompletionsdir,
    output: 'toolbox',
)
endif

if fishcompletionsdir != ''
  custom_target(
    'fish-completion',
    capture: true,
    command: [generate_completions_program, meson.global_source_root() / 'src', 'fish'],
    depends: [toolbox_go],
    install: true,
    install_dir: fishcompletionsdir,
    output: 'toolbox.fish',
)
endif

custom_target(
  'zsh-completion',
  capture: true,
  command: [generate_completions_program, meson.global_source_root() / 'src', 'zsh'],
  depends: [toolbox_go],
  install: true,
  install_dir: zshcompletionsdir,
  output: '_toolbox',
)
