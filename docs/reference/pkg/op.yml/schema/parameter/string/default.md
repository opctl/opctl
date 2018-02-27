Default value for the parameter.

# Examples

## Wrap long line

Simple wrapping of a long line

```yaml
string:
  default: this is a really long default value and it wraps around to the
    second line it is so long. 
```

## Public key

Maintaining line breaks in a public key

```yaml
string:
  default: |
    -----BEGIN PUBLIC KEY-----
    MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgF8/xT6PYwiaJq0PfmeWzWZuEhIk
    tGiKsO6STLByCjNw/SsfqjOjH++gYHZUxDs8YJpVpLQ9L8fxSehHy3pmL76c0FhW
    RK1xE1h1aUNokoOlelt491p7LZe8XVj/xwV17YFUvhIUa7D4+PSAkwPOiam8GLKg
    wqqrZ93IXQOcxUzJAgMBAAE=
    -----END PUBLIC KEY-----
```

