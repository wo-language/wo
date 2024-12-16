### Wo is a fork of Go

The Wo language offers an alternative syntax and functionality to the Go programming language. It aims to be interoperable with Go.

Here's one example:

```go
f, err := os.Open("open_me.go")
if err != nil {
  return err
}
```
in Wo would be:
```go
file = os!Open("open_me.go")
```
or
```go
file, log("Oh no: ", err) = os.Open(fileName)
```

The point of these features is to drop the bantering about the theories of when to boilerplate or how to be readable, and to just try it out to really see what works well before judgement.

Wo also:
- Supplies native string and slice operations like `==`
- Doesn't allow keyword or import overloading like `var int int = 1` and `rune := 'W'`
- Doesn't prefer shortenings like `f` for `file` or function names like `SprintF` for `ConcatFormat`
- Separates the usage of `var` and `:=` amongst initializing, shadowing, and setting variables
- Doesn't allow undeclared variables / "zero values"
- Warns about unused variables instead of giving an error (then compiles them away)
- Allow function overloading
- Plans to do something about null checking somehow in the future (e.g. nonnull or option)
- Will still commit to universal formatting
- Makes you say "woah"

To justify these decisions, I provide a much deeper analysis of the design at [err.nil](https://err.nil/)


Besides syntactical and formatting difference, it also offers functional differences such as `set`:

```go
map... etc.
```

`set` is meant to be more optimized than any personal implementations using map.


### Trademark disclaimer

All activity here should follow all of Go's guidelines at https://go.dev/brand/. If they inform me that anything violates it, then I will quickly comply. It is also preferable to follow https://go.dev/conduct

Do not refer to Wo as anything other than "a fork of Go" at least not in any way that could disparage the Go programming language.
> Unauthorized Naming Conventions: Naming Conventions that disparage the Go programming language, if not permitted as fair use, are unauthorized.

**This is not a source of the Go programming language nor is it affiliated. It is only a fork.**
> ...and may not inaccurately suggest affiliation or endorsement or mislead as to the source.

Also see:
> Modifications that disparage the Go programming language or its reputation without qualifying as fair use, such as the introduction of malicious code, are not compatible with use of the Go Trademarks.

Additionally, do not associate this with Go's logo or mascot.
> In order to accurately identify the Go programming language or any compatible applications, it may be necessary to refer to the language by name (“nominative fair use”). These are the basic rules for nominative fair use of the Go Trademarks:
> 
>    Only use the Go trademark in word mark form, i.e., plain text. Do not use the Go Logo or Go as a stylized form without permission.
>    Only use the Go trademark as much as is necessary. Use should be limited to matter-of-fact statements.
>    Do not use the Go trademark in any way that suggests or implies affiliation with or endorsement from the community or from Google.
