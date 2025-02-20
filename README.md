### Wo is a fork of Go

The Wo language is an interoperable fork of Go that offers alternative syntax and language features aimed at readability. It's not endorsed by Go.

For example, Go's error handling,

```go
f, err := os.Open("hi.wo")
if err != nil {
    return nil, err
}
```

could be done like this in Wo:

```go
var file = os.Open("hi.wo")!
```

###### Pending decisions here. It's a WIP.

The point of these features is to look beyond banter and theories, and to just **try it out** to really see what works well before judgement. I try iterations of these features before listing them, and these were the most notable options. I hope you find it interesting - definitely feel free to give your own suggestions.

## *Currently, <u>none of these necessarily work yet</u>. It's more of a proof of concept.*

### None are fully working currently, this project is like a few weeks old, and I'm one person with a job.

Also see [justifications.md](https://github.com/wo-language/wo/blob/release-branch.go1.23.wo/justifications.md) and [specification.md](https://github.com/wo-language/wo/blob/release-branch.go1.23.wo/specification.md).

### Syntax

|           Syntax Feature           |                          Go Method                          |                     Wo Example                      |
|:----------------------------------:|:-----------------------------------------------------------:|:---------------------------------------------------:|
|           `interface{}`            |              `interface{Length(interface{})}`               |                   `<Length(<>)>`                    |
|       `interface{\|}` union        |                 `interface{int8 \| int16}`                  |                   `int8 \| int16`                   |
|         Enhanced for loop          | `for i, v := range nums {}`<br/>`for _, v := range nums {}` |       `for i, v : nums {}`<br/>`for v : nums`       |
|         Ternary expression         |       `var v int; if high { v = 99 } else { v = 1 }`        |          `var v = if high then 99 else 1`           |
|     Has conditional assignment     |               `if a, cond := call(); cond {}`               |               `if var a = call() {}`                |
|     `_` for multi return value     |                       `_, val = f()`                        |    `func f() (skip, val, skip2)`<br/>`val = f()`    |
|           Function type            |              `func(int) func(func(), int) int`              |           `int -> (() -> _, int) -> int`            |
|    Single line function literal    |            `func(v int) bool { return v == 0 }`             |          `v -> v == 0`, `() -> effects()`           |
| Keep `func` for multi line literal |     `func() int {`<br/>`unlock()`<br/>`return open() }`     | `func() int {`<br/>`unlock()`<br/>`return open() }` |

### Types & Data

|                                Feature                                | Go Method                                                            | Wo Example                                                                                                                 |
|:---------------------------------------------------------------------:|:---------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------|
|                             Native `set`                              | `map[int]struct{}` and/or self-implementation                        | `var s = set[int] { 2, 7 }; exists = s[2]; s.delete(7)`                                                                    |
| Other native collections like stack, list, treeset, and their atomics | self-implementation                                                  | `stack.pop()`, `tree.remove(n)`                                                                                            |
|                         Functional interface                          | wrap interface with function type                                    | `type Doer funcinter {do()}; func(d Doer) {d.do()}`                                                                        |
|                              `enum` type                              | `iota` and switch case                                               | `type E enum {A(true), B(false); b bool}`<br/>`A.b`=true, `B.name`="B", `A.pos`=0                                          |
|                              flags type                               | `1 << iota` and bit operations                                       | `type F enum {R, W, E}` and bit operations                                                                                 |
|                               sum type                                | `struct` state management                                            | `type File enum { Closed, Open(contents string) }`                                                                         |
|                         Functional interface                          | `interface{ Take(any) }`                                             | Inlines                                                                                                                    |
|                          Algebraic data type                          | `type A interface { int \| string }`                                 | `(num int64 \| ByteNum(set[byte]) + Infinity(bool sign), size int8)` struct(union(type, sum), type)                        |
|            Native `strings`, `maps`, and slice operations             | `strings.Contains(str, substr)`                                      | `str.Contains(substr)`                                                                                                     |
|                        Optional type with `?`                         | `v, ok := m[k]; if ok { }`<br/>`func Get() (int, bool)`              | `v int? = m[k]; v?`<br/>`v int = m[k]?`<br/>`func Get() int?`<br/>`.OrElse(v2)`, `.IsPresent()`, etc.                      |
|                     Errable/Result type with `!`                      | `f, err := Open(n); if err == nil { }`<br/>`func Div() (int, error)` | `file *File! = Open(n); file!`<br/>`file *File = Open(n)!`<br/>`func Div(n, d) int!{}`<br/>`.map(fn)` and `.OrElse(file2)` |
|                    Method generic type parameters                     |                                                                      |                                                                                                                            |

### Variables

|                             Design                              | Go Usage                                                                        | Wo Usage                                                          |
|:---------------------------------------------------------------:|:--------------------------------------------------------------------------------|-------------------------------------------------------------------|
| Doesn't allow **import overloading** or **keyword overloading** | `var int int = 1`, `rune := 'W'`<br/>`import { strings }; var strings []string` | *compile error*                                                   |
|                 **Warns** for unused variables                  | `func f() { x := 1 }`<br/>*compile error*                                       | *warning*                                                         |
|       No accessing uninitialized variables (zero values)        | `var s string // = ""`<br/>`s += "." // s == "."`                               | `var s string`<br/>`s += "." // error`<br/>`var t string? = None` |
|               Assign variables with **only** `=`                | `var e int; e, z := 8, 9; e = 7`                                                | `var e = 0; e = 7`                                                |
|  `:=` for shadowing **only** and not mixed with initialization  | `h := 1; { h, m := 2, 5 }`                                                      | `var h = 1; { h := 2; var m = 5 }`                                |
|             Untyped declaration with **only** `var`             | `var a = 1`, `x := 1`                                                           | `var a = 1`<br/>                                                  |
|       Initialize with type with **only** `=` and no `var`       | `var i int = 2`                                                                 | `i int = 2`                                                       |
|                       No multi assignment                       | `p, q = 20, 30`                                                                 | *compile error*                                                   |

### Language Design

Features that change the functionality of the language beyond syntax and design principles.

| Design                                    | Go Method                                                                                                                                   | Wo Usage                                                                                            |
|:------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------|
| Doesn't prefer **name shortenings**       | `f`<br/>`SprintF`                                                                                                                           | `file`<br/>`ConcatFormat`                                                                           |
| Allows **function overloading**           | `print(string)`<br/>`printF(formatter, string)`<br/>`printOF(stdout, formatter, string)`                                                    | `print(string)`<br/>`print(formatter, string)`<br/>`print(stdout, formatter, string)`               |
| Allows **default arguments** in functions | `print(string, stdout, fmt) {`<br/>&emsp;&emsp;`if fmt == nil {formatter = defFmt}`<br/>&emsp;&emsp;`if stdout == nil {stdout = console} }` | `print(string,`<br/>&emsp;&emsp;`formatter = defaultFormatter,`<br/>&emsp;&emsp;`stdout = console)` |
| Export explicitly                         | `func Export() // capital`, `func Xแมว() // add "X"`<br/>`func private()` `func แมว()`                                                      | `func export แมว()`, `export func InKilos()`<br/>`func private()`                                   |
| Prefers all caps consts                   | `const maxSize = 8`                                                                                                                         | `const MAX_SIZE = 8`                                                                                |
| Export to the package but not globally    | *none*                                                                                                                                      | `func pkg Get()`, `type pkg Bog struct`                                                             |
| Make slice append more predictable        | Overrides / Resizes                                                                                                                         | Indicates new allocs                                                                                |
| Run other functions besides main          | `func main() { other() }`                                                                                                                   | `func otherMain() { }`                                                                              |
| More liberal folder structure             | main, mod                                                                                                                                   | *TBD*                                                                                               |

### To justify these decisions, I provide a deeper analysis of the design at [justifications.md](https://github.com/wo-language/wo/blob/release-branch.go1.23.wo/justifications.md).

And see summaries of the changes in [specification.md](https://github.com/wo-language/wo/blob/release-branch.go1.23.wo/specification.md) and their implementation progress.

I am considering making each language features **modular**. If someone likes only the interface syntax, and that's all they want, then I could allow either compiler flags headers in the file to indicate which ones to have turned off.

Wo also:
- Still commits to a **universal formatting**.
- Is open source and free.
- Is a **WIP**, but will always accept change and criticism.
- Has a **wo**mbat mascot.
- Makes you say **"woah"**.

#### See the list below for several unlikely but possible features:
<details>
<summary>
Potential Features...
</summary>

| Potential Features (Unlikely to be added)                                                                                                          | Go                                                             | Wo Proposal                                                                                                            |
|----------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------|
| Remove `func`, and remove parens from the receiver                                                                                                 | `func (C* c) f[A rune](a int) (float32, error) {}`             | `C.f[rune A](int a) float32? {}`                                                                                       |
| Switch the type with the name in variable and struct [declarations](https://go.dev/blog/declaration-syntax), parameters, and function return types | `i int`                                                        | `int i`, `struct s`, `string proc(float32 f)`                                                                          |
| Don't use `type` in declarations                                                                                                                   | `type A interface {}`                                          | `interface A {}`, `struct B {}`                                                                                        |
| Make it more obvious that map and slice are [pointers](https://dave.cheney.net/2017/04/30/if-a-map-isnt-a-reference-variable-what-is-it)           | `map`                                                          | `*map`                                                                                                                 |
| Allow methods to be in their struct                                                                                                                | `struct Bug { func fly() }`<br/>`func (f F*) flee() {f.fly()}` | `struct Bug { fly()`<br/>`flee() { this.fly() } }`<br/>and/or `struct (Bug* bug) { }` to allow `bug` instead of `this` |
</details>

### How to install

I'd rather `wo` were a lite CLI command that just uses the Go compiler that's already installed rather than needing a different build of the entire compiler, but that'd make interoperability almost impossible, so it is its own compiler.

See details in [dev.md](https://github.com/wo-language/wo/blob/release-branch.go1.23.wo/dev.md). You can install it by building it from this source checked out from the right version (not this one, but the [1.23 branch](https://github.com/wo-language/wo/blob/release-branch.go1.23.wo)), as per https://go.dev/doc/install/source#bootstrapFromCrosscompiledSource.

## Trademark disclaimer

All activity here should follow all of Go's guidelines at https://go.dev/brand/. If they inform me that anything violates it, then I will quickly comply. It is also preferable to follow https://go.dev/conduct.

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
