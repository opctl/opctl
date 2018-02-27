Default value for the parameter.

# Examples

## Array of literals

```yaml
array:
  default:
  - item1
  - true
  - 2.2
```

## Object in array

```yaml
array:
  default:
  - item1Prop1: item1Prop1Value
```

## Array in array

```yaml
array:
  default:
  - item1:
    - item2
```
