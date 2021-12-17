go_md2man_command = [
  go_md2man,
  '-in', '@INPUT@',
  '-out', '@OUTPUT@',
]

manuals = {
  '1': [
    'toolbox',
    'toolbox-create',
    'toolbox-enter',
    'toolbox-init-container',
    'toolbox-help',
    'toolbox-list',
    'toolbox-rm',
    'toolbox-rmi',
    'toolbox-run',
  ],
  '5': [
    'toolbox.conf',
  ]
}

foreach section, pages: manuals
  foreach page: pages
    output = page + '.' + section
    input = output + '.md'
    sectiondir = 'man' + section

    custom_target(
      output,
      command: go_md2man_command,
      input: input,
      install: true,
      install_dir: join_paths(get_option('mandir'), sectiondir),
      output: output,
    )
  endforeach
endforeach
