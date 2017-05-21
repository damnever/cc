## cc is a very flexible library for configuration management.

[![Build Status](https://travis-ci.org/damnever/cc.svg?branch=master)](https://travis-ci.org/damnever/cc) [![GoDoc](https://godoc.org/github.com/damnever/cc?status.svg)](https://godoc.org/github.com/damnever/cc)

Only support JSON and YAML.

### Installation

```
go get -u github.com/damnever/cc
```

### Usage

```go
c, _ := cc.NewConfigFromFile("./example/example.yaml")  // file must has extension
_ := c.MergeFromFile("./example/example.json") // do not ignore the errors

c.Must("name")  // panic if not found
c.String("name")

cc := c.Config("map")
cc.Bool("key_one")

list := c.Value("list").List()
list[1].Int()

// environment variables
os.Setenv("float_env", "11.11")
c.Float("float_env")
```

#### Default configs

We may wrote the code like this:
```go
name := "default"
if c.Has("name") {
    name = c.String("name")  // or panic
}
```

Now, we can write code like this:
```go
name := c.StringOr("name", "cc")  // or c.Must("name")
b := c.BoolOr("bool", true)
```

#### Pattern && Validation

If you want check string value whether it is matched by `regexp`:
```go
s, ok := c.StringAnd("name", "^c")
```
Or, the make the string value as a pattern:
```go
p := c.Pattern("pattern_key_name")
ok := p.ValidateString("a string")
```

For `int` and `float`, `cc` use if-like condition to do similar work.
Assume we have `threhold: "N>=30&&N<=80"` in config file, we can use it like this:
```go
p := c.Pattern("threhold")
ok := p.ValidateInt(40)  // or ValidateFloat
```
Or, using a pattern to validate the number:
```go
ni, ok := c.IntAnd("int_key", "N>50")
nf, ok := c.FloatAnd("float_key", "N/100>=0.3")
```
NOTE: bit operation is not supported.


### LICENSE

[The BSD 3-Clause License](https://github.com/damnever/cc/blob/master/LICENSE)
