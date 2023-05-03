- [Design goals](#design-goals)
- [Install](#install)
- [Rule](#rule)
  - [Basic Types](#basic-types)
  - [Constraint](#constraint)
- [Example](#example)
- [TODO](#todo)

## Design Goals

## Install
```shell
go get -u github.com/xuchengeu/invalid
```


## Rule

### Basic Types

- `$obj`  : field object contains sub-fields
- `$str`  : string
- `$bool`: boolean
- `$arr`  : alike type “Array” in Java, contains sub-fields in only one type.
- `$num`  : floating point
- `$int`  : integer
- `$nil`  : NULL value, NULL value’s different from empty string. NULL represent nil in Golang

### Constraint

- `$required` :  $required means fields must exist, $required could be omitted which means fields is required for default.
- `$optional` :  $optional means fields could be omitted.
- `$length` : length of character, valid under type `$str`
- `$reg` : regexp pattern written in string, valid under type `$str`
- `$length.$min` : minimum length of string, valid under constraint `$length`
- `$length.$max` : maximum length of string, valid under constraint `$length`
- `$key-reg` : a regexp written in string to perform key name validation.It can be used in scenario like checking extensible keys only prefix with ‘x’ in Swagger, `key-reg` valid in type `$obj`
- `$constraint` : a type constraint for type $arr , valid for type `$arr`. value of constraint could be a valid basic type or map. checkout array example for more reference.


## Example

### YAML source
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

### YAML Rule
```yaml
---
map: 
  $type: $obj
  str1: 
    $type: $str
    $length: 
      $min: 6
      $max: 12
  bool:
    $type: $bool
  num:
    $type: $float
  int:
    $type: $int
  
list:
  type: $arr
  constraint: $str
data2:
  $type: $obj 
  map3:
    $type: $obj
    map4:
      $type: $str
```

### Code
```gotemplate
    file, err := os.Open(filepath.Join([]string{"your","path","here"}...))
    if err != nil{
	    log.Println(err)
        return
    }
    field, err := NewYAML(file)
    if err != nil{
        log.Println(err)
        return
    }
	
    file, err = os.OpenFile(filepath.Join([]string{"your","path","here"}...), 
            os.O_RDONLY, os.ModeSticky)
    if err != nil{
        log.Println(err)
        return
    }
    rule, err := NewRule(file)
    if err != nil{
        log.Println(err)
        return
    }

    errs := rule.Validate(field)
    log.Println(errs)
```


## TODO

- `$any`  : represent any valid type
- `$seq`  : `$seq` can contain any type as content
- `$range` : range of number or int
- `of` : constraint of valid value in enumeration value, valid under type `$str` ,`$int` ,`$num` or `$any`
- external reference: feature like  "Anchor" & "Extend/Inherit" in YAML Spec 1.2 is an available option.
- Implicit variable declaration, like declaration for type `$obj`. which makes rules more clear.