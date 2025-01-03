
### Index

### Modular annotation

## Conventions

### Variable naming

file, fileName

### Renaming package methods

`Print, Printf, Sprint, Sprintf, Fprint, Fprintf, Sscanf, Fscanf,` etc.

`PrintFormat, Concat, ConcatFormat, FormatterPrint, ScanString, ScanReader` etc.

## Syntax Features

### interface{}

[] Not implemented

```go
func f(a <>, b <bool>) {
}
```

### Unused variables

[] Not implemented

```go
func main() {
    var x = 3 // warning: Unused variable 'x'
}
```

#### Other soft errors

[] Implemented - need to test
- Only ignoring unused

Any of these could be ignored easily in theory:

unused var in switch
no new variables on left side of :=
generic function is missing function body
label %s declared and not used
main - func %s must have no type parameters, func %s must have no arguments and no return values
init - missing function body
%q imported and not used, %q imported as %s and not used
can only use ... with final parameter in list
cannot range over, range over %s permits no iteration variables, range over %s permits only one iteration variable, range clause permits at most two iteration variables
cannot use iteration variable of type %s
cannot use type %s outside a type constraint: interface is (or embeds) comparable
in interface - overlapping terms %s and %s


src/cmd/compile/internal/typecheck/stmt.go
src/internal/types/errors/codes.go
src/cmd/compile/internal/types2/errors.go

### For range

[] Not implemented

```go
for i, v : nums {
    sum += v
}
```

### Variable declaration

[] Not implemented

```go
x int = 8
a var = 8
y const int = 9 // maybe: const y int = 9
b const = y // maybe: b const var = y
x = b
{ x := 0 } // shadow
var ( z, _ (int, error) = count() )
const ( b string = "永", e error = nil )
```

### Multi variable declaration

[] Not implemented

```go
getPoint := func() (x int, y int) { return 1, 2 }
a, b := getPoint()
_, c := getPoint()
d, _ := getPoint()
```
but you would not be able to do `p := getPoint()` nor could you do `y := getPoint()` to select just the `y` part of the return values. One pattern is to use a struct, but let's see what it could look like in a longer flow of functions

```go

```

```go
a, b = 1, 2
```

### Overloading reserved words

[] Not implemented

### Overloading package names

[] Not implemented

```go
import { "strings" }

func combineThem(strings /* Wo Error */ []string) string {
    return strings.Join /* Go error */ (strings, ", ")
}
```

## Overloading functions

[] Not implemented

With the same function names:

```go
payWith(cash)
payWith(creditCardInfo)
payWith(creditCardNumber, zipCode)
payWith(creditCardNumber, city, state)
```


### Ternary

[] Not implemented

```go
if indoorTemp == outdoorTemp {
    hvac = off
} else if indoorTemp < outdoorTemp {
    hvac = heating
} else {
    hvac = ac
}

var hvac = if indoorTemp == outdoorTemp { off } else if indoorTemp < outdoorTemp { heating } else { ac }

### Util

`strings.contains(str, sub)`

vs

`str.contains(sub)`


```go
x == y
```

### set

[] Started implementing

```go
primes hashset[int] = { 2, 3, 5 }  // declaration
ok = primes[4]                 // is ok if contains elem
primes.insert[7]               // insert / add
primes.delete[3]               // delete / remove
```

#### if err != nil

```go
none: no error occured
some(err), an error occured, and it is err
```

```go
type Errable[T] interface {
  t T
  err error
}
```

### enum

[] Not implemented

```go
type Days enum {
    Sunday("sun", false),
    Monday("moon", true)

    root    string,
    workday bool
}

Sunday.position
Monday.working
```

### type

[] Not implemented

`struct S {}`
`interface I {}`

### Scope control

[] Not implemented

## Operators List

### Added

| op            | syntax                                                           |
|---------------|------------------------------------------------------------------|
| SOME          | some(X)                                                          |
| NONE          | none                                                             |
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
 // TODO(bran) 

### Removed

| op     | replacement |
|--------|-------------|
| ORANGE | ENHANCEDFOR |
