
## Usage

### Create the basic Fedora Toolbox image:
```
[user@hostname fedora-toolbox]$ buildah bud --tag fedora-toolbox:28 .
STEP 1: FROM docker://registry.fedoraproject.org/fedora:28
Getting image source signatures
…
…
…
[user@hostname fedora-toolbox]$
```
Modify the Dockerfile to match your taste and Fedora version. The image should
be tagged as `fedora-toolbox` with a suffix matching the host Fedora version.
eg., `fedora-toolbox:29`, etc..

### Create your Fedora Toolbox container:
```
[user@hostname fedora-toolbox]$ ./fedora-toolbox create
[user@hostname fedora-toolbox]$
```
This will create a container, and an image, called
`fedora-toolbox-<your-username>:28` that's specifically customised for your
host user.

### Enter the Toolbox:
```
[user@hostname fedora-toolbox]$ ./fedora-toolbox enter
[user@toolbox ~]$
```

