Human friendly description of the pkg.

Markdown can be used; see [markdown](../markdown.md) for details.

## Examples

### Four liner
Four line description leveraging
[yaml multiline strings](http://yaml-multiline.info/)

```yaml
name: descriptionIsFourLines
description: |
 line 1
 line 2
 
 line 3
run:
  container:
    image: { ref: alpine }
```
