### Wo is a fork of Go

Since Go 1.24 isn't supported on my machine, please see the 1.23 branch https://github.com/wo-language/wo/tree/release-branch.go1.23.wo

The Wo language is an interoperable successor to Go that offers alternative syntax and language features aimed at readability.

For example,

```go
f, err := os.Open("hi.wo")
if err != nil {
    return nil, err
}
```

would be done like this in Wo:

```go
var file = os.Open("hi.wo")! // pending decisions here; it's a WIP
```

(...description continues at https://github.com/wo-language/wo/tree/release-branch.go1.23.wo)
