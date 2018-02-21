Icon for the pkg.

MUST be an absolute path (from pkg root) to a valid
[v1.1 SVG](https://www.w3.org/TR/SVG11/) w/ a 1:1 aspect ratio (equal
height & width)

> SVG is a vector graphic format; it looks the same at any resolution.

## Examples

### Pkg root

icon is at root of pkg

```yaml
name: iconAtPkgRoot
icon: /icon.svg
run:
  container:
    image: { ref: alpine }
```

### Sub directory

icon is in sub-dir of pkg

```yaml
name: iconInSubDir
icon: /sub-dir/icon.svg
run:
  container:
    image: { ref: alpine }
```

