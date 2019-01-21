# Fedora Toolbox — Hacking on OSTree-based Fedoras

[Fedora Toolbox](https://github.com/debarshiray/fedora-toolbox) is a tool that
offers a familiar RPM based environment for developing and debugging software
on locked down [OSTree](https://ostree.readthedocs.io/en/latest/) based Fedora
systems like [Silverblue](https://silverblue.fedoraproject.org/). Such
operating systems are shipped as *immutable* OSTree images, where it's
difficult to setup a development environment with your favorite tools, editors
and SDKs. A toolbox container solves that problem by providing a RPM based
*mutable* container. You can tweak it to your heart's content and use DNF to
install your favorite packages, all without worrying about breaking your
operating system.

The toolbox environment is based on the `fedora-toolbox` image. This image is
then customized for the current user to create a toolbox container that
seamlessly integrates with the rest of the operating system.

## Usage

### Create your Fedora Toolbox container:
```
[user@hostname ~]$ fedora-toolbox create
[user@hostname ~]$
```
This will create a container, and an image, called
`fedora-toolbox-<your-username>:<version-id>` that's specifically customised
for your host user.

### Enter the Toolbox:
```
[user@hostname ~]$ fedora-toolbox enter
🔹[user@toolbox ~]$
```

