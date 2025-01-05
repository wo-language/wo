### Wo is a fork of Go

The Wo language is an interoperable successor to Go that offers alternative syntax and language features aimed at readability. To accomodate preference and situation, features will be modular - to be removed at will during compilation.

For example, Go's error handling

```go
f, err := os.Open("hi.wo")
if err != nil {
    return nil, err
}
```

could be done like this in Wo:

```go
var file = os.Open("hi.wo")




```

(Pending decisions here. It's a WIP)

The point of these features is to look beyond banter and theories, such as how much to boilerplate, or whether to do what people are just used to or not - and to instead just **try it out** to really see what works well before judgement. I try iterations of these features, and these were the most notable options. I hope you find them interesting - definitely feel free to give your own suggestions.

## *Currently, this is a <u>proof of concept</u>, and none of these necessarily work yet. (None do currently, this project is like a couple weeks old)*.

## Syntax

|         Syntax Feature         |                            Go Syntax                            |                              Wo Example                              |
|:------------------------------:|:---------------------------------------------------------------:|:--------------------------------------------------------------------:|
|         `interface{}`          |         `interface{Length(interface{}) interface{int}}`         |                         `<Length(<>) <int>>`                         |
|     `interface{\|}` union      |                   `interface{int8 \| int16}`                    |                           `int8 \| int16`                            |
|       Enhanced for loop        |  `for i, v := range nums {}` <br/> `for _, v := range nums {}`  |               `for i, v : nums {}`<br/>`for v : nums`                |
|       Ternary expression       |         `var v int; if high { v = 99 } else { v = 1 }`          |                 `var v = if high { 99 } else { 1 }`                  |
|   Has conditional assignment   |                 `if a, cond := call(); cond {}`                 |           `if var a = call() {}`<br/>(casted to optional)            |
|   `_` for multi return value   |                         `_, val = f()`                          | `func f() (skip, val, skip2)`<br/>`val = f()` (unless it's an error) |
|         Function type          |                 `func(int) func(int, int) int`                  |                      `int -> (int, int) -> int`                      |
|  Single line function literal  |               `func(v) bool -> { return v == 0 }`               |                   `v -> v == 0`, `() -> effects()`                   |
| `func` for multi line or `{ }` | `func() int {`<br/>`unlock()\open()`<br/>`return getFunds()\}`  |    `func() int {`<br/>`unlock()\open()`<br/>`return getFunds()\}`    |

## Language Design

| Design                                    | Go Usage                                                                                                                                    | Wo Usage                                                                                            |
|:------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------|
| Doesn't prefer **name shortenings**       | `f`<br/>`SprintF`                                                                                                                           | `file`<br/>`ConcatFormat`                                                                           |
| Allows **function overloading**           | `print(string)`<br/>`printF(formatter, string)`<br/>`printOF(stdout, formatter, string)`                                                    | `print(string)`<br/>`print(formatter, string)`<br/>`print(stdout, formatter, string)`               |
| Allows **default arguments** in functions | `print(string, stdout, fmt) {`<br/>&emsp;&emsp;`if fmt == nil {formatter = defFmt}`<br/>&emsp;&emsp;`if stdout == nil {stdout = console} }` | `print(string,`<br/>&emsp;&emsp;`formatter = defaultFormatter,`<br/>&emsp;&emsp;`stdout = console)` |

Wo could also address **null checking** somehow (e.g. `nonnil`) and pointer/value receivers. And as for error handling, maybe take inspiration from Rust's [result](https://doc.rust-lang.org/std/result/) or Scala's [canthrow](https://docs.scala-lang.org/scala3/reference/experimental/canthrow.html).

## Types & Data

|                                Feature                                | Go Method                                                            | Wo Example                                                                                                                           |
|:---------------------------------------------------------------------:|:---------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------|
|                             Native `set`                              | `map[int]struct{}` and/or self-implementation                        | `var s = set[int] { 2, 7 }; exists = s[2]; s.delete(7)`                                                                              |
| Other native collections like stack, list, treeset, and their atomics | self-implementation                                                  | `stack.pop()`, `tree.remove(n)`                                                                                                      |
|                         Functional interface                          | wrap interface with function type                                    | `type Doer funcinter {do()}; func(d Doer) {d.do()}`                                                                                  |
|                              `enum` type                              | `iota` and switch case                                               | `type E enum {A(true), B(false); b bool}`<br/>`A.b`=true, `B.name`="B", `A.pos`=0                                                    |
|                              flags type                               | `1 << iota` and bit operations                                       | `type F enum {R, W, E}` and bit operations                                                                                           |
|                                sum type                               | `struct` state management                                            | `type File enum { Closed, Open(contents string) }`                                                                                   |
|                         Functional interface                          |                                                                      |                                                                                                                                      |
|                          Algebraic data type                          | `type A interface { int \| string }`                                 | `(num int64 \| ByteNum(set[byte]) + Infinity(bool sign), size int8)` struct(union(type, sum), type)                                  |
|            Native `strings`, `maps`, and slice operations             | `strings.Contains(str, substr)`                                      | `str.Contains(substr)`                                                                                                               |
|                        Optional type with `?`                         | `v, ok := m[k]; if ok { }`<br/>`func Get() (int, bool)`              | `v int? = m[k]; v?`<br/>`v int = m[k]?`<br/>`func Get() int?`<br/>`.OrElse(v2)`, `.IsPresent()`, etc.                                |
|                     Errable/Result type with `!`                      | `f, err := Open(n); if err == nil { }`<br/>`func Div() (int, error)` | `file *File! = Open(n); file!`<br/>`file *File = Open(n)!` (must check)<br/>`func Div() int!`<br/>`.OrElse(file2)`, `.Erred()`, etc. |
|                    Method generic type parameters                     |                                                                      |                                                                                                                                      |

## Variables

|                               Design                                | Go Usage                                                                        | Wo Usage                               |
|:-------------------------------------------------------------------:|:--------------------------------------------------------------------------------|----------------------------------------|
|   Doesn't allow **import overloading** or **keyword overloading**   | `var int int = 1`, `rune := 'W'`<br/>`import { strings }; var strings []string` | *compile error*                        |
|                   Warns for **unused variables**                    | `func f() { x := 1 }`<br/>*compile error*                                       | *warning*                              |
|                No accessing uninitialized variables                 | `var s string // = ""`<br/>`s += "." // s == "."`                               | `var s string`<br/>`s += "." // error` |
|                 Assign variables with **only** `=`                  | `var e int; e, z := 8, 9; e = 7`                                                | `var e = 0; e = 7`                     |
| `:=` for shadowing **only** and not when mixing with initialization | `h := 1; { h, m := 2, 5 }`                                                      | `var h = 1; { h := 2; var m = 5 }`     |
|               Untyped declaration with **only** `var`               | `var a = 1`, `x := 1`                                                           | `var a = 1`<br/>                       |
|         Initialize with type with **only** `=` and no `var`         | `var i int = 2`                                                                 | `i int = 2`                            |
|                         No multi assignment                         | `p, q = 20, 30`                                                                 | *compile error*                        |

## Language Features

Features that change the functionality of the language beyond syntax and design principles.

| Feature                                | Go Method                                                                             | Wo Usage                                                      |
|:---------------------------------------|---------------------------------------------------------------------------------------|---------------------------------------------------------------|
| Export explicitly                      | `func Export() // apital`, `func Xแมว() // add "X"`<br/>`func private()` `func แมว()` | `func export แมว()`, `export const Kilo`<br/>`func private()` |
| Export to the package but not globally | *none*                                                                                | `func pkg Get()`, `type pkg Bog struct`                       |
| Make slice append more predictable     | Overrides / Resizes                                                                   | Indicates new allocs                                          |
| Run other functions besides main       | `func main() { other() }`                                                             | `func otherMain() { }`                                        |
| More liberal folder structure          | main, mod                                                                             | TBD                                                           |

### To justify these decisions, I provide a deeper analysis of the design at ~~[err.nil](https://err.nil/)~~ [justifications.md](/justifications.md).

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
| Use the arrow return style in `func`s, and for function types                                                                                      | `var f func(func(float64) int) string`                         | `(float64 -> int) -> string f`                                                                                         |
| Switch the type with the name in variable and struct [declarations](https://go.dev/blog/declaration-syntax), parameters, and function return types | `i int`                                                        | `int i`, `struct s`, `string proc(float32 f)`                                                                          |
| Have tuples as an assignable type                                                                                                                  | `a, b; return a, b`                                            | `t; return t`                                                                                                          |
| Don't use `type` in declarations                                                                                                                   | `type A interface {}`                                          | `interface A {}`, `struct B {}`                                                                                        |
| Make it more obvious that map and slice are [pointers](https://dave.cheney.net/2017/04/30/if-a-map-isnt-a-reference-variable-what-is-it)           | `map`                                                          | `*map`                                                                                                                 |
| Allow methods to be in their struct                                                                                                                | `struct Bug { func fly() }`<br/>`func (f F*) flee() {f.fly()}` | `struct Bug { fly()`<br/>`flee() { this.fly() } }`<br/>and/or `struct (Bug* bug) { }` to allow `bug` instead of `this` |
</details>

### How to install

I'd rather `wo` were a lite CLI command that just uses the Go compiler that's already installed rather than needing a different build of the entire compiler, but that'd make interoperability almost impossible, so it is its own compiler.

You can install it by building it from this source checked out from the right version, as per https://go.dev/doc/install/source#bootstrapFromCrosscompiledSource. Currently, it's not guaranteed to work.

### Status

To summarize which code has actually been written, I've implemented tiny parts of a couple things, and the inner functionality and tree structure of `set`, but not its operators yet to test it. I have gotten it (the compiler) to compile and run ".wo" files.

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
