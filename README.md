### Wo is a fork of Go

The Wo language is an interoperable successor to Go. It offers alternative syntax and language features aimed at readability.

For example,

```go
f, err := os.Open("hi.wo")
if err != nil {
    return nil, err
}
```

would be done like this in Wo:

```go
file var = os.Open!("hi.wo") // pending decisions here; it's a WIP
```

Some other ways to handle errors in Wo:

```go
file var = os.Open!!("hi.wo") // panic
file, log("Error:", err)   var = os.Open("hi.wo")
file, handle(err)          var = os!Open("hi.wo") // handle and throw
file, return(none, 3, err) var = os.Open("hi.wo") // with other return values
file, if(err)              var = os.Open("hi.wo") { handle(err) } // similar to Swift's `try?`
```

The point of these features is to drop the bantering about the theories of how much to boilerplate or whether to copy what people have been used to, and to just **try it out** to really see what works well before judgement. I've tried iterations of this myself, and these were the most notable options, but it is a **proof of concept** and I have not necessarily got any of these working yet.


| Rule                                                                            | Usage                                                                                                                                            |
|---------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------|
| Uses `<>`, not `interface{}` in type parameters                                 | `f(a interface{})` -> `f(a <>)` or `interface{Length() int}` to `<Length() int>`                                                                 |
| Allows **function overloading**                                                 | `print(string)`<br/>`print(formatter, string)`<br/>`print(stdout, formatter, string)`                                                            |
| ...but does **allow** **default arguments** in functions anyway                 | `print(stdout = console,`<br/>&emsp;`formatter = defaultFormatter,`<br/>&emsp;`string)`<br/>like how `[:]` already does: `slice(start=0, end=0)` |
| **Doesn't** allow import overloading or **keyword overloading**                 | `var int int = 1` and `rune := 'W'` give compile error                                                                                           |
| **Doens't** use "`range`" in **enhanced for** loops `for i, v := range nums {}` | `for i, v : nums {}`<br/>`for v : nums` (values instead of `_, v`)                                                                               |
| **Doesn't** prefer name shortenings                                             | `f` for `file` or function names like `SprintF` for `ConcatFormat` (isn't enforced)                                                              |
| Has the **ternary operator** for<br/>`if a, cond := call(); cond {}`            | `a, if(cond) = call() {} else {}`<br/>or maybe `?:` and  `a, cond? = call() {} : {}`                                                             |
| Uses `[]` after the type for arrays to reflect `map[key]`                       | **`int[]`** and `int[...][3]`                                                                                                                    |

(In the future) Wo also...
- Reworks variables by
    - Not giving an **error for unused variables**, (just warns and compiles them away)
    - Not allowing undeclared variables or **"zero values"** like `var x`
    - Allow optionals for when the zero value would have had a double meaning like `string?` to avoid `""`
      - However, I could allow zero value initialization with `none`, like `int x` would mean `int x = none`
    - Separating the usage of **`var`, `:=`, and `=`** amongst initializing, shadowing, and setting variables without any overlapping functionality
      - Uses `=` for initialization and setting, requiring `:=` for only shadowing, and then use **`i int = 5`** syntax for initializing or `var i = 5` for untyped variables
    - Making `_, val = f()` redundant (like `for i = range` has it optionally) by accessing only specific values from multi-return values: `w, o = f()` where `func f() (w, skip, o)`
      - Unless `f` were to return an `error`, maybe requiring something like `val, err <!= f()` when that could happen
- Will still commit to universal formatting
- Is a **WIP**, but will always accept change and criticism
- Makes you say **"woah"**

Wo is open source and free. Yes, the mascot is a **wo**mbat.

Besides syntactical and formatting difference, Wo also offers functional differences such as
- A native implementation of `set`
- Could address **null checking** somehow (e.g. `nonnull` or `option`) and pointer/value receivers
- Error values: a few potential options: not returning `nil` if there is no error, but something like `status.isErr()` being true, maybe like Rust's [result](https://doc.rust-lang.org/std/result/). Or `error` overriding all other return values like an exception: `io.Read` returns either `n` or throws `error` like `errable io.Read() n` or [canthrow](https://docs.scala-lang.org/scala3/reference/experimental/canthrow.html)
- (compile time) explicit `enum` and `union`
- Make slice append more predictable
- Tuples as an assignable type
- Native string and slice operations like `==` and `"".contains`
- Run other functions besides main / package control

### See the list below for other several other more unlikely features:
<details>
<summary>
Potential Features
</summary>

