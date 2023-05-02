# Schema Validation

## Schema Rule

### Basic Types

- `$obj`  : field object contains sub-fields
- `$str`  : string
- `$bool`: boolean
- `$arr`  : alike type “Array” in Java, contains sub-fields in only one type.
- `$num`  : floating point
- `$int`  : integer
- `$nil`  : NULL value, NULL value’s different from empty string. NULL represent nil in Golang

### Constraint

- `required` :  fields must exist,  exists under type `$obj`
- `optional` :  fields which are optional,  exists under type `$obj` alike `required`
- `length` : length of character, valid under type `$obj`
- `reg` : regexp pattern written in string, valid in type `$str`
- `length.min` : minimum length of string, valid in type `$str`
- `length.max` : maximum length of string, valid in type `$str`
- `key-reg` : a regexp written in string to perform key name validation.It can be used in scenario like checking extensible keys only prefix with ‘x’ in Swagger, `key-reg` valid in type `$obj`
- `constraint` : a type constraint for type $arr , valid for type `$arr`

### TODO

- `$any`  : represent any valid type 🌶️
- `$seq`  : sequence block in YAML specification, `$seq` can contain any type as content
- `$range` : range of number or int
- `of` : constraint of valid value in enumeration value, valid in type `$str` ,`$int` ,`$num` or `$any`
- external reference, feature like  "Anchor" & "Extend/Inherit" in YAML Spec 1.2 is an available option.
- Implicit variable declaration

### Implementation Schedule

### Reference

### Example

```yaml
---

map:
  str1: value2
  bool: true
  num: 12e3
  int: 20
list:
  - list_value1
  - list_value2
  - list_value3
data2:
  map3:
    map4: value4
```

```yaml
---

map: 
  type: $obj
  str1: 
    type: $str
    length: 
      min: 6
      max: 12
  bool:
    type: $bool
  num:
    type: $float
  int:
    type: $int
  
list:
  type: $arr
  constraint: $str
data2:
  type: $obj
  optional: 
map3:
    type: $obj
    map4:
      type: $str
```