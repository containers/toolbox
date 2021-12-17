#!/usr/bin/env python3

import os
import subprocess
import sys

destdir = os.environ.get('DESTDIR', '')

if not destdir and not os.path.exists('/run/.containerenv'):
    print('Calling systemd-tmpfiles --create ...')

    try:
        subprocess.run(['systemd-tmpfiles', '--create'], check=True)
    except subprocess.CalledProcessError as e:
        print('Returned non-zero exit status', e.returncode)
        sys.exit(e.returncode)

sys.exit(0)
