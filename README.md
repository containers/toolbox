Toolbox website
===============

To set up jekyll locally on Fedora:

install `ruby`, `rubygems` and `rubygem-bundler`

In the git checkout, do a `bundle install`. This installs all 
the needed modules in their appropriate vesions.

To run a local web server to test the site:

```
bundle exec jekyll s
```

The CI should deploy the site automatically.