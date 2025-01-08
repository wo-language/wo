## *Remember, many of these things are pending, and waiting to be tested properly, and any of these could be scrapped or altered!*

### This doc is heavily subjected to change

The point of this file is to give a minimal depiction of each feature. See [justifications.md](/justifications.md) for the theory and elaboration of each section.

[Examples from this file in code](https://github.com/wo-language/wo-info/blob/main/examples/specification.wo)

### Index

1. [Operators](#Operators-List)
2. [Syntax](#Syntax)
   1. [Interface](#Interface)
   2. [For range](#For-range)
   3. [Ternary](#Ternary)
   4. [Array/Slice](#ArraySlice)
   5. [Map](#Map)
   6. [Function](#Function)
   7. [Type keyword](#Type-keyword)
3. [Data Types](#Data-Types)
   1. [Set](#Set)
   2. [Option](#Option)
   3. [Complex](#i)
4. [Data Models](#Data-Models)
   1. [Tab Enum](#Tab-Enum)
   2. [Flags](#Flags)
   3. [Sum](#Sum)
   4. [Union](#Union)
   5. [Functional interface](#Functional-interface)
   6. [Algebraic types](#Algebraic-types)
   7. [Pattern Matching](#Pattern-Matching)
5. [Generics](#Generics)
   1. [Parameterized methods](#Parameterized-methods)
6. [Variables](#Variables)
   1. [Unused variables](#Unused-variables)
   2. [Variable declaration](#Variable-declaration)
   3. [Multi variable declaration](#Multi-variable-declaration)
7. [Error handling](#Error-handling)
   1. [nil](#nil)
8. [Design](#Design)
   1. [Standard library](#Standard-library)
   2. [Variables](#Variable-naming)
   3. [CONST naming](#CONST)
   4. [Package methods](#Renaming-package-methods)
   5. [Overloading](#Overloading)
      1. [Package names](#Overloading-package-names)
      2. [Reserved words](#Overloading-Reserved-words)
      3. [Functions](#Overloading-Functions)
   6. [Import compatibility](#Import-compatibility)
   7. [Export](#Export)
   8. [Scope control](#Scope-control)
   9. [Array/Slice clarity](#ArraySlice-clarity)
   10. [Modularity](#Modularity)

### Priorities

0. A test feature / necessary starting feature
1. Most important - to do after a test feature
2. Very High importance
3. High
4. Important
5. Medium 
6. Low 
7. To consider later on

## Operators List

### Added

Operators added from Go base:

| op            | syntax                                           |
|---------------|--------------------------------------------------|
| Set           |                                                  |
| ENHANCEDFOR   | for Key, Value : X { Body } variant of RANGE     |
| UNWRAP        | X?                                               |
| OPTION        | X.Type?                                          |
| UNWRAPERR     | X!                                               |
| ERRABLE       | X.Type!                                          |
| ARROWCLOSURE  | Type -> { Func.Closure.Body } variant of CLOSURE |
| DCLFUNC       | export? (r)? func f() - modification             |
| DCLSHADOW     | X; X := Y                                        |
| INTERFACETAGS | <Type{List}>                                     |
| ENUMDECL      | type X enum { ...ENUMLIT }                       |
| ENUMLIT       | X(List) \| X                                     |
| ENUMADD       | ENUMLIT + ENUMLIT \| ENUMADD + ENUMLIT           |
| TUPLE         | (X, Y...)                                        |

### Replaced (Modified)

| op    | replacement |
|-------|-------------|
| RANGE | ENHANCEDFOR |

### Modularity

###### Not Implemented. Priority: 2

`Disable[Set]`

## Syntax

### Interface

###### Not Implemented. Priority: 4

`interface{T}` &#8594; `<T>`

Also see [union](#Union).

### For range

###### Not Implemented. Priority: 5

`for i, v := range vs` &#8594; `for i, v : vs`

`for _, v := range vs` &#8594; `for v : vs`

`for i := range vs` &#8594; `for i, _ : vs`

### Ternary

###### Not Implemented. Priority: 4

`if cond then A else B`

### Array/Slice

`[]elem` (no changes)

### Map

`map[K]V` (no changes)

### Function

###### Not Implemented. Priority: 3

Multi line (no changes):

```go
func f() {
    line1()
    return result()
}
```

Single line:

`func()` &#8594; `() -> _`,
`func() {}` &#8594; `() -> {}`

`func(I) O` &#8594; `I -> O`,
`func(I) O { o() }` &#8594; `i -> o()`

`func(i I) (o O)` &#8594; `(i I) -> (o O)`

`func(f func()) (r R)` &#8594; `(f() -> _) -> (r R)`

`func() func()` &#8594; `() -> () -> _`,
`func() func() { return func() {}}` &#8594; `() -> () -> {}`

### Type keyword

No changes.

### Struct tags



## Data Types

### Set

###### Partially Implemented. Priority: 0

```go
primes set[int] = { 2, 3, 5 }  // declaration
ok = primes[4]                 // is ok if contains elem
primes.insert[7]               // insert / add
primes.delete[3]               // delete / remove
```

### Option

###### Not Implemented. Priority: 3

`Some(v)`

`None`

`T?` - `Option[T]`

`Option(v)?` - unwrap to v if Some, panic if `None`

`IsPresent() bool`

`Map() Option`

`func f() (T, bool)` &#8594; `func f() T?` // formatting AND interpretation from .go functions

## Error Handling

### nil

###### Not Implemented. Priority: 3

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

`func f() (T, err)` &#8594; `func f() T!` // formatting AND interpretation from .go functions

### i

No changes.

## Data Models

`struct`, tuple, `interface`, Union `interface`, Functional `interface`, Sum `enum`, Tab `enum`, Flags `enum`

Example [go](https://github.com/wo-language/wo-info/blob/main/examples/go/datamodels.go) and [wo](https://github.com/wo-language/wo-info/blob/main/examples/datamodels.wo) file

### Tab Enum

###### Not Implemented. Priority: 3

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

###### Not Implemented. Priority: 3

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

###### Not Implemented. Priority: 3

```go
type SumName enum {
    SumVal1 // no parens or args
    SumVal2() // empty parens to differentiate with flags
    SumVal3(fields)
    SumVal4(ExternalEnumVal, uint, named ...uint) // access uint with `.uint`
}
```

```go
return switch sumNameVal {
    case SumVal1 => 0
    case SumVal2() => 8
    case SumVal3(m, _, o) => m * o
    case SumVal4(other, num, named) =>{
        switch other {
        Mul => num * named[0]
        Add => num + named[0]
        }
    }
}
```

### Union

###### Not Implemented. Priority: 4

`interface{A | B}` &#8594; `A | B`

### Functional interface

###### Not Implemented. Priority: 4

`interface{f()}` &#8594; `<f()>`

`interface{func()}` &#8594; `<() -> _>`

### Algebraic types

###### Not Implemented. Priority: 5

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

type FlipNest = (On + Off, None + Some(int | int8 | int16))    // tuple(flags, sum)
type Combine  = ([]int, string) | A + B | <Field(int)>         // union(tuple, sum, functional interface)
type Tagged   = (a []int, b string) | c A + B | <d Field(int)> // union(tuple, sum, functional interface)
```

### Pattern Matching

###### Not Implemented. Priority: 

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

###### Not Implemented. Priority: 5

## Variables

### Unused variables

###### Not Implemented. Priority: 0

```go
func main() {
    var x = 3 // error: Unused variable 'x'
}
```


### Variable declaration

###### Not Implemented. Priority: 4



### Multi variable declaration

###### Not Implemented. Priority: 4



## Design

### Standard library

###### Not Implemented. Priority: 5

[sets](/src/sets/sets.go), [set](/src/runtime/set.go), option, enum, and collections

Unexported:

[set_fast32](/src/runtime/set_fast32.go), [set_fast64](/src/runtime/set_fast64.go), [set_faststr](/src/runtime/set_faststr.go)

### Variable naming

Full names of variables like `file` and `fileName`

### CONST

`const fawn = 1` &#8594; `const FAWN = 1` (style)

### Renaming package methods

###### Not Implemented. Priority: 5

Full, unabbreviated names of functions like `ConcatFormat` for `SprintF`

`Print, Printf, Sprint, Sprintf, Fprint, Fprintf, Sscanf, Fscanf`

`PrintFormat, Concat, ConcatFormat, FormatterPrint, ScanString, ScanReader`

### Overloading

#### Overloading package names

###### Not Implemented. Priority: 6

Can't name a variable the same as an existing package name:

```go
import ( "strings" )

func combineThem(strings /* Wo Error */ []string) string {
    return strings.Join /* Go error */ (strings, ", ")
}
```

#### Overloading reserved words

###### Not Implemented. Priority: 5

I assume one of the reasons it allows overloading reserved words (`int`, `nil`) is because of backwards compatability, which means I don't need that since this is a fresh start for syntax. Allowing the ability to override those is always confusing and unsafe. Words spelled the same with different meanings used in the same exact contexts, which can be done by accident, is confusing. Enough said.

#### Overloading functions

###### Not Implemented. Priority: 4

```go
payWith(cash)
payWith(creditCardInfo)
payWith(creditCardNumber, zipCode)
payWith(creditCardNumber, city, state)
```

### Import compatibility

###### Not Implemented. Priority: 4

`$ident` to import something from Go that overlaps with a reserved word in Wo

### Export

###### Not Implemented. Priority: 6

`export func Sew(string) {}`

### Scope control

###### Not Implemented. Priority: 6

`func innerSew(string) {}` // is not visible publicly

### Array/Slice clarity

###### Not Implemented. Priority: 6