- *MAYBE* remove `func`, and remove parens from the receiver like `func (C* c) f[A rune](a int) (float32, error) {}` to `C.f[rune A](int a) float32? {}`
- Signify errored outputs like `f() errable (int, string)` means `f() error? | (int, string)?` where only one is some and the other is none
- Use the arrow return style in `func`s, and for function types: `var f func(func(float64) int) string` for `(float64 -> int) -> string f`
- *Undecided* whether to switch the type with the name in variable and struct [declarations](https://go.dev/blog/declaration-syntax), parameters, and function return types like `int i`, `struct s`, `string proc(float32 f)`
- *MAYBE* don't use `type` from `type A interface {}`
- *MAYBE* Make it more obvious that map and slice are pointers? https://dave.cheney.net/2017/04/30/if-a-map-isnt-a-reference-variable-what-is-it
- *MAYBE* (probably won't) allow methods to be in their struct like
- *MAYBE* not allowing **mixing shadowed** and initialized variable declarations


`struct Bug { func fly() }   func (f F*) flee() {f.fly()}` ->

`struct Bug { fly()   flee() { this.fly() } }` or `struct (Bug* bug) { }` to allow `bug` instead of `this`

</details>

To justify these decisions, I provide a deeper analysis of the design at ~~[err.nil](https://err.nil/)~~ justifications.md for now

I'd rather `wo` were a lite CLI command that just uses the Go compiler that's already installed rather than needing a different build of the entire compiler, but I'm making it a separate build for now.

### Code example

```go
import { "strings" }
type FilePath interface {
  string | url
}

type Program struct {
  executable [...]byte
}

func (p Program*) output() string {
  return p.executable[:strings.LastIndex(p.executable, ".exe"))
}

func (p Program*) len() int {
  return len(p.executable)
}

func runProgram() string {
  output, err = runProgram("/")
  if err != nil {
    log(err)
  }
  return output
}

var fs = map[FilePath]string{"/app/host": "server.ts", "/", "Main.java"}

func runProgramO(dir interface{string|url}) (int, *string, error) {
  f, ok = fs[dir]
  if (!ok) {
    return nil, errors.New("invalid filepath")
  }

  r, err := os.Open(f)
  if err != nil {
    return nil, err
  }

  defer func() {
    if err := r.Close(); err != nil {
      return nil, err
    }
  }()

  if err := reader.Sync(); err != nil {
    return nil, err
  }

  p := myCompiler.build(reader)

  return p.len(), *p.outputPath(), nil
}
```
a possible design for Wo:
```go
runProgram(<string|url> directory) -> errable (int, string) { // members reversed to order by relevancy
  fileName, if(!ok) = runnableFiles[directory] {
    return errors.New("invalid filepath") // like throw
  }
  reader *File = os.Open(fileName)
  defer reader.Close!()
  reader.Sync!()
  program Program = myCompiler.build(reader)
  return program.outputPath() // converts it to some(string) and error as nil/none
}

runnableFiles = map[FilePath, string]{"/app/host": "server.ts", "/", "Main.java"}

runProgram() -> string {
  output, log(err) = runProgram("/")
  return output
}

Program struct {
  byte[...] executable
  outputPath() -> string {
    return executable[:executable.LastIndex(".exe")]
  }

  len() -> int {
    len(executable)
  }
}

FilePath interface {
  string | url
}
```


|go| wo with types before name                                                                                                                                                |
|------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|<pre>var fs = map[FilePath]string</pre>| <pre>map[FilePath, string] runnableFiles</pre>                                                                                                                           |
|<pre>func runProgramO(dir interface{string:url})</pre>| <pre>runProgram(<string:url> directory)</pre>                                                                                                                            |
|<pre>&emsp;f, ok = fs[dir]<br>&emsp;if (!ok) {<br>&emsp;&emsp;return nil, errors.New("invalid filepath")<br>&emsp;}</pre>| <pre>&emsp;fileName, if(!ok) = runnableFiles[directory] {<br>&emsp;&emsp;throw error("invalid filepath")<br>&emsp;}<br><br></pre>                                        |
|<pre>&emsp;r, err := os.Open(f)<br>&emsp;if err != nil {<br>&emsp;&emsp;return nil, err`<br>&emsp;}| <pre>&emsp;*File reader = os.Open!(fileName)<br><br><br><br></pre>                                                                                                       |
|<pre>&emsp;defer func() {<br>&emsp;&emsp;if err := r.Close(); err != nil {<br>&emsp;&emsp;&emsp;return nil, err<br>&emsp;&emsp;}<br>&emsp;}()</pre>| <pre>&emsp;defer reader.Close!()<br><br><br><br><br></pre>                                                                                                               |
|<pre>type Program struct {<br>&emsp;executable [...]byte<br>}<br>func (p Program*) output() string {<br>&emsp;return p.executable[:strings.LastIndex(p.executable, ".exe"))<br>}</pre>| <pre>struct Program {<br>&emsp;byte[...] executable<br>&emsp;string outputPath() {<br>&emsp;&emsp;return executable[:executable.LastIndex(".exe")]<br>&emsp;}<br>}</pre> |


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
