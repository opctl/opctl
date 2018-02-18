Human friendly identifier for the pkg.

> It's considered good practice to make `name` unique by using domain
> &/or path based namespacing.

Packages MAY be network resolvable; therefore `name` MUST be a valid
[uri-reference](https://tools.ietf.org/html/rfc3986#section-4.1)

## Examples

### With namespacing

```yaml
name: `github.com/opspec-pkgs/jwt.encode`
run:
  container:
    image: { ref: alpine }
```
