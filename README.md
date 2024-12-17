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

A few other potential ways:
```c
file = os!!Open("hi.wo") // panic
file, log("couldn't open:", err) = os.Open("hi.wo")
file, handle(err) = os.Open("hi.wo")
file, return(none, 3, err) = os.Open("hi.wo") // with other return values
file, if(err) = os.Open("hi.wo") { handle(err) } // similar to Swift's `try?`
```

The point of these features is to drop the bantering about the theories of when to boilerplate or how to be readable, and to just try it out to really see what works well before judgement.

Wo also...
- Uses `interface{}` for `<>` in type parameters, e.g. `f(a interface{})` -> `f(a <>)` or `interface{Length() int}` to `<Length() int>`
- Allows **function overloading** like `print(string), print(formatter, string), print(stdout, formatter, string)`
- But allows **default arguments** in functions anyway like `print(stdout = console, formatter = defaultFormatter, string)`
- Doesn't allow import or **keyword overloading** like `var int int = 1` and `rune := 'W'`
- Doens't use `range` in **enhanced for** loops or `_,` to ignore the index like `for _, v := range nums {}` for `for v : nums {}`
- Doesn't prefer shortenings like `f` for `file` or function names like `SprintF` for `ConcatFormat` (isn't enforced)
- Reworks variables by
  - not giving an **error for unused variables**, (just warn and compile them away)
  - not allowing undeclared variables or **"zero values"**
  - *MAYBE* make `_, val = f()` redundant by accessing only specific values from multi-return values -> `val = f()` where `val` matches the name in `func f() (other, val)` unless it is returning an `error`, maybe needing something like `<=` when that happens
  - *MAYBE* removing **mixing shadowed** and initialized variable declarations together
  - separating the usage of **`var`, `:=`, and `=`** amongst initializing, shadowing, and setting variables without any overlapping functionality
  - *MAYBE* use `=` for initialization and setting, requiring `:=` for shadowing (but not `for range`), and then use **`int i = 5`** (good old C) syntax for initializing with `var i = 5` for vague untyped variables (or just remove untyped variables syntactically)
- *MAYBE* switch type with the name of parameters, put the return types before the function, remove `func`, use `errable`, and generic types before the fuction name like `func (c C*) f[A rune](a int) (float32, error) {}` to `float32 (C* c) [rune A] f(int a) errable {}` or arrow style, `(C* c) [rune A] f(int a) -> float32 | error {}` (or `!float32`)
  - and *MAYBE* do similarly for function types: `var f func(func(float64) int) string` for `string func(int func(float64)) f`, `string f(int _(float64))`, or `(float64 -> int) -> string f`
- Uses `interface A {}`, as well as `struct B {}`, unlike `type A interface {}`
- *MAYBE* allow methods to be in their struct like `struct Bug { func fly() }   func (f F*) flee() {f.fly()}` for `struct Bug { fly()   flee() { this.fly() } }` or `struct (Bug* bug) { }` to allow `bug` instead of `this`
- Has the **ternary operator** `if cond {} else {}` (or ?: upon more deliberation) and `if a, cond := call(); cond {}` for `a, if(cond) = call() {}` or maybe `a, cond? = call() : {}`
- Uses `[]` after the type for arrays like **`int[]`** and `int[...][3]`. A `map` of arrays would be ambiguous, so `map[int][]int` becomes **`map[int, int[]]`** and `map` in general uses `[A, B]`
- Allows optionals for when the zero value has a double meaning like `string?` or `File?` to not be `""` or `nil` which means I could allow zero value initialization without declaration by setting it to `none` like `int x` would mean `int x = none`
- *MAYBE* Make it more obvious that map and slice are pointers like `*map[string, string]`
- Will still commit to universal formatting
- Is a **WIP**, but will always accept change and criticism
- Makes you say **"woah"**

To justify these decisions, I provide a deeper analysis of the design at [err.nil](https://err.nil/)

Besides syntactical and formatting difference, Wo also offers functional differences such as
- a native `set`, which is meant to be more optimized than implementations using map
- could address **null checking** somehow (e.g. `nonnull` or `option`) and pointer/value receivers
- error values: a few potential options: not returning `nil` if there is no error, but something like `status.isErr()` being true, maybe like Rust's [result](https://doc.rust-lang.org/std/result/). Or `error` overriding all other return values like an exception: `io.Read` returns either `n` or throws `error` like `errable io.Read() n` or [canthrow](https://docs.scala-lang.org/scala3/reference/experimental/canthrow.html)
- (complie time) **enum**
- native string and slice operations like `==` and `.contains`
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
func runProgram() string {
  output, err = runProgram("/")
  if err != nil {
    log(err)
  }
  return output
}
var fs = map[FilePath]string{"/app/host": "server.ts", "/", "Main.java"}
func runProgramO(dir interface{string|url}) (*string, error) {
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
    return *p.outputPath(), nil
}
```
a possible design for Wo:
```c
string? runProgram(<string|url> directory) errable { // members reversed to order by relevancy
    fileName, if(!ok) = runnableFiles[directory] {
      return errors.New("invalid filepath") // like throw
    }
    *File reader = os!Open(fileName)
    defer reader!Close()
    reader!Sync()
    Program program = myCompiler.build(reader)
    string directory := program.outputPath() // shadowing
    return directory // converts it to some(string) and error as nil/none
}
map[FilePath, string] runnableFiles = {"/app/host": "server.ts", "/", "Main.java"}
string runProgram() {
  output, log(err) = runProgram("/")
  return output
}
struct Program {
  byte[...] executable
  string outputPath() {
    return executable[:executable.LastIndex(".exe")]
  }
}
interface FilePath {
  string | url
}
```

Again, the types before variable names is TBD. It isn't really a problem, but just fits with other syntax choices.

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
