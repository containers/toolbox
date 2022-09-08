test('toolbox go binary tests',
  go,
  args: ['test', '-v', './...', '-args', toolbox.full_path(), meson.project_build_root() / 'src' / 'go-build-overrides'],
  depends: toolbox,
  workdir: meson.current_source_dir(),
)
