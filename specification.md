
### Index

## Conventions

### Variable naming

There is one situation where shortened variable names might be acceptable, which is lambda function calls like `starfruits.map(s -> s.weight() * 2.2)`, or generally very short function calls, where you can easily see what `s` means.

Wo also does not allow variable names like `π` or `__`.


### Renaming package methods

`Print, Printf, Sprint, Sprintf, Fprint, Fprintf, Sscanf, Fscanf,` etc.

`PrintFormat, Concat, ConcatFormat, FormatterPrint, ScanString, ScanReader` etc.

## Syntax Features

### interface{}

```go
func f(a <>, b <bool>) {
}
```

### Unused variables

```go
func main() {
    var x = 3 // warning: Unused variable 'x'
}
```

### For range

```go
for i, v : nums {
    sum += v
}
```

### Variable declaration

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


### Overloading package names

```go
import { "strings" }

func combineThem(strings /* Wo Error */ []string) string {
    return strings.Join /* Go error */ (strings, ", ")
}
```

## Overloading functions

With the same function names:

```go
payWith(cash)
payWith(creditCardInfo)
payWith(creditCardNumber, zipCode)
payWith(creditCardNumber, city, state)
```


### Ternary

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


```go
primes set[int] = { 2, 3, 5 }  // declaration
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

`struct S {}`
`interface I {}`

### Scope control
