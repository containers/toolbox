![Toolbox logo](data/logo/toolbox-logo-landscape.svg)

# Contributing to Toolbox

Thank you for wanting to contribute to Toolbox! We greatly appreciate your
interest!

# Table of contents

- [Contributing to Toolbox](#contributing-to-toolbox)
- [Table of contents](#table-of-contents)
- [Reporting Bugs](#reporting-bugs)
  - [Before Submiting a Bug Report](#before-submiting-a-bug-report)
  - [Writing a Bug Report](#writing-a-bug-report)
- [Making Suggestions](#making-suggestions)
  - [Before Submitting a Suggestion](#before-submitting-a-suggestion)
  - [Writing a Suggestion](#writing-a-suggestion)
- [First Contribution](#first-contribution)
- [Pull Requests](#pull-requests)
  - [Creating a Pull Request](#creating-a-pull-request)
  - [After Creating a Pull Request](#after-creating-a-pull-request)
- [Little Style Guide](#little-style-guide)

# Reporting Bugs

## Before Submiting a Bug Report

- Check if your issue is already reported in our [bug tracker](https://github.com/containers/toolbox/issues)
  - If the issue is already reported and is marked as **OPEN**, comment on it
    and if possible and needed, share info about the issue just as if you were
    submiting a new issue
  - If the issue is marked as **CLOSED**, check if your version of Toolbox is
    up-to-date or if there are some steps, described in the closed issue, that
    you should follow. If you are still experiencing the issue, please file a
    new issue
- See our [documentation](https://docs.fedoraproject.org/en-US/fedora-silverblue/toolbox/)
  if there are some steps that could help you solve your issue
- Sometimes a bug is not reported in our bug tracker but instead people ask for
  help somewhere else: IRC ([Freenode](https://freenode.net) - #silverblue,
  #containers, #fedora, #fedora-devel,..), [Fedora forum](https://discussion.fedoraproject.org/tag/toolbox),
  or somewhere else. In such cases we'd like you to still report the bug and
  share with us any info that could be gathered from those places

## Writing a Bug Report

Writing good bug reports is a nice way to make the job of the maintainers and
other contributors a bit easier.

When writing a bug report:

- **Use a clear and descriptive title**
- **Describe the problem** - Can you reproduce the bug reliably? What first
  triggered the problem? Did it start happening after upgrading your system?
- **Provide steps how to reproduce** - It's easier for us to fix a bug if we can
  reproduce it.
- **Describe the behavior you received and what you expected** - Sometimes it
  may not be clear what the *right* behavior should look like.
- **Provide info about the version of used software** - What version of Toolbox
  and Podman do you use?
- **Provide info about your system** - What distribution do you use? Which
  desktop environment? Is it a VM or a real machine?

# Making Suggestions

Toolbox is not feature-complete and some of it's functionality is not-there-yet.
We are thankful for all suggestions and ideas but be ready that your suggestion
may be rejected.

## Before Submitting a Suggestion

- Check if your suggestion has not already been made in our [bug tracker](https://github.com/containers/toolbox/issues)
  - If it has and is marked as **OPEN**, go ahead and share your own thoughts
    about the topic!
  - If it has and is marked as **CLOSED**, please read the ticket and depending
    on whether the suggestion was accepted or not consider if it is worth
    opening a new issue or not.
- Consider if the suggestion is not too out of scope of the project.

## Writing a Suggestion

When writing a suggestion:

- **Use a clear and descriptive title**
- **Describe the idea** - What parts of Toolbox does it affect? Is it a major
  functionality or a minor tweak?
- **Provide step-by-step description of the suggested behavior** so that we
  will understand.
- **Explain why would this idea be useful** - It sounds good to have a lot of
  options but sometimes less is more. See this [article](https://ometer.com/preferences.html).

# First Contribution

Toolbox is written in [Go](https://golang.org) and uses [Meson](https://mesonbuild.com)
as it's buildsystem.

Instructions for building Toolbox from source are in our [README](https://github.com/containers/toolbox/blob/main/README.md).

> You may not need to build the project from source if your contribution is not
> related to the code of Toolbox itself (e.g., documentation, updating CI
> config, playing with image definitions,...).

Here are some ideas of what you could contribute with:

- Check our [bug tracker](https://github.com/containers/toolbox/issues)
  and look for tickets marked with labels `good-first-issue` or `help-wanted`.
- Write tests - Go has [tools](https://golang.org/pkg/testing/) for writing tests.
  There are also [some](https://github.com/stretchr/testify) [libraries](https://github.com/onsi/ginkgo)
  used for creating even more sophisticated tests.
- Play with custom images - Toolbox currently officially works with Fedora-based
  images. Ultimately there should be a wide variety of supported distro images.
  You can help with testing other people's image definitions or creating your
  own. **Beware**, maintainers still don't have a clear idea of how the image
  infrustructure should look like.
- Write documentation - Some functions in Toolbox's code don't have comments and
  it's not very clear what they do. Toolbox has it's [documentation](https://docs.fedoraproject.org/en-US/fedora-silverblue/toolbox/)
  hosted by Fedora. It's not very large and could use some attention.
- Hack on the code and share the result - Seriously! Sometimes random ideas are
  the best.

Toolbox currently does not have an infrastructure for translations. You can help
us to set it up!

# Pull Requests

All pull requests are welcome! Features, bug fixes, fixing of typos, tests,
documentation, code comments and much more.

## Creating a Pull Request

- Document well your changes - This applies to the description of your PR and to
  your commit messages.
- If possible add additional test cases - If there are no tests for the part of
  code you're contributing to, consider opening another PR if you want to
  implement it yourself or file an issue so that somebody else can pick it up.
- Update documentation to reflect your changes - Manual pages can be found in
  directory `doc`. If your changes affect Toolbox's [documentation](https://docs.fedoraproject.org/en-US/fedora-silverblue/toolbox/),
  consider creating a PR there (but to save yourself time, you can do it
  after your changes are accepted), too.

## After Creating a Pull Request

It may take the us some time to review your changes and sometimes even longer to
actually merge them. Please, don't interpret this as an act of not appreciating
your efforts! We really appreciate them! Sometimes we may be stuck in different
parts of our lives.

If it takes us a very long time to even respond to your Pull Request, you can
try to @ping us, request a review or try to reach to us on IRC ([Freenode](https://freenode.net/);
#silverblue, #containers, #fedora-devel,..) or [Fedora Forum](https://discussion.fedoraproject.org).

Toolbox has a simple CI (Continuos Integration) setup for running system tests (
can be found under directory `test/system`). Their goal is to check if your
changes don't affect adversely Toolbox's functionality. Sometimes these tests
mail fail with a false-positive. If you are not sure about the outcome of the
tests, reach out to the maintainers!

Toolbox's CI system is [Zuul](https://zuul-ci.org/) hosted at [softwarefactory](https://softwarefactory-project.io/).

# Little Style Guide

Toolbox is written in [Go](https://golang.org) and uses its default set of tools
including `gofmt` and `golint`.

Here are some good materials to learn from about the way how to write nice and
idiomatic code in Go:

- [A Tour of Go](https://tour.golang.org/welcome)
- [How To Write Go Code](https://golang.org/doc/code.html)
- [Effective Go](https://golang.org/doc/effective_go.html)

Overall, the [Go Blog](https://blog.golang.org/) is a good place to learn more
about Go.

If you are using Visual Studio Code, there are [plugins](https://marketplace.visualstudio.com/items?itemName=golang.Go)
that include all this functionality and throw a warning if you're doing
something wrong.

