### Wo is a fork of Go

The Wo language offers an alternative syntax and functionality to the Go programming language and interoperates with Go.

Here's one example:
```go
f, err := os.Open("hi.wo")
if err != nil {
  return err
}
```

would be done like this:

```c
file = os!Open("hi.wo")
```

And it would return with any other return values filled in as their zero value.

Some other ways to deal with error handling:
```c
file = os!!Open("hi.wo") // panic
file, log("couldn't open:", err) = os.Open("hi.wo")
file, handle(err) = os.Open("hi.wo")
file, return(none, 3, err) = os.Open("hi.wo") // with other return values
file, if(err) = os.Open("hi.wo") { handle(err) } // similar to Swift's `try?`
```

The point of these features is to drop the bantering about the theories of when to boilerplate, how to be readable, whether to copy what people are used to, and to just try it out to really see what works well before judgement.

(In the future) Wo also...
- Uses `interface{}` for `<>` in type parameters, e.g. `f(a interface{})` -> `f(a <>)` or `interface{Length() int}` to `<Length() int>`
- Allows **function overloading** like `print(string), print(formatter, string), print(stdout, formatter, string)`
- But allows **default arguments** in functions anyway like `print(stdout = console, formatter = defaultFormatter, string)` (like [:] is already doing: `slice(start=0, end=0)`)
- Doesn't allow import or **keyword overloading** like `var int int = 1` and `rune := 'W'`
- Doens't use `range` in **enhanced for** loops like `for i, v := range nums {}` for `for i, v : nums {}`, and *MAYBE* make the default the value, not index: `for v : nums`
- Doesn't prefer shortenings like `f` for `file` or function names like `SprintF` for `ConcatFormat` (isn't enforced)
- Reworks variables by
  - not giving an **error for unused variables**, (just warns and compiles them away)
  - not allowing undeclared variables or **"zero values"**
  - allow optionals for when the zero value would have had a double meaning like `string?` or `File?` to not be `""` or `nil` which means I could allow zero value initialization without declaration by setting it to `none` like `int x` would mean `int x = none`

  - *MAYBE* make `_, val = f()` redundant by accessing only specific values from multi-return values: `val = f()` where `val` matches the name in `func f() (other, val)` unless `f` were returning an `error`, maybe requiring something like `val, err <!= f()` when that could happen
  - *MAYBE* not allowing **mixing shadowed** and initialized variable declarations
  - separating the usage of **`var`, `:=`, and `=`** amongst initializing, shadowing, and setting variables without any overlapping functionality
  - *MAYBE* use `=` for initialization and setting, requiring `:=` for shadowing (but not `for range`), and then use **`i: int = 5`** with type syntax for initializing or `var i = 5` for vague untyped variables (or just remove untyped variables syntactically)
- *MAYBE* remove `func`, use "`errable`", and put remove parens from the receiver like `func (C* c) f[A rune](a int) (float32, error) {}` to `C.f[rune A](int a) errable float32? {}`
  - arrow return style `-> float32 | error` (or `!float32` or !> for errable) `f() errable (int, string)` means `f() error? | (int, string)?` where only one is some and the other is none
  - and *MAYBE* do similarly for function types: `var f func(func(float64) int) string` for `f(_(float64) int) string`, or `(float64 -> int) -> string f`
  - *Undecided* whewher to switch the type with the name in variable and struct [declarations](https://go.dev/blog/declaration-syntax), parameters, and function return types like `int i`, `struct s`, `string proc(float32 f)`
- Doesn't use `type` from `type A interface {}` (if the base compiler will allow it)
- *MAYBE* (probably won't) allow methods to be in their struct like `struct Bug { func fly() }   func (f F*) flee() {f.fly()}` for `struct Bug { fly()   flee() { this.fly() } }` or `struct (Bug* bug) { }` to allow `bug` instead of `this`
- Has the **ternary operator** `if cond {} else {}` (or `?:` upon more deliberation) and `if a, cond := call(); cond {}` for `a, if(cond) = call() {}` or maybe `a, cond? = call() : {}`
- Uses `map[A, B]`
- Uses `[]` after the type for arrays like **`int[]`** and `int[...][3]`. A `map` of arrays would be ambiguous, so `map[int][]int` becomes **`map[int, int[]]`**
- *MAYBE* Make it more obvious that map and slice are pointers like `*map[string, string]`
- Will still commit to universal formatting
- Is a **WIP**, but will always accept change and criticism
- Makes you say **"woah"**

To justify these decisions, I provide a deeper analysis of the design at [err.nil](https://err.nil/)

Besides syntactical and formatting difference, Wo also offers functional differences such as
- a native `set`, which is meant to be more optimized than implementations using map
- could address **null checking** somehow (e.g. `nonnull` or `option`) and pointer/value receivers
- error values: a few potential options: not returning `nil` if there is no error, but something like `status.isErr()` being true, maybe like Rust's [result](https://doc.rust-lang.org/std/result/). Or `error` overriding all other return values like an exception: `io.Read` returns either `n` or throws `error` like `errable io.Read() n` or [canthrow](https://docs.scala-lang.org/scala3/reference/experimental/canthrow.html)
- (complie time) **enum**s and unions
- native string and slice operations like `==` and `"".contains`
- being able to run other functions besides main

I'd rather `wo` were a lite CLI command that just uses the Go compiler rather than needing a different build of the entire compiler, but I'm making it a separate build for now.

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
```c
runProgram(<string|url> directory) -> errable (int, string)? { // members reversed to order by relevancy
    fileName, if(!ok) = runnableFiles[directory] {
      return errors.New("invalid filepath") // like throw
    }
    reader *File = os!Open(fileName)
    defer reader!Close()
    reader!Sync()
    program Program = myCompiler.build(reader)
    directory := program.outputPath() // shadowing
    return directory // converts it to some(string) and error as nil/none
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
Comparison with types and names switched:
|go|wo|
|---------|--------|
|<pre>var fs = map[FilePath]string</pre>|<pre>map[FilePath, string] runnableFiles</pre>|
|<pre>func runProgramO(dir interface{string:url}) (*string, error) {</pre>|<pre>string? runProgram(<string:url> directory) errable {</pre>|
|<pre>&emsp;f, ok = fs[dir]<br>&emsp;if (!ok) {<br>&emsp;&emsp;return nil, errors.New("invalid filepath")<br>&emsp;}</pre>|<pre>&emsp;fileName, if(!ok) = runnableFiles[directory] {<br>&emsp;&emsp;throw error("invalid filepath")<br>&emsp;}<br><br></pre>|
|<pre>&emsp;r, err := os.Open(f)<br>&emsp;if err != nil {<br>&emsp;&emsp;return nil, err`<br>&emsp;}|<pre>&emsp;*File reader = os!Open(fileName)<br><br><br></pre>|
|<pre>&emsp;defer func() {<br>&emsp;&emsp;if err := r.Close(); err != nil {<br>&emsp;&emsp;&emsp;return nil, err<br>&emsp;&emsp;}<br>&emsp;}()</pre>|<pre>&emsp;defer reader!Close()<br><br><br><br></pre>|
|<pre>type Program struct {<br>&emsp;executable [...]byte<br>}<br>func (p Program*) output() string {<br>&emsp;return p.executable[:strings.LastIndex(p.executable, ".exe"))<br>}</pre>|<pre>struct Program {<br>&emsp;byte[...] executable<br>&emsp;string outputPath() {<br>&emsp;&emsp;return executable[:executable.LastIndex(".exe")]<br>&emsp;}<br>}</pre>|

Yes, the mascot is a **wo**mbat.

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
