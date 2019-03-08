go_md2man_command = [
  go_md2man,
  '-in', '@INPUT@',
  '-out', '@OUTPUT@',
]

manuals = [
  'toolbox.1',
  'toolbox-create.1',
  'toolbox-enter.1',
  'toolbox-list.1',
  'toolbox-rm.1',
  'toolbox-rmi.1',
  'toolbox-run.1',
]

foreach manual: manuals
  input = manual + '.md'
  output = manual

  custom_target(
    output,
    command: go_md2man_command,
    input: input,
    install: true,
    install_dir: join_paths(get_option('mandir'), 'man1'),
    output: output,
  )
endforeach
