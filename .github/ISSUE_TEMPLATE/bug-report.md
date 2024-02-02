---
name: Bug report
about: Toolbx's bug report template
title: ''
labels: 1. Bug
assignees: ''

---

**Describe the bug**
A clear and concise description of *what the bug is*. If possible, re-run the command(s) with `--log-level debug` and put the output here.

**Steps how to reproduce the behaviour**
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected behaviour**
A clear and concise description of what you *expected to happen*.

**Actual behaviour**
A clear and concise description of what *actually happened*.

**Screenshots**
If applicable, add screenshots to help explain your problem.

**Output of `toolbox --version` (v0.0.90+)**
e.g., `toolbox version 0.0.90`

**Toolbx package info (`rpm -q toolbox`)**
e.g., `toolbox-0.0.18-2.fc32.noarch`

**Output of `podman version`**
e.g.,
```
Version:            1.9.2
RemoteAPI Version:  1
Go Version:         go1.14.2
OS/Arch:            linux/amd64
```

**Podman package info (`rpm -q podman`)**
e.g., `podman-1.9.2-1.fc32.x86_64`

**Info about your OS**
e.g., Fedora Silverblue 32

**Additional context**
Add any other context about the problem here.
When did the issue start occurring? After an update (what packages were updated)?
If the issue is about operating with containers/images (creating, using, deleting,..), share here what image you used. If you're unsure, share here the output of `toolbox list -i` (shows all Toolbx images on your system).

If you see an error message saying: `Error: invalid entry point PID of container <name-of-container>`, add to the ticket output of command `podman start --attach <name-of-container>`.
