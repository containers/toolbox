This started out as a fork of the (now) containers/toolbox shell script
code.  Currently, it can be used opt-in.

Installation
---

Be sure you have [cargo installed](https://doc.rust-lang.org/cargo/getting-started/installation.html).

Then:
`cargo install --path .`

In the future we may invest in packaging this for different distributions, or
see about shipping it with e.g. podman by default.

Getting started
---

One time setup

```
$ toolbox create
<answer questions>
$
```

Now, each time you want to enter the toolbox:

```
$ toolbox run
```

One suggestion is to add a "profile" or configuration to your terminal
emulator that runs `coretoolbox run` by default, so that you can
easily create new tabs/windows in the toolbox.
