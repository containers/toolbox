#!/usr/bin/env python3
#
# Copyright © 2022 Ondřej Míchal
# Copyright © 2022 Red Hat Inc.
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
import subprocess
import sys

if len(sys.argv) != 3:
    print('{}: wrong arguments'.format(sys.argv[0]), file=sys.stderr)
    print('Usage: {} [SOURCE DIR] [COMPLETION TYPE]'.format(sys.argv[0]), file=sys.stderr)
    sys.exit(1)

source_dir = sys.argv[1]
completion_type = sys.argv[2]

os.chdir(source_dir)
output = subprocess.run(['go', 'run', '.', 'completion', completion_type], check=True)

sys.exit(0)
