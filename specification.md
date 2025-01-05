## *Remember, many of these things are pending, and waiting to be tested properly, and any of these could be scrapped or altered!*

The point of this file is to give a minimal depiction of each feature. See [justifications.md](/justifications.md) for the theory and elaboration of each section.

[Examples in code](https://github.com/wo-language/wo-info/blob/main/examples/specification.wo)

### Index

1. [Operators](#Operators List)
2. [Syntax](#Syntax)
   1. [Interface](#Interface)
   2. [For range](#For range)
   3. [Ternary](#Ternary)
   4. [Array/Slice](#Array/Slice)
   5. [Map](#Map)
   6. [Function](#Function)
   7. [Type keyword](#Type keyword)
3. [Data Types](#Data Types)
   1. [Set](#Set)
   2. [Optional](#Optional)
   3. [Complex](#i)
4. [Data Models](#Data Models)
   1. [Enum](#Enum)
   2. [Flags](#Flags)
   3. [Sum](#Sum)
   4. [Union](#Union)
   5. [Functional interface](#Functional interface)
   6. [Overview](#Overview)
   7. [Algebraic types](#Algebraic types)
   8. [Pattern Matching](#Pattern Matching)
5. [Generics](#Generics)
   1. [Parameterized methods](#Parameterized methods)
6. [Variables](#Variables)
   1. [Unused variables](#Unused variables)
   2. [Variable declaration](#Variable declaration)
   3. [Multi variable declaration](#Multi variable declaration)
7. [Error handling](#Error handling)
   1. [nil](#nil)
8. [Design](#Design)
   1. [Standard library](#Standard library)
   2. [Naming](#Naming)
      1. [Variables](#Variables)
      2. [Package methods](#Package methods)
   3. [Overleading](#Overleading)
      1. [Package names](#Package names)
      2. [Reserved words](#Reserved words)
      3. [Functions](#Functions)
   4. [Keywords](#Keywords)
   5. [Import compatibility](#Import compatibility)
   6. [Export](#Export)
   7. [Scope control](#Scope control)
   8. [Array/Slice Clarity](#Array/Slice Clarity)
   9. [Modularity](#Modularity)


## Operators List

### Added

| op            | syntax                                                           |
|---------------|------------------------------------------------------------------|
| ENHANCEDFOR   | for Key, Value : X { Body } variant of ORANGE                    |
| UNWRAP        | X?                                                               |
| OPTION        | X.Type?                                                          |
| UNWRAPERR     | X!                                                               |
| ERRABLE       | X.Type!                                                          |
| ARROWCLOSURE  | Type -> { Func.Closure.Body } variant of OCLOSURE                |
| DCLFUNC       | export? (r)? func f() - modification                             |
| ENUMLIT       | Type(List) (composite literal, Type is enum) - done in enum decl |
| DCLSHADOW     | X; X := Y                                                        |
| INTERFACETAGS | <Type{List}>                                                     |

### Removed

| op     | replacement |
|--------|-------------|
| ORANGE | ENHANCEDFOR |

### Modularity

`Disable[Set]`

## Syntax

### Interface

`interface{T}` -> `<T>`

Also see [union](#Union).

### For range

`for i, v := range vs` -> `for i, v : vs`

`for _, v := range vs` -> `for v : vs`

`for i := range vs` -> `for i, _ : vs`

### Ternary

`if cond then A else B`

### Array/Slice

`[]elem` (no changes)

### Map

`map[K]V` (no changes)

### Function

Multi line (no changes):

```go
func f() {
    line1()
    return result()
}
```

Single line:

`func()` -> `() -> _`,
`func() {}` -> `() -> {}`

`func(I) O` -> `I -> O`,
`func(I) O { o() }` -> `i -> o()`

`func(i I) (o O)` -> `(i I) -> (o O)`

`func(f func()) (r R)` -> `(f() -> _) -> (r R)`

`func() func()` -> `() -> () -> _`,
`func() func() { return func() {}}` -> `() -> () -> {}`

### Type keyword

No changes.

## Data Types

### Set

```go
primes set[int] = { 2, 3, 5 }  // declaration
ok = primes[4]                 // is ok if contains elem
primes.insert[7]               // insert / add
primes.delete[3]               // delete / remove
```

### Optional

`Some(v)`

`None`

`T?` - `Option[T]`

`Option(v)?` - unwrap to v if Some, panic if `None`

`IsPresent() bool`

`Map() Option`

`func f() (T, bool)` -> `func f() T?` // formatting AND interpretation from .go functions

### i

No changes.

## Data Models

### Enum

```go
type Enum enum {
    EnumVal1(values...), EnumVal2(values...), EnumVal3(values...), ... // delim by commas or newlines

    fields... // types must match with values, same style as struct fields
}
```

```go
type Data enum {
    Val1(5, false, "A", "B")
    Val2(-2.2, true)
    Val3(0e3, true, "aaaaaaaa")
    
    float32 // unnamed
    ok bool // named
    ...string // varargs
}
```

### Flags

```go
type Flags enum {
  Flag1 // 1
  Flag2 // = 2
  Flag3 // = 4
}
```

`Flags.Full` = 0b111 = 7

`Flag1 & Flag2 &^ Flags.Full | Flag3` // bin ops

### Sum

```go
type Sum enum {
  SumVal1 // no parens or args
  SumVal2() // empty parens to differentiate with flags
  SumVal3(fields)
  SumVal4(ExternalEnumVal, uint, named ...uint) // access uint with `SumVal4.uint`
}

```

### Union

`interface{A | B}` -> `A | B`

### Functional interface

`interface{f()}` -> `<f()>`

`interface{func()}` -> `<() -> _>`

### Data Models

`struct`, tuple, `interface`, Union `interface`, Functional `interface`, `enum`, Flags `enum`, Sum `enum` (Wo)

Example [go](https://github.com/wo-language/wo-info/blob/main/examples/go/datamodels.go) and [wo](https://github.com/wo-language/wo-info/blob/main/examples/datamodels.wo) file

### Algebraic types

Precedence: `+` < `|` < `,`

Class hierarchy: Sum <- Enum <- Flags

```go
type Sing = int                  // type
type Uni  = int8 | int16         // union
type Nest = Sing | Uni | float32 // union
type Flip = On + Off             // flags
type IntO = None + val Some(int) // sum
type Pt   = (x, y int)           // tuple
type Pair = (f Flip, IntO)       // tuple

type FlipNest = (On + Off, None + Some(int | int8 | int16)) // tuple(flags, sum)
type Combine  = ([]int, string) | A + B | <Field(int)>                    // union(tuple, sum, functional interface)
type Tagged   = (a []int, b string) | c A + B | <d Field(int)>                    // union(tuple, sum, functional interface)
```

### Pattern Matching

```go
type Length enum {
    Cm(float32)
    M(float32)
    FtInch(int, float32) 
}

func (length Length) ToInches() float32 {
    return switch length {
        case Cm(cm) => cm / 2.54
        case M(m) => m / 0.0254
        case FtInch(ft, in) => float32(ft) * 12 + in
    }
}
```

## Generics

### Parameterized methods



## Variables

### Unused variables

```go
func main() {
    var x = 3 // error: Unused variable 'x'
}
```

### Variable declaration



### Multi variable declaration



## Error Handling

### nil

```go
type Errable[T] = T + error

T: no error occured
err: an error occured, err
```

```go
var file                       = os.Open("hi.wo")!  // return err
var file, log("Error:", err)   = os.Open("hi.wo")
var file, handle(err)          = os.Open("hi.wo")!  // handle and throw
var file, return(none, 3, err) = os.Open("hi.wo")   // with other return values
var file, if(err)              = os.Open("hi.wo") { handle(err) } // similar to Swift's `try?`
if var file                    = os.Open("hi.wo") { /*main code*/ }    // Swift/Rust
var file                       = os.Open("hi.wo")!! // panic
var file                       = os.Open("hi.wo")?  // unwrap or panic
var file                       = os.Open("hi.wo").orElse(newFile)
//var file                       = os.Open("hi.wo")? else newFile
```

`func f() (T, err)` -> `func f() T!` // formatting AND interpretation from .go functions

## Design

### Standard library

[sets](/src/sets/sets.go), [set](/src/runtime/set.go), option, enum, and collections

Unexported:

[set_fast32](/src/runtime/set_fast32.go), [set_fast64](/src/runtime/set_fast64.go), [set_faststr](/src/runtime/set_faststr.go)

### Variable naming

Full names of variables like `file` and `fileName`

### Renaming package methods

Full, unabbreviated names of functions like `ConcatFormat` for `SprintF`

`Print, Printf, Sprint, Sprintf, Fprint, Fprintf, Sscanf, Fscanf`

`PrintFormat, Concat, ConcatFormat, FormatterPrint, ScanString, ScanReader`

### Overloading

#### Overloading package names

Can't name a variable the same as an existing package name:

```go
import ( "strings" )

func combineThem(strings /* Wo Error */ []string) string {
    return strings.Join /* Go error */ (strings, ", ")
}
```

#### Overloading reserved words

I assume one of the reasons it allows overloading reserved words (`int`, `nil`) is because of backwards compatability, which means I don't need that since this is a fresh start for syntax. Allowing the ability to override those is always confusing and unsafe. Words spelled the same with different meanings used in the same exact contexts, which can be done by accident, is confusing. Enough said.

#### Overloading functions

```go
payWith(cash)
payWith(creditCardInfo)
payWith(creditCardNumber, zipCode)
payWith(creditCardNumber, city, state)
```

### Import compatibility

`$ident` to import something from Go that overlaps with a reserved word in Wo

### Export

`export func Sew(string) {}`

### Scope control

`func innerSew(string) {}` // is not visible publicly

### Array/Slice Clarity


