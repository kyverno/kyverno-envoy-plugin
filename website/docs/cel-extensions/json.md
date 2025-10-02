# JSON library

The JSON CEL library provides functions for parsing JSON strings into CEL values.

## Functions

### json.Unmarshal

The `json` function parses a JSON-encoded string and returns the corresponding CEL value. This allows you to work with JSON data directly in CEL expressions.

#### Signature

```
json.Unmarshal(<string> jsonString) -> Map[string][any]
```

### Example

```
# Check JSON nested value
json.Unmarshal("{\"item1\": { \"item2\": 123 } \"}").item1.item2 == 123
```