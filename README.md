# Char

Char is an organizations best friend, it helps with code sharing when either
using a multi-repo strategy or mono repo. The goal of the project is to
facilitate sharing of presets and plugins to reduce the required complexity
developers have to take on.

This project is best suited with a standard library of plugins (which serves the
function of libraries), as well `kjuulh/bust` which is a platform agnostic
CI/task setup

This is in very early stages, and for now it officially supports scaffolding.

## Example

The `examples` folder shows how to load plugins, though presets are still
pending.

```yaml
# file: .char.yml
registry: git.front.kjuulh.io
plugins:
  "kjuulh/char#/plugins/gocli":
    vars:
      name: "char"
```
