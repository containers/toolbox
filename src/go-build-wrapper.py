#!/usr/bin/env python3
#
# Copyright © 2022 Ondřej Míchal
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

import os
import shlex
import subprocess
import sys

prog_name = os.path.basename(sys.argv[0])


def get_canonical_libc_dirname(cc):
    try:
        libc_path = subprocess.run(
            [cc, '--print-file-name=libc.so'],
            capture_output=True, check=True, text=True).stdout.strip()
    except subprocess.CalledProcessError as e:
        print(f'{prog_name}: failed to read the path to libc.so', file=sys.stderr)
        print(e.stderr, file=sys.stderr)
        sys.exit(1)

    libc_path_canonical = os.path.realpath(libc_path)

    return os.path.dirname(libc_path_canonical)


def create_rpath(libc_path_dirname):
    return '/run/host' + libc_path_dirname


def alter_dynamic_linker_path(dynamic_linker):
    dynamic_linker_basename = os.path.basename(dynamic_linker)
    dynamic_linker_canonical = os.path.realpath(dynamic_linker)
    dynamic_linker_canonical_dirname = os.path.dirname(
        dynamic_linker_canonical)

    return f'/run/host{dynamic_linker_canonical_dirname}/{dynamic_linker_basename}'


if len(sys.argv) != 8:
    print('{}: wrong arguments'.format(prog_name), file=sys.stderr)
    print('''Usage: {} [SOURCE DIR]
                           [OUTPUT DIR]
                           [VERSION]
                           [C COMPILER]
                           [DYNAMIC LINKER]
                           [MIGRATION PATH FORCOREOS/TOOLBOX]
                           [BUILD FLAGS]'''.format(prog_name), file=sys.stderr)
    sys.exit(1)

source_dir = sys.argv[1]
output_dir = sys.argv[2]
version = sys.argv[3]
cc = sys.argv[4]
dynamic_linker = sys.argv[5]
# Boolean is passed from Meson as a string
coreos_migration = True if sys.argv[6] == "true" else False
build_flags = sys.argv[7]

try:
    os.chdir(source_dir)
except OSError as e:
    print(f'{prog_name}: failed to change directory to {source_dir}',
          file=sys.stderr)
    print(e, file=sys.stderr)
    sys.exit(1)

libc_path_dirname = get_canonical_libc_dirname(cc)
dynamic_linker = alter_dynamic_linker_path(dynamic_linker)
rpath = create_rpath(libc_path_dirname)

# Write overriden binary values to a file for testing
go_build_override_values_file = os.path.join(output_dir, 'go-build-overrides')
with open(go_build_override_values_file, mode='w') as f:
    f.write(f'{dynamic_linker} {rpath}')

build_cmd = [
    'go', 'build',
    '-trimpath',
    '-o', os.path.join(output_dir, 'toolbox')
]

if coreos_migration:
    build_cmd.extend(['-tags', 'migration_path_for_coreos_toolbox'])

if build_flags:
    build_cmd.extend(shlex.split(build_flags))

build_cmd.extend([
    '-ldflags',
    f'-extldflags \'-Wl,-dynamic-linker,{dynamic_linker} -Wl,-rpath,{rpath}\' -linkmode external -X github.com/containers/toolbox/pkg/version.currentVersion={version}'
])

go_build = subprocess.run(build_cmd)
sys.exit(go_build.returncode)
